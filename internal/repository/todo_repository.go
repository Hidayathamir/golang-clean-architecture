package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/TodoRepository.go -pkg=mock . TodoRepository

type TodoRepository interface {
	Create(ctx context.Context, db *gorm.DB, todo *entity.Todo) error
	Update(ctx context.Context, db *gorm.DB, todo *entity.Todo) error
	Delete(ctx context.Context, db *gorm.DB, todo *entity.Todo) error
	FindByIDAndUserID(ctx context.Context, db *gorm.DB, todo *entity.Todo, id string, userID string) error
	List(ctx context.Context, db *gorm.DB, req *model.ListTodoRequest) (entity.TodoList, int64, error)
}

var _ TodoRepository = &TodoRepositoryImpl{}

type TodoRepositoryImpl struct {
	Config *viper.Viper
	Log    *logrus.Logger
}

func NewTodoRepository(cfg *viper.Viper, log *logrus.Logger) *TodoRepositoryImpl {
	return &TodoRepositoryImpl{
		Config: cfg,
		Log:    log,
	}
}

func (r *TodoRepositoryImpl) Create(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	if err := db.Create(todo).Error; err != nil {
		return errkit.AddFuncName("repository.(*TodoRepositoryImpl).Create", err)
	}
	return nil
}

func (r *TodoRepositoryImpl) Update(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	if err := db.Save(todo).Error; err != nil {
		return errkit.AddFuncName("repository.(*TodoRepositoryImpl).Update", err)
	}
	return nil
}

func (r *TodoRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	if err := db.Delete(todo).Error; err != nil {
		return errkit.AddFuncName("repository.(*TodoRepositoryImpl).Delete", err)
	}
	return nil
}

func (r *TodoRepositoryImpl) FindByIDAndUserID(ctx context.Context, db *gorm.DB, todo *entity.Todo, id string, userID string) error {
	err := db.Where("id = ? AND user_id = ?", id, userID).Take(todo).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName("repository.(*TodoRepositoryImpl).FindByIDAndUserID", err)
	}

	return nil
}

func (r *TodoRepositoryImpl) List(ctx context.Context, db *gorm.DB, req *model.ListTodoRequest) (entity.TodoList, int64, error) {
	var todos entity.TodoList
	if err := db.Scopes(r.filterTodos(req)).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Order("created_at DESC").
		Find(&todos).Error; err != nil {
		return nil, 0, errkit.AddFuncName("repository.(*TodoRepositoryImpl).List", err)
	}

	var total int64
	if err := db.Model(&entity.Todo{}).Scopes(r.filterTodos(req)).Count(&total).Error; err != nil {
		return nil, 0, errkit.AddFuncName("repository.(*TodoRepositoryImpl).List", err)
	}

	return todos, total, nil
}

func (r *TodoRepositoryImpl) filterTodos(req *model.ListTodoRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("user_id = ?", req.UserID)

		if req.Title != "" {
			title := "%" + req.Title + "%"
			tx = tx.Where("title ILIKE ?", title)
		}

		if req.IsCompleted != nil {
			tx = tx.Where("is_completed = ?", *req.IsCompleted)
		}

		return tx
	}
}
