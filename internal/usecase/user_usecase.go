package usecase

import (
	"context"
	"errors"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/pkg/errkit"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase interface {
	Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error)
	Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error)
	Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error)
	Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error)
	Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error)
}

var _ UserUseCase = &UserUseCaseImpl{}

type UserUseCaseImpl struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate

	// repository
	UserRepository repository.UserRepository

	// producer
	UserProducer messaging.UserProducer
}

func NewUserUseCase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,

	// repository
	userRepository repository.UserRepository,

	// producer
	userProducer messaging.UserProducer,
) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		DB:       db,
		Log:      logger,
		Validate: validate,

		// repository
		UserRepository: userRepository,

		// producer
		UserProducer: userProducer,
	}
}

func (u *UserUseCaseImpl) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByToken(ctx, tx, user, req.Token); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return &model.Auth{ID: user.ID}, nil
}

func (u *UserUseCaseImpl) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	total, err := u.UserRepository.CountById(ctx, tx, req.ID)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if total > 0 {
		err = errors.New("user already exists")
		err = errkit.Conflict(err)
		return nil, errkit.AddFuncName(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	user := &entity.User{
		ID:       req.ID,
		Password: string(password),
		Name:     req.Name,
	}

	if err := u.UserRepository.Create(ctx, tx, user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err = u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCaseImpl) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(ctx, tx, user, req.ID); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	user.Token = uuid.New().String()
	if err := u.UserRepository.Update(ctx, tx, user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToTokenResponse(user), nil
}

func (u *UserUseCaseImpl) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(ctx, tx, user, req.ID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCaseImpl) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return false, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(ctx, tx, user, req.ID); err != nil {
		return false, errkit.AddFuncName(err)
	}

	user.Token = ""

	if err := u.UserRepository.Update(ctx, tx, user); err != nil {
		return false, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return false, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return false, errkit.AddFuncName(err)
	}

	return true, nil
}

func (u *UserUseCaseImpl) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(ctx, tx, user, req.ID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errkit.AddFuncName(err)
		}
		user.Password = string(password)
	}

	if err := u.UserRepository.Update(ctx, tx, user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}
