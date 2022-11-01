package repos

import "github.com/h4n-openschool/classes/models"

// ClassRepository defines a common interface for querying Class data
type ClassRepository interface {
	// GetAll returns all of the available Class items, paginated based on the
	// passed arguments.
	GetAll(perPage int, page int) ([]models.Class, error)

	// Get returns a single Class by its ID.
	Get(id string) (*models.Class, error)

	// Update takes a class object that has been mutated and persists it to the
	// data store, returning the modified object and possibly an error.
	Update(class *models.Class) (*models.Class, error)

	// Create takes a class object that has been populated with data and creates
	// a record for it in the data store, returning the filled record and
	// possibly an error.
	Create(class models.Class) (*models.Class, error)

	// Delete takes a class object that includes at least an ID and deletes the
	// relevant record for it in the data store.
	Delete(class models.Class) error
}
