package ordersUsecases

import (
	"fmt"
	"math"

	"github.com/IzePhanthakarn/kawaii-shop/modules/entities"
	"github.com/IzePhanthakarn/kawaii-shop/modules/orders"
	"github.com/IzePhanthakarn/kawaii-shop/modules/orders/ordersRepositories"
	"github.com/IzePhanthakarn/kawaii-shop/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrder(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
}

type ordersUsecase struct {
	ordersRepository     ordersRepositories.IOrdersRepository
	productsRepositories productsRepositories.IProductsRepository
}

func OrderUsecase(ordersRepository ordersRepositories.IOrdersRepository, productsRepositories productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepository:     ordersRepository,
		productsRepositories: productsRepositories,
	}
}

func (u *ordersUsecase) FindOneOrder(orderId string) (*orders.Order, error) {
	order, err := u.ordersRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (u *ordersUsecase) FindOrder(req *orders.OrderFilter) *entities.PaginateRes {
	orders, count := u.ordersRepository.FindOrder(req)
	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *ordersUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {
	// Check product is exist
	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product not nil")
		}

		product, err := u.productsRepositories.FindOneProduct(req.Products[i].Product.Id)
		if err != nil {
			return nil, err
		}
		
		// Set price
		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = product
	}
	
	orderId, err := u.ordersRepository.InsertOrder(req)
	if err != nil {
		return nil, err
	}
	
	order, err := u.ordersRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}
