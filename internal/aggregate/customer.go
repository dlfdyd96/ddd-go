package aggregate

import (
	"github.com/dlfdyd96/ddd-go/internal/entity"
	"github.com/dlfdyd96/ddd-go/internal/valueobject"
)

// Customer is a aggregate that combines all entities needed to represent a customer
type Customer struct {
	// person is the root entiy of a customer
	// which means the person.ID is the main identifier for this aggregate
	person *entity.Person
	// a customer can hold many products
	products []*entity.Item
	// a customer can perform many transaction
	// the value objects are held as nonpointers though since they cannot change state.
	transactions []valueobject.Transaction
}
