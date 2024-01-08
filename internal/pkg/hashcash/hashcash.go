package hashcash

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

const (
	zeroBit    = '0'
	dateLayout = "20060102150405"
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

func New(bits int, resource string) (*HashCash, error) {
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
		resource: resource,
		rand:     rand.Bytes(),
	}, nil
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
func (h *HashCash) EqualResource(resource string) bool {
	return h.resource == resource
}

// IsActual - check if hashcash expiration exceeded ttl
func (h *HashCash) IsActual(ttl time.Duration) bool {
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

// Key - returns string presentation of hashcash without counter
// Key is using to match original hashcash with solved hashcash
func (h *HashCash) Key() string {
	return fmt.Sprintf("%d:%d:%s:%d", h.bits, h.date.Unix(), h.resource, binary.BigEndian.Uint32(h.rand))
}

// Header - returns string presentation of hashcash to share it
func (h *HashCash) Header() Header {
	return Header(fmt.Sprintf("1:%d:%s:%s:%s:%s:%s",
		h.bits,
		h.date.Format(dateLayout),
		h.resource,
		h.extension,
		base64.StdEncoding.EncodeToString(h.rand),
		base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.counter))),
	))
}
