package hashcash

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

const (
	defaultVersion = 1
	zeroBit        = '0'
	dateLayout     = "20060102150405"
)

// HashCash - hashcash structure (https://en.wikipedia.org/wiki/Hashcash)
type HashCash struct {
	version   int       // format version, 1 (which supersedes version 0)
	bits      int       // number of "partial pre-image" (zero) bits in the hashed code
	date      time.Time // time that the message was sent, in the format YYMMDD[hhmm[ss]]
	resource  string    // resource data string being transmitted, e.g., an IP address or email address
	extension string    // extension field, optional; ignored in version 1
	rand      []byte    // random characters, encoded in base-64 format
	counter   int       // binary counter, encoded in base-64 format.
}

func New(ctx context.Context, bits int) (*HashCash, error) {
	rand, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		return nil, err
	}

	if bits <= 0 {
		return nil, errors.New("zero bits must be more than zero")
	}

	return &HashCash{
		bits:     bits,
		date:     time.Now().UTC().Truncate(time.Second),
		resource: ctx.Value("clientID").(string),
		rand:     rand.Bytes(),
	}, nil
}

func (h *HashCash) IsValid(ctx context.Context, ttl time.Duration) (bool, error) {
	if !h.equalResource(ctx.Value("clientID").(string)) {
		return false, errors.New("hashcash resource is not equal with clientID")
	}

	if !h.isActual(ttl) {
		return false, errors.New("hashcash is not actual")
	}

	isHashCorrect, err := h.Header().IsHashCorrect(h.Bits())
	if err != nil {
		return false, errors.New("hashcash is not correct")
	}
	if !isHashCorrect {
		return false, errors.New("hashcash is not correct")
	}

	return true, nil
}

// Bits - returns number of zero bits
func (h *HashCash) Bits() int {
	return h.bits
}

// Counter - returns counter
func (h *HashCash) Counter() int {
	return h.counter
}

// EqualResource - check if input resource is equal with hashcash resource
func (h *HashCash) equalResource(resource string) bool {
	return h.resource == resource
}

// IsActual - check if hashcash expiration exceeded ttl
func (h *HashCash) isActual(ttl time.Duration) bool {
	return h.date.Add(ttl * time.Second).After(time.Now().UTC())
}

// Compute - compute hash with enough zero bits in the begining
// Increase counter if hash does't have enough zero bits in the begining
func (h *HashCash) Compute(maxAttempts int) error {
	if maxAttempts > 0 {
		h.counter = 0
		for h.counter <= maxAttempts {
			ok, err := h.Header().IsHashCorrect(h.bits)
			if err != nil {
				return err
			}
			if ok {
				return nil
			}

			h.counter++
		}
	}

	return errors.New("max attempts exceeded")
}

// Header - returns string presentation of hashcash to share it
func (h *HashCash) Header() Header {
	return Header(fmt.Sprintf("%d:%d:%s:%s:%s:%s:%s",
		defaultVersion,
		h.bits,
		h.date.Format(dateLayout),
		h.resource,
		h.extension,
		base64.StdEncoding.EncodeToString(h.rand),
		base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.counter))),
	))
}
