// Package eserr provides generic error-chain utilities for working with
// Elasticsearch errors returned by the go-elasticsearch typed client.
package eserr

// UnwrapErr calls Unwrap on the error if available, returning the wrapped error.
// Returns nil if err does not implement the Unwrap method.
func UnwrapErr(err error) error {
	type unwrapper interface {
		Unwrap() error
	}
	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}
	return nil
}

// FindErrorInChain walks the error chain looking for the first error whose
// concrete type matches T. If found, target (if non-nil) is set to the matched
// value and true is returned. Returns false if no matching error is found.
//
// T should be an error pointer type (e.g. *MyError). This allows type inference
// when passing &target where target is of that pointer type.
//
// Example:
//
//	var esErr *types.ElasticsearchError
//	if eserr.FindErrorInChain(err, &esErr) {
//	    // esErr is now set
//	}
func FindErrorInChain[T any](err error, target *T) bool {
	e := err
	for e != nil {
		if te, ok := e.(T); ok {
			if target != nil {
				*target = te
			}
			return true
		}
		e = UnwrapErr(e)
	}
	return false
}
