package isbn

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

type ISBN [13]byte

func (isbn ISBN) String() string {
	return string(isbn[:])
}

var (
	defaultPrefix = "978"
)

var (
	ErrTodo   = fmt.Errorf("todo")
	ErrValue  = fmt.Errorf("invalid ISBN value")
	ErrFormat = fmt.Errorf("invalid ISBN format")
)

type invalidLengthError struct{ len int }

func (err invalidLengthError) Error() string {
	return fmt.Sprintf("invalid ISBN length %d", err.len)
}

type invalidTypeError struct{ value reflect.Type }

func (err invalidTypeError) Error() string {
	return fmt.Sprintf("failed to scan type %+v for value", err.value)
}

// Parse s into an ISBN 13 or returns an error. Supports the forms
// ISBN 13 XXXXXXXXXXXXX (XXX-X-XXXX-XXXX-X) and
// ISBN 10 XXXXXXXXXX (X-XXXX-XXXX-X)
func Parse(s string) (isbn ISBN, err error) {
	switch len(s) {
	case 13 + 4: //XXX-X-XXXX-XXXX-X
	case 10: //XXXXXXXXXX
		s = defaultPrefix + s
		fallthrough
	case 13: //XXXXXXXXXXXXX or XXX-XX-XXX-XXXX-X
		if !strings.Contains(s, "-") {
			return check13(s)
		}
		s = defaultPrefix + "-" + s
	default:
		return isbn, invalidLengthError{len(s)}
	}
	// s is now in the format XXX-X-XXXX-XXXX-X
	return check13(strings.ReplaceAll(s, "-", ""))
}

// ParseBytes is like Parse, except it parses a byte slice instead of a string.
func ParseBytes(b []byte) (isbn ISBN, err error) { return Parse(string(b)) }

func check13(s string) (isbn ISBN, err error) {
	var acc [2]int

	for i := 0; i < len(s); i++ {
		switch v := int(s[i] - '0'); {
		case v >= 10:
			return isbn, ErrFormat
		default:
			acc[i%2] += v
			isbn[i] = s[i]
		}
	}

	if (acc[0]+acc[1]*3)%10 != 0 {
		return isbn, ErrValue
	}
	return isbn, nil
}

func (isbn ISBN) MarshalJSON() ([]byte, error) {
	return isbn[:], nil
}

func (isbn *ISBN) UnmarshalJSON(b []byte) (err error) {
	*isbn, err = ParseBytes(b[1 : len(b)-1])
	return err
}

// into SQL
func (isbn ISBN) Value() (driver.Value, error) {
	return isbn[:], nil
}

// from SQL
func (isbn *ISBN) Scan(v any) (err error) {
	switch u := v.(type) {
	case []byte: // sqlite3
		*isbn, err = ParseBytes(u)
	case string: // driver.Value is []byte therefore this won't run
		*isbn, err = Parse(u)
	default:
		return invalidTypeError{reflect.TypeOf(v)}
	}
	return err
}

// Support for in future implentations would be nice
// type ISBN struct {
// 	group
//  publisher
//  title
//  check digit
// }

// Useful
// wiki https://en.wikipedia.org/wiki/ISBN
