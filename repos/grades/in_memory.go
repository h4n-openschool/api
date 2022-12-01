package grades

import (
	"errors"
	"math/rand"
	"time"

	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/repos/classes"
	"github.com/h4n-openschool/api/utils"
	"github.com/lucsky/cuid"
)

var (
	GradeDoesNotExist = errors.New("no existing grade found by that id")
)

// InMemoryGradeRepository implements the [GradeRepository] interface using an
// in-memory slice of [models.Grade] items.
type InMemoryGradeRepository struct {
	// Items is the slice of [models.Grade] items stored in memory.
	Items []models.Grade
}

// NewInMemoryGradeRepository creates a new instance of
// [NewInMemoryGradeRepository]
func NewInMemoryGradeRepository(cr classes.ClassRepository) InMemoryGradeRepository {
	var items []models.Grade

  pq := utils.NewPaginationQuery()
  c, _ := cr.GetAll(pq)

  for _, class := range c {
    for _, stu := range class.StudentIds {
      // Generate grades in-memory to use with repo methods.
      id := cuid.New()

      grade := rand.Intn(9)

      items = append(items, models.Grade{
        BaseMetadata: models.BaseMetadata{
          Id:        id,
          CreatedAt: time.Now(),
          UpdatedAt: time.Now(),
        },
        StudentId: stu,
        Value: grade,
      })
    }
  }

	// Return the new repository to the caller
	return InMemoryGradeRepository{Items: items}
}

func (r *InMemoryGradeRepository) GetAll(classId string, pq utils.PaginationQuery) ([]models.Grade, error) {
	var items []models.Grade

	offset := pq.Offset()
	for i := offset; i < (offset + pq.PerPage); i++ {
    if r.Items[i].ClassId == classId {
      item := r.Items[i]
      items = append(items, item)
    }
	}

	return items, nil
}

func (r *InMemoryGradeRepository) Get(id string) (*models.Grade, error) {
	var found *models.Grade

	for _, v := range r.Items {
		if v.Id == id {
			found = &v
			break
		}
	}

	return found, nil
}

func (r *InMemoryGradeRepository) Update(grade *models.Grade) (*models.Grade, error) {
	var found *models.Grade

	for k, v := range r.Items {
		if v.Id == grade.Id {
			if grade.StudentId != "" && grade.StudentId != v.StudentId {
				return nil, errors.New("you cannot update StudentId after creation")
			}

      v.Value = grade.Value

			v.UpdatedAt = time.Now()

			found = &v

			r.Items[k] = *found

			break
		}
	}

	if found == nil {
		return nil, GradeDoesNotExist
	}

	return found, nil
}

func (r *InMemoryGradeRepository) Create(grade models.Grade) (*models.Grade, error) {
	model := models.Grade{
		BaseMetadata: models.BaseMetadata{
			Id:        cuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		StudentId: grade.StudentId,
		Value:     grade.Value,
	}

	r.Items = append(r.Items, model)

	return &model, nil
}

func (r *InMemoryGradeRepository) Delete(grade models.Grade) error {
	var newItems []models.Grade

	var found *models.Grade
	for _, c := range r.Items {
		if c.Id == grade.Id {
			found = &c
			break
		}
	}
	if found == nil {
		return GradeDoesNotExist
	}

	for _, c := range r.Items {
		if c.Id != grade.Id {
			newItems = append(newItems, c)
		}
	}

	r.Items = newItems

	return nil
}

func (r *InMemoryGradeRepository) Count() (int, error) {
	return len(r.Items), nil
}
