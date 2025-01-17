package ptr

// Ref returns the pointer reference to the value in the argument.
func Ref[T any](v T) *T {
	return &v
}

// Value returns the dereference of the pointer in the argument
func Value[T any](v *T) T {
	if v == nil {
		var r T

		return r
	}

	return *v
}

// SameValue check if given pointers has the same value
func SameValue[T comparable](a, b *T) bool {
	if a != nil && b != nil {
		return Value(a) == Value(b)
	}

	return a == nil && b == nil
}

// NilIfZero returns nil if the provided value is a zero value.
//
// Otherwise, it will return the pointer to the value
func NilIfZero[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// If bool provided is true, provided value will be returned.
// Else, a zero value will be returned.
func If[T any](clause bool, v T) T {
	if clause {
		return v
	}
	var zero T
	return zero
}
