package websocket

import "sync"

const clientSendBuffer = 64

type Client struct {
	UserID int64
	Send   chan []byte
	once   sync.Once
}

func NewClient(userID int64) *Client {
	return &Client{
		UserID: userID,
		Send:   make(chan []byte, clientSendBuffer),
	}
}

// Close makes the client's send channel safe to abandon; it is idempotent so
// both the reader and writer loop can call it.
func (c *Client) Close() {
	c.once.Do(func() {
		close(c.Send)
	})
}
