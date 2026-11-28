package websocket_test

import (
	"testing"
	"time"

	ws "github.com/Hidayathamir/golang-clean-architecture/internal/inbound/websocket"
	"github.com/stretchr/testify/require"
)

func TestHub_Send_DeliversToRegisteredClient(t *testing.T) {
	hub := ws.NewHub()
	client := ws.NewClient(7)
	hub.Register(7, client)

	hub.Send(7, []byte(`{"type":"follow"}`))

	select {
	case got := <-client.Send:
		require.Equal(t, `{"type":"follow"}`, string(got))
	case <-time.After(time.Second):
		t.Fatal("expected message, got none")
	}
}

func TestHub_Send_NoopForUnknownUser(t *testing.T) {
	hub := ws.NewHub()
	require.NotPanics(t, func() {
		hub.Send(99, []byte(`{}`))
	})
}

func TestHub_Unregister_StopsDelivery(t *testing.T) {
	hub := ws.NewHub()
	client := ws.NewClient(7)
	hub.Register(7, client)
	hub.Unregister(7, client)

	require.Equal(t, 0, hub.ConnectionCount(7))
	hub.Send(7, []byte(`{}`))
	select {
	case <-client.Send:
		t.Fatal("expected no message after unregister")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestHub_Send_DoesNotBlockOnSlowClient(t *testing.T) {
	hub := ws.NewHub()
	client := ws.NewClient(7)
	hub.Register(7, client)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			hub.Send(7, []byte(`{}`))
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Send blocked on a slow client")
	}
}
