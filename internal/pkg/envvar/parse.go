package envvar

import (
	"encoding"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/mail"
	"net/url"
	"os"
	"strconv"
	"time"
)

func ParseURL(key string) (*url.URL, error) {
	u, err := url.Parse(os.Getenv(key))
	if err != nil {
		return nil, fmt.Errorf("invalid url for %s: %w", key, err)
	}
	return u, nil
}

func ParseEmailAddress(key string) (mail.Address, error) {
	a, err := mail.ParseAddress(os.Getenv(key))
	if err != nil {
		return mail.Address{}, fmt.Errorf("invalid email address: %w", err)
	}
	return *a, nil
}

func ParseDuration(key string) (time.Duration, error) {
	d, err := time.ParseDuration(os.Getenv(key))
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
	}
	return d, nil
}

func ParseDateTime(key string, format string) (time.Time, error) {
	d, err := time.Parse(format, os.Getenv(key))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date time for %s: %w", key, err)
	}
	return d, nil
}

var ErrIntOverflow = errors.New("integer overflow")

func ParseInt(key string) (int, error) {
	n, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0, fmt.Errorf("invalid integer for %s: %w", key, err)
	}
	return n, nil
}

func ParseInt32(key string) (int32, error) {
	v, err := ParseInt(key)
	if err != nil {
		return 0, err
	}

	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, ErrIntOverflow
	}

	return int32(v), nil
}

func ParseHexString(key string) ([]byte, error) {
	h, err := hex.DecodeString(os.Getenv(key))
	if err != nil {
		return nil, fmt.Errorf("not a hex string for %s: %w", key, err)
	}
	return h, nil
}

func ParseBool(key string) (bool, error) {
	b, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false, fmt.Errorf("not a valid bool for %s: %w", key, err)
	}
	return b, nil
}

func OptionalString(key string, format ...func(string) string) *string {
	s, ok := os.LookupEnv(key)
	if !ok {
		return nil
	}

	for _, it := range format {
		s = it(s)
	}

	return &s
}

func ParseOptionalEnvFunc[T any](key string, f func(string) (T, error)) (T, error) {
	s, ok := os.LookupEnv(key)
	if !ok {
		var empty T
		return empty, nil
	}
	return f(s)
}

func UnmarshalJSON(key string, a any) error {
	err := json.Unmarshal([]byte(os.Getenv(key)), &a)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json: %w", err)
	}
	return nil
}

func UnmarshalText(key string, a encoding.TextUnmarshaler) error {
	err := a.UnmarshalText([]byte(os.Getenv(key)))
	if err != nil {
		return fmt.Errorf("cannot unmarshal text: %w", err)
	}
	return nil
}
