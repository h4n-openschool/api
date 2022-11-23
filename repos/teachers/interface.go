package teachers

import (
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/utils"
)

// TeacherRepository defines a common interface for querying Teacher data
type TeacherRepository interface {
	// GetAll returns all of the available teacher items, paginated based on the
	// passed arguments.
	GetAll(pq utils.PaginationQuery) ([]models.Teacher, error)

	// Get returns a single teacher by its ID.
	Get(id string) (*models.Teacher, error)

	// Update takes a teacher object that has been mutated and persists it to the
	// data store, returning the modified object and possibly an error.
	Update(teacher *models.Teacher) (*models.Teacher, error)

	// Create takes a teacher object that has been populated with data and creates
	// a record for it in the data store, returning the filled record and
	// possibly an error.
	Create(teacher models.Teacher) (*models.Teacher, error)

	// Delete takes a teacher object that includes at least an ID and deletes the
	// relevant record for it in the data store.
	Delete(teacher models.Teacher) error

	// Count returns the total number of teachers in the datastore.
	Count() (int, error)
}
