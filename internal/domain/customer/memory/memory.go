package memory

import (
	"sync"

	"github.com/google/uuid"

	"github.com/dlfdyd96/ddd-go/internal/aggregate"
)

// Repository fulfills the CustomerRepository interface
type Repository struct {
	customers map[uuid.UUID]aggregate.Customer
	sync.Mutex
}

// New is a factory function to generate a new repository of customers
func New() *Repository {
	return &Repository{
		customers: make(map[uuid.UUID]aggregate.Customer),
	}
}

// Get finds a customer by ID
func (r *Repository) Get(uuid.UUID) (aggregate.Customer, error) {
	return aggregate.Customer{}, nil
}

// Add will add a new customer to the repository
func (r *Repository) Add(aggregate.Customer) error {
	return nil
}

// Update will replace an existing customer information with the new customer information
func (r *Repository) Update(aggregate.Customer) error {
	return nil
}
