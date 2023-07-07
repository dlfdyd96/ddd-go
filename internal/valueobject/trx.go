package valueobject

import (
	"github.com/google/uuid"
	"time"
)

// Transaction is a payment between two parties
// strong consist
type Transaction struct {
	// all values lowercase since they are immutable
	amount    int
	from      uuid.UUID
	to        uuid.UUID
	createdAt time.Time
}
