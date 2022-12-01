package students

import (
	"errors"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/repos/classes"
	"github.com/h4n-openschool/api/utils"
	"github.com/lucsky/cuid"
)

var (
	StudentDoesNotExist = errors.New("no existing student found by that id")
)

// InMemoryStudentRepository implements the [studentRepository] interface using an
// in-memory slice of [models.Student] items.
type InMemoryStudentRepository struct {
	// Items is the slice of [models.Student] items stored in memory.
	Items []models.Student
}

// NewInMemoryStudentRepository creates a new instance of
// [NewInMemoryStudentRepository]
func NewInMemoryStudentRepository(cr classes.ClassRepository, itemCount int) InMemoryStudentRepository {
  classes, _ := cr.GetAll(utils.NewPaginationQuery())

  var items []models.Student
  for _, c := range classes {
    var classStudents []string

    // Generate studentes in-memory to use with repo methods.
    for i := 0; i < itemCount; i++ {
      id := cuid.New()

      items = append(items, models.Student{
        BaseMetadata: models.BaseMetadata{
          Id:        id,
          CreatedAt: time.Now(),
          UpdatedAt: time.Now(),
        },
        FullName: faker.FirstName() + " " + faker.LastName(),
      })

      classStudents = append(classStudents, id)
    }

    c.StudentIds = classStudents
    _, _ = cr.Update(&c)
  }

	// Return the new repository to the caller
	return InMemoryStudentRepository{Items: items}
}

func (r *InMemoryStudentRepository) GetAll(pq utils.PaginationQuery) ([]models.Student, error) {
	var items []models.Student

	offset := pq.Offset()
	for i := offset; i < (offset + pq.PerPage); i++ {
		item := r.Items[i]
		items = append(items, item)
	}

	return items, nil
}

func (r *InMemoryStudentRepository) Get(id string) (*models.Student, error) {
	var found *models.Student

	for _, v := range r.Items {
		if v.Id == id {
			found = &v
			break
		}
	}

	return found, nil
}

func (r *InMemoryStudentRepository) Update(student *models.Student) (*models.Student, error) {
	var found *models.Student

	for k, v := range r.Items {
		if v.Id == student.Id {
      v.FullName = student.FullName
			v.UpdatedAt = time.Now()

			found = &v

			r.Items[k] = *found

			break
		}
	}

	if found == nil {
		return nil, StudentDoesNotExist
	}

	return found, nil
}

func (r *InMemoryStudentRepository) Create(student models.Student) (*models.Student, error) {
	model := models.Student{
		BaseMetadata: models.BaseMetadata{
			Id:        cuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
    FullName: student.FullName,
	}

	r.Items = append(r.Items, model)

	return &model, nil
}

func (r *InMemoryStudentRepository) Delete(student models.Student) error {
	var newItems []models.Student

	var found *models.Student
	for _, c := range r.Items {
		if c.Id == student.Id {
			found = &c
			break
		}
	}
	if found == nil {
		return StudentDoesNotExist
	}

	for _, c := range r.Items {
		if c.Id != student.Id {
			newItems = append(newItems, c)
		}
	}

	r.Items = newItems

	return nil
}

func (r *InMemoryStudentRepository) Count() (int, error) {
	return len(r.Items), nil
}
