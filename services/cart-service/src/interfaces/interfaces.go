package interfaces

import (
	"github.com/online-order-system/cart-service/models"
)

// CartService defines the interface for cart service
type CartService interface {
	CreateCart(req models.CreateCartRequest) (models.Cart, error)
	GetCartByID(id string) (models.Cart, error)
	GetCartByUserID(userID string) (models.Cart, error)
	AddCartItem(cartID string, req models.AddCartItemRequest) (models.Cart, error)
	UpdateCartItem(cartID string, itemID string, req models.UpdateCartItemRequest) (models.Cart, error)
	RemoveCartItem(cartID string, itemID string) (models.Cart, error)
	DeleteCart(id string) error
	DeleteCartByUserID(userID string) error
}

// CartProducer defines the interface for cart producer
type CartProducer interface {
	PublishCartUpdated(cart models.Cart) error
	PublishCartCleared(cartID string, customerID string) error
	Close() error
}

// CartConsumer defines the interface for cart consumer
type CartConsumer interface {
	ProcessOrderEvent(event models.OrderEvent) error
}
