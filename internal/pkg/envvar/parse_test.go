package envvar_test

import (
	"encoding/hex"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"go-starter/internal/pkg/envvar"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseInt(t *testing.T) {
	testCases := []struct {
		given  string
		exp    int
		expErr bool
	}{
		{given: "5", exp: 5},
		{given: "30", exp: 30},
		{given: "0.5", expErr: true},
		{given: "abc", expErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			// Given:
			key := "TEST_INT"
			t.Setenv(key, tc.given)

			// When:
			act, err := envvar.ParseInt(key)

			// Then:
			if tc.expErr {
				require.Error(t, err)
				assert.Zero(t, act)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.exp, act)
		})
	}
}

func TestParseDuration(t *testing.T) {
	testCases := []struct {
		given  string
		exp    time.Duration
		expErr bool
	}{
		{given: "5s", exp: 5 * time.Second},
		{given: "30m", exp: 30 * time.Minute},
		{given: "2h", exp: 2 * time.Hour},
		{given: "60", expErr: true},
		{given: "abc", expErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			// Given:
			key := "TEST_DURATION"
			t.Setenv(key, tc.given)

			// When:
			act, err := envvar.ParseDuration(key)

			// Then:
			if tc.expErr {
				require.Error(t, err)
				assert.Zero(t, act)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.exp, act)
		})
	}
}

func TestParseURL(t *testing.T) {
	testCases := []struct {
		given  string
		exp    string
		expErr string
	}{
		{given: "http://localhost:8080", exp: "http://localhost:8080"},
		{given: "localhost:8080", exp: "localhost:8080"},
		{given: "\"http://localhost:8080\"", expErr: `invalid url for TEST_URL: parse "\"http://localhost:8080\"": first path segment in URL cannot contain colon`},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			// Given:
			key := "TEST_URL"
			t.Setenv(key, tc.given)

			// When:
			act, err := envvar.ParseURL(key)

			// Then:
			if tc.expErr != "" {
				require.Error(t, err)
				assert.Equal(t, tc.expErr, err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.exp, act.String())
		})
	}
}

func TestParseHexString(t *testing.T) {
	testCases := []struct {
		given  string
		exp    []byte
		expErr bool
	}{
		{given: hex.EncodeToString([]byte("ABC")), exp: []byte("ABC")},
		{given: "ABC", expErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			// Given:
			key := "TEST_HEXSTRING"
			t.Setenv(key, tc.given)

			// When:
			act, err := envvar.ParseHexString(key)

			// Then:
			if tc.expErr {
				require.Error(t, err)
				assert.Zero(t, act)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.exp, act)
		})
	}
}

func TestParseBool(t *testing.T) {
	testCases := []struct {
		given  string
		exp    bool
		expErr bool
	}{
		{given: "TRUE", exp: true},
		{given: "true", exp: true},
		{given: "T", exp: true},
		{given: "t", exp: true},
		{given: "f", exp: false},
		{given: "F", exp: false},
		{given: "false", exp: false},
		{given: "FALSE", exp: false},
		{given: "FALSE", exp: false},
		{given: "abc", expErr: true},
		{given: "123", expErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			// Given:
			key := "TEST_BOOLSTRING"
			t.Setenv(key, tc.given)

			// When:
			act, err := envvar.ParseBool(key)

			// Then:
			if tc.expErr {
				require.Error(t, err)
				assert.Zero(t, act)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.exp, act)
		})
	}
}

func TestParseOptionalEnvFunc(t *testing.T) {
	testCases := []struct {
		given  string
		setenv bool
		exp    bool
		expErr bool
	}{
		{given: "TRUE", setenv: true, exp: true},
		{given: "true", setenv: true, exp: true},
		{given: "T", setenv: true, exp: true},
		{given: "t", setenv: true, exp: true},
		{given: "f", setenv: true, exp: false},
		{given: "F", setenv: true, exp: false},
		{given: "false", setenv: true, exp: false},
		{given: "FALSE", setenv: true, exp: false},
		{given: "FALSE", setenv: true, exp: false},
		{given: "abc", setenv: true, expErr: true},
		{given: "123", setenv: true, expErr: true},
		// no env var
		{given: "TRUE", setenv: false},
		{given: "true", setenv: false},
		{given: "T", setenv: false},
		{given: "t", setenv: false},
		{given: "f", setenv: false},
		{given: "F", setenv: false},
		{given: "false", setenv: false},
		{given: "FALSE", setenv: false},
		{given: "FALSE", setenv: false},
		{given: "abc", setenv: false},
		{given: "123", setenv: false},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			// Given:
			if tc.setenv {
				t.Setenv("TEST_OPTIONALENVFUNC", tc.given)
			}

			// When:
			act, err := envvar.ParseOptionalEnvFunc("TEST_OPTIONALENVFUNC", strconv.ParseBool)

			// Then:
			if tc.expErr {
				require.Error(t, err)
				assert.Zero(t, act)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.exp, act)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type exp struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	testCases := []struct {
		given  string
		exp    exp
		expErr bool
	}{
		{given: `{}`, exp: exp{}},
		{given: `{"wrong_key":"value"}`, exp: exp{}},
		{
			given: `{"a":"aaaa","b":"bbbb"}`,
			exp: exp{
				A: "aaaa",
				B: "bbbb",
			},
		},
		{given: "ABC", expErr: true},
		{given: "[]", expErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			// Given:
			key := "TEST_UNMARSHALJSON"
			t.Setenv(key, tc.given)

			// When:
			var act exp
			err := envvar.UnmarshalJSON(key, &act)

			// Then:
			if tc.expErr {
				require.Error(t, err)
				assert.Zero(t, act)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.exp, act)
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	testCases := []struct {
		given  string
		exp    slog.Level
		expErr bool
	}{
		{given: `debug`, exp: slog.LevelDebug},
		{given: `InFo`, exp: slog.LevelInfo},
		{given: "invalid", expErr: true},
		{given: "", expErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.given, func(t *testing.T) {
			// Given:
			key := "TEST_UNMARSHALTEXT"
			t.Setenv(key, tc.given)

			// When:
			var act slog.Level
			err := envvar.UnmarshalText(key, &act)

			// Then:
			if tc.expErr {
				require.Error(t, err)
				assert.Zero(t, act)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.exp, act)
		})
	}
}
