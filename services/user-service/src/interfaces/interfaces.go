package interfaces

import (
	"github.com/online-order-system/user-service/models"
)

// UserService defines the interface for user service
type UserService interface {
	CreateUser(req models.CreateCustomerRequest) (models.Customer, error)
	GetUserByID(id string) (models.Customer, error)
	GetUserByEmail(email string) (models.Customer, error)
	GetUsers() ([]models.Customer, error)
	UpdateUser(id string, req models.UpdateCustomerRequest) (models.Customer, error)
	DeleteUser(id string) error
	VerifyUser(req models.VerifyCustomerRequest) (models.VerifyCustomerResponse, error)
	GetUserOrders(userID string) ([]models.CustomerOrder, error)
	AddUserOrder(userID string, orderID string, orderStatus string) error
	UpdateUserOrderStatus(userID string, orderID string, orderStatus string) error
}

// UserProducer defines the interface for user producer
type UserProducer interface {
	PublishUserVerified(user models.Customer) error
	PublishUserUpdated(user models.Customer) error
	Close() error
}

// UserConsumer defines the interface for user consumer
type UserConsumer interface {
	ProcessOrderEvent(event models.OrderEvent) error
}
