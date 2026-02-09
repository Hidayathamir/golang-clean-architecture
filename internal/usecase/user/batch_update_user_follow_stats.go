package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) BatchUpdateUserFollowStats(ctx context.Context, req *dto.BatchUpdateUserFollowStatsRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	for _, v := range req.UserIncreaseFollowerFollowingCountList {
		var err error = nil
		switch {
		case v.HasFollowerCountAndFollowingCount():
			err = u.UserStatRepository.IncrementFollowerCountAndFollowingCountByID(ctx, u.DB, v.UserID, v.FollowerCount, v.FollowingCount)
		case v.HasFollowerCount():
			err = u.UserStatRepository.IncrementFollowerCountByID(ctx, u.DB, v.UserID, v.FollowerCount)
		case v.HasFollowingCount():
			err = u.UserStatRepository.IncrementFollowingCountByID(ctx, u.DB, v.UserID, v.FollowingCount)
		default:
			x.Logger.WithContext(ctx).WithField("v", v).Warn("invalid follower count and following count")
		}
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn()
		}
	}

	return nil
}
