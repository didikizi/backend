package item

import (
	"context"
	"errors"
	"net/http"

	"shop/internal/pkg/models"
)

type ItemsUsecase interface {
	Create(ctx context.Context, item *models.Item) error
	List(ctx context.Context, filter models.ItemFilter) ([]models.Item, error)
	Update(ctx context.Context, item models.Item) error
	Delete(ctx context.Context, id int) error
}

type RESTDelivery interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

var (
	ErrItemNotFound = errors.New("item not found")
)
