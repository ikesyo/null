// Package null provides an opinionated yet reasonable way of handling null values.
package null

import (
	"database/sql"
	"encoding/json"
)

// String is an even nuller nullable string.
type String struct {
	sql.NullString
}

// StringFrom creates a new String that will never be blank.
func StringFrom(s string) String {
	return NewString(s, true)
}

// StringFrom creates a new String that be null if s is nil.
func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	str := NewString(*s, true)
	return str
}

// NewString creates a new String
func NewString(s string, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input produces a null String.
// It also supports unmarshalling a sql.NullString.
func (s *String) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	json.Unmarshal(data, &v)
	switch x := v.(type) {
	case string:
		s.String = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &s.NullString)
	case nil:
		s.Valid = false
		return nil
	}
	s.Valid = (err == nil) && (s.String != "")
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this String is null.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (s *String) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null or empty strings, for future omitempty support. (Go 1.4?)
// Will return false s is blank but non-null.
func (s String) IsZero() bool {
	return !s.Valid
}
