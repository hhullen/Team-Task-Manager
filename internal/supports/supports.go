package supports

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	defaultHashLen = 97
)

type HashSettings struct {
	memory     uint32
	iterations uint32
	threads    uint8
	saltLength uint32
	keyLength  uint32
}

var hashSettings = HashSettings{
	memory:     64 * 1024,
	iterations: 3,
	threads:    4,
	saltLength: 16,
	keyLength:  32,
}

func ArgonHash(s string) (string, error) {
	salt := make([]byte, hashSettings.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(s), salt, hashSettings.iterations,
		hashSettings.memory, hashSettings.threads, hashSettings.keyLength)

	var b strings.Builder
	b.Grow(defaultHashLen)

	b.WriteString("$argon2id$v=")
	b.WriteString(strconv.Itoa(argon2.Version))
	b.WriteString("$m=")
	b.WriteString(strconv.FormatUint(uint64(hashSettings.memory), 10))
	b.WriteString(",t=")
	b.WriteString(strconv.FormatUint(uint64(hashSettings.iterations), 10))
	b.WriteString(",p=")
	b.WriteString(strconv.FormatUint(uint64(hashSettings.threads), 10))
	b.WriteByte('$')
	b.WriteString(base64.RawStdEncoding.EncodeToString(salt))
	b.WriteByte('$')
	b.WriteString(base64.RawStdEncoding.EncodeToString(hash))

	return b.String(), nil
}

func IsStringArgonHash(s, hash string) (bool, error) {
	vals := strings.Split(hash, "$")

	if len(vals) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	memory, iterations, parallelism, err := ParseParams(vals[3])
	if err != nil {
		return false, err
	}

	unb64salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false, err
	}

	unb64hash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return false, err
	}

	comparisonHash := argon2.IDKey([]byte(s), unb64salt, iterations, memory, parallelism, uint32(len(unb64hash)))

	return subtle.ConstantTimeCompare(unb64hash, comparisonHash) == 1, nil
}

func ParseParams(paramsStr string) (m, t uint32, p uint8, err error) {
	parts := strings.Split(paramsStr, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid params format")
	}

	for _, part := range parts {
		val, err := strconv.ParseUint(part[2:], 10, 32)
		if err != nil {
			return 0, 0, 0, err
		}

		switch part[0:1] {
		case "m":
			m = uint32(val)
		case "t":
			t = uint32(val)
		case "p":
			p = uint8(val)
		}
	}
	return m, t, p, nil
}
