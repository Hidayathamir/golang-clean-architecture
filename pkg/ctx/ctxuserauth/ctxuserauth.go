package ctxuserauth

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
)

type userAuthContextKey struct{}

var userAuthKey = userAuthContextKey{}

func Set(ctx context.Context, userAuth *model.UserAuth) context.Context {
	return context.WithValue(ctx, userAuthKey, userAuth)
}

func Get(ctx context.Context) *model.UserAuth {
	if val, ok := ctx.Value(userAuthKey).(*model.UserAuth); ok {
		return val
	}
	return nil
}
