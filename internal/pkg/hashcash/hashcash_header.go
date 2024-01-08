package hashcash

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseHeader - parse hashcah from header
func ParseHeader(header string) (hashcash *HashCash, err error) {
	parts := strings.Split(header, ":")

	if len(parts) < 7 {
		return nil, errors.New("incorrect header format")
	}
	if len(parts) > 7 {
		for i := 0; i < len(parts)-7; i++ {
			parts[3] += ":" + parts[3+i+1]
		}
		parts[4] = parts[len(parts)-3]
		parts[5] = parts[len(parts)-2]
		parts[6] = parts[len(parts)-1]
		parts = parts[:7]
	}
	if parts[0] != "1" {
		return nil, errors.New("incorrect header format")
	}

	hashcash = &HashCash{}

	hashcash.bits, err = strconv.Atoi(parts[1])
	if err != nil {
		return nil, errors.New("incorrect header format")
	}

	hashcash.date, err = time.ParseInLocation(dateLayout, parts[2], time.UTC)
	if err != nil {
		return nil, errors.New("incorrect header format")
	}

	hashcash.resource = parts[3]
	hashcash.extension = parts[4]

	hashcash.rand, err = base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, errors.New("incorrect header format")
	}

	counterStr, err := base64.StdEncoding.DecodeString(parts[6])
	if err != nil {
		return nil, errors.New("incorrect header format")
	}

	hashcash.counter, err = strconv.Atoi(string(counterStr))
	if err != nil {
		return nil, errors.New("incorrect header format")
	}

	return
}

// Header - string presentation of hashcash
// Format - 1:bits:date:resource:externsion:rand:counter
type Header string

// IsHashCorrect - check header hash contain zero bits
func (header Header) IsHashCorrect(bits int) (ok bool, err error) {
	if bits <= 0 {
		return false, errors.New("zero bits must be more than zero")
	}

	hash, err := header.sha1()
	if err != nil {
		return ok, err
	}
	if len(hash) < bits {
		return false, errors.New("incorrect hash length")
	}

	ok = true
	for _, s := range hash[:bits] {
		if s != zeroBit {
			ok = false
			break
		}
	}

	return
}

func (header Header) sha1() (hash string, err error) {
	hasher := sha1.New()
	if _, err = hasher.Write([]byte(header)); err != nil {
		return
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
