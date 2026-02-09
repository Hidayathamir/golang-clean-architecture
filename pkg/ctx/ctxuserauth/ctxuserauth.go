package ctxuserauth

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
)

type userAuthContextKey struct{}

var userAuthKey = userAuthContextKey{}

func Set(ctx context.Context, userAuth *dto.UserAuth) context.Context {
	return context.WithValue(ctx, userAuthKey, userAuth)
}

func Get(ctx context.Context) *dto.UserAuth {
	if val, ok := ctx.Value(userAuthKey).(*dto.UserAuth); ok {
		return val
	}
	return nil
}
