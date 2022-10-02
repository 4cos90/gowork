package data

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"

	ctV1 "github.com/go-kratos/beer-shop/api/catalog/service/v1"
	"github.com/go-kratos/beer-shop/app/shop/interface/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

var _ biz.CatalogRepo = (*catalogRepo)(nil)

type catalogRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewBeerRepo(data *Data, logger log.Logger) biz.CatalogRepo {
	return &catalogRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/beer")),
		sg:   &singleflight.Group{},
	}
}

func (r *catalogRepo) GetBeer(ctx context.Context, id int64) (*biz.Beer, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_beer_by_id_%d", id), func() (interface{}, error) {
		reply, err := r.data.bc.GetBeer(ctx, &ctV1.GetBeerReq{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		images := make([]biz.Image, 0)
		for _, x := range reply.Image {
			images = append(images, biz.Image{URL: x.Url})
		}
		return &biz.Beer{
			Id:          reply.Id,
			Name:        reply.Name,
			Description: reply.Description,
			Count:       reply.Count,
			Images:      images,
		}, err
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.Beer), nil
}

func (r *catalogRepo) ListBeer(ctx context.Context, pageNum, pageSize int64) ([]*biz.Beer, error) {
	reply, err := r.data.bc.ListBeer(ctx, &ctV1.ListBeerReq{
		PageNum:  pageNum,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.Beer, 0)
	for _, x := range reply.Results {
		images := make([]biz.Image, 0)
		for _, img := range x.Image {
			images = append(images, biz.Image{URL: img.Url})
		}
		rv = append(rv, &biz.Beer{
			Id:          x.Id,
			Description: x.Description,
			Count:       x.Count,
			Images:      images,
		})
	}
	return rv, err
}
