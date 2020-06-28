package model

import (
	"context"
	"time"
)

// Article struct
type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    Author    `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// SetCreatedAt method
func (a *Article) SetCreatedAt() {
	a.CreatedAt = time.Now()
}

// SetUpdatedAt method
func (a *Article) SetUpdatedAt() {
	a.UpdatedAt = time.Now()
}

// ArticleUsecase interface
type ArticleUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Article, string, error)
	GetByID(ctx context.Context, id int64) (Article, error)
	Update(ctx context.Context, ar *Article) error
	GetByTitle(ctx context.Context, title string) (Article, error)
	Store(ctx context.Context, ar *Article) error
	Delete(ctx context.Context, id int64) error
}

// ArticleRepository interface
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Article, string, error)
	GetByID(ctx context.Context, id int64) (Article, error)
	Update(ctx context.Context, ar *Article) error
	GetByTitle(ctx context.Context, title string) (Article, error)
	Store(ctx context.Context, ar *Article) error
	Delete(ctx context.Context, id int64) error
}
