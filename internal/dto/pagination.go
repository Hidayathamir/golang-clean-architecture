package dto

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/pagination"
)

const (
	// DefaultPageLimit is used when the request omits limit.
	DefaultPageLimit = 20
	// MaxPageLimit caps the page size for every list endpoint.
	MaxPageLimit = 100
)

// PageRequest is the shared cursor-pagination input for list endpoints.
// An empty Cursor means "first page".
type PageRequest struct {
	Limit  int
	Cursor string
}

// NormalizedLimit returns the request limit clamped to [1, MaxPageLimit],
// falling back to DefaultPageLimit when unset. The requested limit + 1 is
// fetched so callers can tell whether a further page exists.
func (p PageRequest) NormalizedLimit() int {
	switch {
	case p.Limit <= 0:
		return DefaultPageLimit
	case p.Limit > MaxPageLimit:
		return MaxPageLimit
	default:
		return p.Limit
	}
}

// DecodeCursor decodes the request cursor. An empty cursor is the first page
// (ok=false), not an error.
func (p PageRequest) DecodeCursor() (pagination.Cursor, bool, error) {
	return pagination.Decode(p.Cursor)
}

// CursorPage is the shared list envelope. NextCursor is empty on the last page.
type CursorPage[T any] struct {
	Data       T      `json:"data"`
	NextCursor string `json:"next_cursor"`
}

// PageFrom builds the envelope, deriving NextCursor from the last row when a
// further page exists.
func PageFrom[T any](data T, hasMore bool, lastCreatedAt time.Time, lastID int64) CursorPage[T] {
	page := CursorPage[T]{Data: data}
	if hasMore {
		page.NextCursor = pagination.Encode(pagination.Cursor{CreatedAt: lastCreatedAt, ID: lastID})
	}
	return page
}
