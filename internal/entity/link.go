package entity

import (
	"encoding/binary"
	"errors"
)

type Link struct {
	ShortLink string `gorm:"primary_key;size:10;not null;"`
	FullLink  string `gorm:"size:100;not null;"`
}

func NewLink(fullLink string) (*Link, error) {
	if fullLink == "" {
		return nil, errors.New("cannot create link from empty string")
	}
	return &Link{
		ShortLink: CreateLink(fullLink),
		FullLink:  fullLink,
	}, nil
}

func CreateLink(s string) string {
	const seed uint64 = 0x749e3e6989df617

	b := make([]byte, 8)
	hash := murmurOAAT64([]byte(s), seed)

	for !validate(hash) {
		binary.LittleEndian.PutUint64(b, hash)
		hash = murmurOAAT64(b, seed)
	}

	return hash2string(hash)
}

// convert 64bit number to 10-length string
// last 4 bits are omitted
// every 6 bits represents a character (a-zA-Z0-9_):
// 0-9 = 0-9
// 10-35 = A-Z
// 36-61 = a-z
// 62 = _
func hash2string(hash uint64) string {
	s := make([]byte, 10)
	var b uint64

	for i := 0; i < 10; i++ {
		b = hash & 63

		if b < 10 {
			s[i] = byte(48 + b)
		} else if b < 36 {
			s[i] = byte(55 + b)
		} else if b < 62 {
			s[i] = byte(61 + b)
		} else {
			s[i] = byte(95)
		}

		hash >>= 6
	}

	return string(s)
}

// check if there is no 6-bit character with value 63 in hash
// there is only 63 characters (0-62) and 63 is forbidden
func validate(hash uint64) bool {
	var b uint64
	for i := 0; i < 10; i++ {
		b = hash & 63
		if b == 63 {
			return false
		}
		hash >>= 6
	}
	return true
}

// murmur one-byte-at-a-time 64-bit implementation
func murmurOAAT64(data []byte, seed uint64) uint64 {
	var hash = seed
	for _, b := range data {
		hash ^= uint64(b)
		hash *= 0x5bd1e9955bd1e995
		hash ^= hash >> 47
	}
	return hash
}
