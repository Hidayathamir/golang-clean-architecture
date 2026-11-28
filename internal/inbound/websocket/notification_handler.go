package websocket

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

// TokenVerifier is the narrow slice of the user usecase the websocket handler
// needs. Declaring it here (rather than importing userusecase) keeps the
// notifusecase -> websocket -> userusecase -> notifusecase edge from forming a
// cycle.
type TokenVerifier interface {
	Verify(ctx context.Context, req dto.VerifyUserRequest) (dto.UserAuth, error)
}

type NotificationHandler struct {
	Hub           *Hub
	TokenVerifier TokenVerifier
}

func NewNotificationHandler(hub *Hub, tokenVerifier TokenVerifier) *NotificationHandler {
	return &NotificationHandler{
		Hub:           hub,
		TokenVerifier: tokenVerifier,
	}
}

// Handle authenticates the connection from the token query parameter and then
// runs the read/write loop for its lifetime. Browsers cannot set headers on a
// WebSocket handshake, hence the query parameter.
func (h *NotificationHandler) Handle(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	userAuth, err := h.TokenVerifier.Verify(ctx.UserContext(), dto.VerifyUserRequest{Token: token})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	userID := userAuth.ID

	return websocket.New(func(conn *websocket.Conn) {
		client := NewClient(userID)
		h.Hub.Register(userID, client)
		defer func() {
			h.Hub.Unregister(userID, client)
			client.Close()
			_ = conn.Close()
		}()

		// Writer: pump hub messages to the socket.
		go func() {
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()
			for {
				select {
				case payload, ok := <-client.Send:
					_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
					if !ok {
						_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
						return
					}
					if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
						return
					}
				case <-ticker.C:
					_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
						return
					}
				}
			}
		}()

		// Reader: keep the connection alive and detect disconnects.
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(pongWait))
		})
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	})(ctx)
}
