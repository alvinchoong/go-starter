package ptr_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"go-starter/internal/pkg/ptr"

	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		given any
	}{
		{"min int", math.MinInt},
		{"max int", math.MaxInt},
		{"empty string", ""},
		{"some string", "something"},
		{"char", 'r'},
		{"rune", '日'},
		{"empty time.Time", time.Time{}},
		{"time.Time", time.Date(2100, 12, 31, 23, 59, 59, 999999999, time.UTC)},
		{"byte", byte('b')},
		{"boolean: true", true},
		{"boolean: false", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When:
			got := ptr.Ref(tc.given)

			// Then:
			assert.Equal(t, tc.given, ptr.Value(got))
		})
	}
}

func TestRef_Int(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		given any
	}{
		{"zero", 0},
		{"min int", math.MinInt},
		{"min int8", math.MinInt8},
		{"min int16", math.MinInt16},
		{"min int32", math.MinInt32},
		{"min int64", math.MinInt64},
		{"max int", math.MaxInt},
		{"max int8", math.MaxInt8},
		{"max int16", math.MaxInt16},
		{"max int32", math.MaxInt32},
		{"max int64", math.MaxInt64},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When:
			actual := ptr.Ref(tc.given)

			// Then:
			assert.Equal(t, tc.given, ptr.Value(actual))
		})
	}
}

func TestRef_Uint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		given uint
	}{
		{"zero", 0},
		{"max uint", math.MaxUint},
		{"max uint8", math.MaxUint8},
		{"max uint16", math.MaxUint16},
		{"max uint32", math.MaxUint32},
		{"max uint64", math.MaxUint64},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given:

			// When:
			actual := ptr.Ref(tc.given)

			// Then:
			assert.Equal(t, tc.given, ptr.Value(actual))
		})
	}
}

func TestRef_Float32(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		given float32
	}{
		{"zero", 0},
		{"smallest non-zero float32", math.SmallestNonzeroFloat32},
		{"neg max float32", -math.MaxFloat32},
		{"max float32", math.MaxFloat32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given:

			// When:
			actual := ptr.Ref(tc.given)

			// Then:
			assert.InDelta(t, tc.given, ptr.Value(actual), 0)
		})
	}
}

func TestRef_Float64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		given float64
	}{
		{"zero", 0},
		{"smallest non-zero float64", math.SmallestNonzeroFloat64},
		{"neg max float64", -math.MaxFloat64},
		{"max float64", math.MaxFloat64},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.given), func(t *testing.T) {
			t.Parallel()
			// Given:

			// When:
			actual := ptr.Ref(tc.given)

			// Then:
			assert.InDelta(t, tc.given, ptr.Value(actual), 0)
		})
	}
}

func TestSame_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		a, b *string
		want bool
	}{
		{
			desc: "nil == nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			desc: `nil != ""`,
			a:    nil,
			b:    ptr.Ref(""),
			want: false,
		},
		{
			desc: "b != nil",
			a:    ptr.Ref("b"),
			b:    nil,
			want: false,
		},
		{
			desc: "a != b",
			a:    ptr.Ref("a"),
			b:    ptr.Ref("b"),
			want: false,
		},
		{
			desc: "a == a",
			a:    ptr.Ref("a"),
			b:    ptr.Ref("a"),
			want: true,
		},
		{
			desc: "a != A",
			a:    ptr.Ref("a"),
			b:    ptr.Ref("A"),
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// When:
			got := ptr.SameValue(tc.a, tc.b)

			// Then:
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSame_Bool(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		a, b *bool
		want bool
	}{
		{
			desc: "nil == nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			desc: "nil != true",
			a:    nil,
			b:    ptr.Ref(true),
			want: false,
		},
		{
			desc: "nil != false",
			a:    nil,
			b:    ptr.Ref(false),
			want: false,
		},
		{
			desc: "true != false",
			a:    ptr.Ref(true),
			b:    ptr.Ref(false),
			want: false,
		},
		{
			desc: "true == true",
			a:    ptr.Ref(true),
			b:    ptr.Ref(true),
			want: true,
		},
		{
			desc: "false == false",
			a:    ptr.Ref(false),
			b:    ptr.Ref(false),
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// When:
			got := ptr.SameValue(tc.a, tc.b)

			// Then:
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSame_Int(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		a, b *int
		want bool
	}{
		{
			desc: "nil == nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			desc: "nil != 0",
			a:    nil,
			b:    ptr.Ref(0),
			want: false,
		},
		{
			desc: "1 != 10",
			a:    ptr.Ref(1),
			b:    ptr.Ref(10),
			want: false,
		},
		{
			desc: "1 == 1",
			a:    ptr.Ref(1),
			b:    ptr.Ref(1),
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// When:
			got := ptr.SameValue(tc.a, tc.b)

			// Then:
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNilIfZero(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		assert.Nil(t, ptr.NilIfZero(0))

		intVal := 5
		assert.Equal(t, &intVal, ptr.NilIfZero(intVal))
	})

	t.Run("float", func(t *testing.T) {
		assert.Nil(t, ptr.NilIfZero(0.0))

		floatVal := 3.14
		assert.Equal(t, &floatVal, ptr.NilIfZero(floatVal))
	})

	t.Run("string", func(t *testing.T) {
		assert.Nil(t, ptr.NilIfZero(""))

		strVal := "hello"
		assert.Equal(t, &strVal, ptr.NilIfZero(strVal))
	})

	t.Run("boolean", func(t *testing.T) {
		assert.Nil(t, ptr.NilIfZero(false))

		boolVal := true
		assert.Equal(t, &boolVal, ptr.NilIfZero(boolVal))
	})

	t.Run("byte", func(t *testing.T) {
		assert.Nil(t, ptr.NilIfZero(byte(0)))

		byteVal := byte(255)
		assert.Equal(t, &byteVal, ptr.NilIfZero(byteVal))
	})

	t.Run("rune", func(t *testing.T) {
		assert.Nil(t, ptr.NilIfZero(rune(0)))

		runeVal := rune('a')
		assert.Equal(t, &runeVal, ptr.NilIfZero(runeVal))
	})

	t.Run("complex", func(t *testing.T) {
		assert.Nil(t, ptr.NilIfZero(complex64(0)))

		complexVal := complex64(1 + 2i)
		assert.Equal(t, &complexVal, ptr.NilIfZero(complexVal))
	})

	t.Run("uint", func(t *testing.T) {
		assert.Nil(t, ptr.NilIfZero(uint(0)))

		uintVal := uint(10)
		assert.Equal(t, &uintVal, ptr.NilIfZero(uintVal))
	})
}
