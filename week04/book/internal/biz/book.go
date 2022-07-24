package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrBookNotFound = errors.NotFound("1", "book not found")
)

type Book struct {
	Name  string
	Count int
}

type BookRepo interface {
	Save(context.Context, *Book) (*Book, error)
	Update(context.Context, *Book) (*Book, error)
	FindByName(context.Context, *Book) (*Book, error)
	ListAll(context.Context) ([]*Book, error)
}

type BookUsecase struct {
	repo BookRepo
	log  *log.Helper
}

func NewBookUsecase(repo BookRepo, logger log.Logger) *BookUsecase {
	return &BookUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *BookUsecase) CreateBook(ctx context.Context, g *Book) (*Book, error) {
	uc.log.WithContext(ctx).Infof("CreateBook: %v", g.Name)
	return uc.repo.Save(ctx, g)
}
