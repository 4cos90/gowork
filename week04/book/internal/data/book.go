package data

import (
	"context"

	"book/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type bookRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewBookRepo(data *Data, logger log.Logger) biz.BookRepo {
	return &bookRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *bookRepo) Save(ctx context.Context, g *biz.Book) (*biz.Book, error) {
	g.Count = g.Count + 1
	return g, nil
}

func (r *bookRepo) Update(ctx context.Context, g *biz.Book) (*biz.Book, error) {
	return g, nil
}

func (r *bookRepo) FindByName(ctx context.Context, g *biz.Book) (*biz.Book, error) {
	return g, nil
}

func (r *bookRepo) ListAll(context.Context) ([]*biz.Book, error) {
	return nil, nil
}
