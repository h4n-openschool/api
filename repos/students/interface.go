package students

import (
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/utils"
)

// StudentRepository defines a common interface for querying Student data
type StudentRepository interface {
	// GetAll returns all of the available Student items, paginated based on the
	// passed arguments.
	GetAll(pq utils.PaginationQuery) ([]models.Student, error)

	// Get returns a single Student by its ID.
	Get(id string) (*models.Student, error)

	// Update takes a Student object that has been mutated and persists it to the
	// data store, returning the modified object and possibly an error.
	Update(Student *models.Student) (*models.Student, error)

	// Create takes a Student object that has been populated with data and creates
	// a record for it in the data store, returning the filled record and
	// possibly an error.
	Create(Student models.Student) (*models.Student, error)

	// Delete takes a Student object that includes at least an ID and deletes the
	// relevant record for it in the data store.
	Delete(Student models.Student) error

	// Count returns the total number of students in the datastore.
	Count() (int, error)
}
