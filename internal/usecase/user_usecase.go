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
		u.Log.Warnf("Invalid request body : %+v", err)
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByToken(tx, user, req.Token); err != nil {
		u.Log.Warnf("Failed find user by token : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	return &model.Auth{ID: user.ID}, nil
}

func (u *UserUseCaseImpl) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(req)
	if err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	total, err := u.UserRepository.CountById(tx, req.ID)
	if err != nil {
		u.Log.Warnf("Failed count user from database : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	if total > 0 {
		err = errors.New("user already exists")
		u.Log.Warnf("User already exists : %+v", err)
		err = errkit.Conflict(err)
		return nil, errkit.AddFuncName(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	user := &entity.User{
		ID:       req.ID,
		Password: string(password),
		Name:     req.Name,
	}

	if err := u.UserRepository.Create(tx, user); err != nil {
		u.Log.Warnf("Failed create user to database : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err = u.UserProducer.Send(event); err != nil {
		u.Log.WithError(err).Error("failed to publish user created event")
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCaseImpl) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.Warnf("Invalid request body  : %+v", err)
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, req.ID); err != nil {
		u.Log.Warnf("Failed find user by id : %+v", err)
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		u.Log.Warnf("Failed to compare user password with bcrype hash : %+v", err)
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName(err)
	}

	user.Token = uuid.New().String()
	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed save user : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(event); err != nil {
		u.Log.WithError(err).Error("Failed publish user login event")
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToTokenResponse(user), nil
}

func (u *UserUseCaseImpl) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, req.ID); err != nil {
		u.Log.Warnf("Failed find user by id : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCaseImpl) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		err = errkit.BadRequest(err)
		return false, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, req.ID); err != nil {
		u.Log.Warnf("Failed find user by id : %+v", err)
		return false, errkit.AddFuncName(err)
	}

	user.Token = ""

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed save user : %+v", err)
		return false, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return false, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(event); err != nil {
		u.Log.WithError(err).Error("Failed publish user logout event")
		return false, errkit.AddFuncName(err)
	}

	return true, nil
}

func (u *UserUseCaseImpl) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, req.ID); err != nil {
		u.Log.Warnf("Failed find user by id : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			u.Log.Warnf("Failed to generate bcrype hash : %+v", err)
			return nil, errkit.AddFuncName(err)
		}
		user.Password = string(password)
	}

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.Warnf("Failed save user : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(event); err != nil {
		u.Log.WithError(err).Error("Failed publish user updated event")
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}
