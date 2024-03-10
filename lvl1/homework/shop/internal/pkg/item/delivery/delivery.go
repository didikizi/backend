package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"shop/internal/pkg/item"
	"shop/internal/pkg/models"
	"strconv"
)

type delivery struct {
	items item.ItemsUsecase
}

func New(items item.ItemsUsecase) item.RESTDelivery {
	return delivery{
		items: items,
	}
}

func (d delivery) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var response interface{}
	var err error
	switch r.Method {
	case http.MethodGet:
		response, err = d.handleGet(ctx, r)
	case http.MethodPost:
		response, err = d.handlePost(ctx, r)
	case http.MethodPut:
		response, err = d.handlePut(ctx, r)
	case http.MethodDelete:
		response, err = d.handleDelete(ctx, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		log.Println(err)
	}
}

func (d delivery) handleGet(ctx context.Context, r *http.Request) (interface{}, error) {
	filter := models.ItemFilter{}
	max := r.FormValue("price_max")
	min := r.FormValue("price_min")

	if max != "" {
		val, err := strconv.Atoi(max)
		if err != nil {
			return nil, err
		}
		filter.PriceMax = val
	}

	if min != "" {
		val, err := strconv.Atoi(min)
		if err != nil {
			return nil, err
		}
		filter.PriceMin = val
	}

	return d.items.List(ctx, filter)
}

func (d delivery) handlePost(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (d delivery) handlePut(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (d delivery) handleDelete(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, errors.New("not implemented")
}
