package milestone

import (
	"crypto/rand"
	"crypto/sha256"
	"strconv"
	"strings"

	"github.com/oklog/ulid/v2"
)

// idLen is the length of a milestone ID. Six Crockford characters give 32^6
// values, so a hundred milestones collide with a chance far below one in a
// thousand; four digits would have made that one in three.
const idLen = 6

// NewID rolls a fresh milestone ID: idLen random characters of the ticket
// handle alphabet. An ID of digits only is rolled again, because it could not
// be told from the tail of a name such as sprint-202611.
func NewID() string {
	for {
		var b [idLen]byte
		if _, err := rand.Read(b[:]); err != nil {
			panic(err)
		}
		var sb strings.Builder
		for _, c := range b {
			sb.WriteByte(ulid.Encoding[c&31])
		}
		if id := sb.String(); hasLetter(id) {
			return id
		}
	}
}

// IsID reports whether s is a milestone ID: exactly six characters, upper
// case, from the ticket handle alphabet, at least one of them a letter.
func IsID(s string) bool {
	if len(s) != idLen {
		return false
	}
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(ulid.Encoding, s[i]) < 0 {
			return false
		}
	}
	return hasLetter(s)
}

func hasLetter(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > '9' {
			return true
		}
	}
	return false
}

// LegacyID is the ID a milestone written before IDs existed is given: the
// first 30 bits of the SHA-256 of its old file name. It is computed rather
// than rolled so that every clone, offline and without talking to anybody,
// moves the same old file to the same ID. A hash of digits only is hashed
// again with a counter until a letter is in it.
func LegacyID(name string) string {
	in := name
	for n := 1; ; n++ {
		sum := sha256.Sum256([]byte(in))
		v := uint32(sum[0])<<24 | uint32(sum[1])<<16 | uint32(sum[2])<<8 | uint32(sum[3])
		var b [idLen]byte
		for i := range b {
			b[i] = ulid.Encoding[(v>>(2+5*uint(idLen-1-i)))&31]
		}
		if id := string(b[:]); hasLetter(id) {
			return id
		}
		in = name + "\n" + strconv.Itoa(n)
	}
}

// SplitFile reads a milestone file's base name (without .md) as <name>-<ID>.
// ok is false for a file written before IDs existed, whose whole base name is
// the name.
func SplitFile(base string) (name, id string, ok bool) {
	i := strings.LastIndexByte(base, '-')
	if i <= 0 || !IsID(base[i+1:]) {
		return "", "", false
	}
	return base[:i], base[i+1:], true
}

// FileNamed is the test for a file base name (no .md) that belongs to a
// milestone called name, with or without an ID in it.
func FileNamed(name string) func(base string) bool {
	return func(base string) bool {
		if base == name {
			return true
		}
		n, _, ok := SplitFile(base)
		return ok && n == name
	}
}
