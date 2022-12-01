package grades

import (
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/utils"
)

// GradeRepository defines a common interface for querying Grade data
type GradeRepository interface {
	// GetAll returns all of the available Grade items, paginated based on the
	// passed arguments.
	GetAll(classId string, pq utils.PaginationQuery) ([]models.Grade, error)

	// Get returns a single Grade by its ID.
	Get(id string) (*models.Grade, error)

	// Update takes a grade object that has been mutated and persists it to the
	// data store, returning the modified object and possibly an error.
	Update(grade *models.Grade) (*models.Grade, error)

	// Create takes a grade object that has been populated with data and creates
	// a record for it in the data store, returning the filled record and
	// possibly an error.
	Create(grade models.Grade) (*models.Grade, error)

	// Delete takes a grade object that includes at least an ID and deletes the
	// relevant record for it in the data store.
	Delete(grade models.Grade) error

	// Count returns the total number of grades in the datastore.
	Count() (int, error)
}
