package product

import (
	"errors"
	"github.com/dlfdyd96/ddd-go/internal/aggregate"
	"github.com/google/uuid"
)

var (
	//ErrProductNotFound is returned when a product is not found
	ErrProductNotFound = errors.New("the product was not found")
	//ErrProductAlreadyExist is returned when trying to add a product that already exists
	ErrProductAlreadyExist = errors.New("the product already exists")
)

// ProductRepository is the repository interface to fulfill to use the product aggregate
type ProductRepository interface {
	GetAll() ([]aggregate.Product, error)
	GetByID(uuid.UUID) (aggregate.Product, error)
	Add(aggregate.Product) error
	Update(aggregate.Product) error
	Delete(uuid.UUID) error
}
