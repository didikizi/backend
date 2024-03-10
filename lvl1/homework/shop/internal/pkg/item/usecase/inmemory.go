package usecase

import (
	"context"
	"errors"

	"shop/internal/pkg/item"
	"shop/internal/pkg/models"
)

type inmemory struct {
	iterator int
	items    []models.Item
}

func (in *inmemory) Create(ctx context.Context, item *models.Item) error {
	item.ID = in.iterator
	in.iterator++

	in.items = append(in.items, *item)

	return nil
}

func (in *inmemory) List(ctx context.Context, filter models.ItemFilter) ([]models.Item, error) {
	res := make([]models.Item, 0, in.iterator)
	for _, val := range in.items {

		if filter.ID != 0 && filter.ID != val.ID {
			continue
		}

		if filter.PriceMax != 0 && filter.PriceMax < val.Price {
			continue
		}

		if filter.PriceMin != 0 && filter.PriceMin > val.Price {
			continue
		}

		res = append(res, val)
	}

	return res, nil
}

func (in *inmemory) Update(ctx context.Context, item models.Item) error {
	found := false
	for i, itm := range in.items {
		if itm.ID == item.ID {
			in.items = append(append(in.items[:i], item), in.items[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		errors.New("not found")
	}

	return nil
}

func (in *inmemory) Delete(ctx context.Context, id int) error {
	found := false
	for i, itm := range in.items {
		if itm.ID == id {
			in.items = append(in.items[:i], in.items[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		errors.New("not found")
	}

	return nil
}

func NewInmemory() item.ItemsUsecase {
	return &inmemory{
		items: []models.Item{
			{
				ID:          1,
				Name:        "Snack",
				Description: "Just some cnack",
				Price:       100,
			},
		},
		iterator: 2,
	}
}
