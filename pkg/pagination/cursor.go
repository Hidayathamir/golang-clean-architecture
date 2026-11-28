// Package pagination implements an opaque cursor used by every list endpoint.
//
// A cursor encodes the sort position of the last row on the previous page
// (created_at + id) so that rows inserted between two requests never cause
// duplicates or skips. Ordering must always tie-break by id.
package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

// ErrInvalidCursor is returned when a cursor cannot be decoded.
var ErrInvalidCursor = errors.New("invalid cursor")

// Cursor is the decoded sort position of the last row of a page.
type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
}

// Encode serializes a cursor into an opaque, URL-safe string.
func Encode(c Cursor) string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// Decode parses an opaque cursor string. An empty string means "first page"
// and returns a zero cursor with ok=false.
func Decode(s string) (Cursor, bool, error) {
	if s == "" {
		return Cursor{}, false, nil
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return Cursor{}, false, ErrInvalidCursor
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return Cursor{}, false, ErrInvalidCursor
	}
	return c, true, nil
}
