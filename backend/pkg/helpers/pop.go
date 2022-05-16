package helpers

import "github.com/gobuffalo/nulls"

// NullStringify returns a NullString if the given string is empty
func NullStringify(s string) nulls.String {
	if s != "" {
		return nulls.NewString(s)
	}
	return nulls.String{}
}
