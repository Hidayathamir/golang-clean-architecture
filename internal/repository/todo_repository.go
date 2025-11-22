package repository

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/TodoRepository.go -pkg=mock . TodoRepository

type TodoRepository interface {
	Create(ctx context.Context, db *gorm.DB, todo *entity.Todo) error
	Update(ctx context.Context, db *gorm.DB, todo *entity.Todo) error
	Delete(ctx context.Context, db *gorm.DB, todo *entity.Todo) error
	FindByIDAndUserID(ctx context.Context, db *gorm.DB, todo *entity.Todo, id int64, userID int64) error
	List(ctx context.Context, db *gorm.DB, req *model.ListTodoRequest) (entity.TodoList, int64, error)
}

var _ TodoRepository = &TodoRepositoryImpl{}

type TodoRepositoryImpl struct {
	Config *viper.Viper
}

func NewTodoRepository(cfg *viper.Viper) *TodoRepositoryImpl {
	return &TodoRepositoryImpl{
		Config: cfg,
	}
}

func (r *TodoRepositoryImpl) Create(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	if err := db.WithContext(ctx).Create(todo).Error; err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *TodoRepositoryImpl) Update(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	if err := db.WithContext(ctx).Save(todo).Error; err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *TodoRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	if err := db.WithContext(ctx).Delete(todo).Error; err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (r *TodoRepositoryImpl) FindByIDAndUserID(ctx context.Context, db *gorm.DB, todo *entity.Todo, id int64, userID int64) error {
	err := db.WithContext(ctx).
		Where(map[string]any{
			entity.TodoColumnID:     id,
			entity.TodoColumnUserID: userID,
		}).
		Take(todo).Error
	if err != nil {
		err = errkit.NotFound(err)
		return errkit.AddFuncName(err)
	}

	return nil
}

func (r *TodoRepositoryImpl) List(ctx context.Context, db *gorm.DB, req *model.ListTodoRequest) (entity.TodoList, int64, error) {
	var todos entity.TodoList
	if err := db.WithContext(ctx).Scopes(r.filterTodos(req)).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Order(fmt.Sprintf("%s DESC", entity.TodoColumnCreatedAt)).
		Find(&todos).Error; err != nil {
		return nil, 0, errkit.AddFuncName(err)
	}

	var total int64
	if err := db.Model(&entity.Todo{}).Scopes(r.filterTodos(req)).Count(&total).Error; err != nil {
		return nil, 0, errkit.AddFuncName(err)
	}

	return todos, total, nil
}

func (r *TodoRepositoryImpl) filterTodos(req *model.ListTodoRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where(map[string]any{entity.TodoColumnUserID: req.UserID})

		if req.Title != "" {
			title := "%" + req.Title + "%"
			tx = tx.Where(fmt.Sprintf("%s ILIKE ?", entity.TodoColumnTitle), title)
		}

		if req.IsCompleted != nil {
			tx = tx.Where(fmt.Sprintf("%s = ?", entity.TodoColumnIsCompleted), *req.IsCompleted)
		}

		return tx
	}
}
