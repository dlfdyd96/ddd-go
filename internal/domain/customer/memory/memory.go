package memory

import (
	"fmt"
	"github.com/dlfdyd96/ddd-go/internal/domain/customer"
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
func (r *Repository) Get(id uuid.UUID) (aggregate.Customer, error) {
	if customer, ok := r.customers[id]; ok {
		return customer, nil
	}
	return aggregate.Customer{}, customer.ErrCustomerNotFound
}

// Add will add a new customer to the repository
func (r *Repository) Add(c aggregate.Customer) error {
	if r.customers == nil {
		// Saftey check if customers is not create, shouldn't happen if using the Factory, but you never know
		r.Lock()
		r.customers = make(map[uuid.UUID]aggregate.Customer)
		r.Unlock()
	}

	// Make sure Customer isn't already in the repository
	if _, ok := r.customers[c.GetID()]; ok {
		return fmt.Errorf("customer already exists: %w", customer.ErrFailedToAddCustomer)
	}
	r.Lock()
	r.customers[c.GetID()] = c
	r.Unlock()
	return nil
}

// Update will replace an existing customer information with the new customer information
func (r *Repository) Update(c aggregate.Customer) error {
	// Make sure Customer is in the repository
	if _, ok := r.customers[c.GetID()]; !ok {
		return fmt.Errorf("customer does not exist: %w", customer.ErrUpdateCustomer)
	}
	r.Lock()
	r.customers[c.GetID()] = c
	r.Unlock()
	return nil
}
