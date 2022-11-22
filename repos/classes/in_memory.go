package classes

import (
	"errors"
	"fmt"
	"time"

	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/utils"
	"github.com/lucsky/cuid"
)

// InMemoryClassRepository implements the [ClassRepository] interface using an
// in-memory slice of [models.Class] items.
type InMemoryClassRepository struct {
  // Items is the slice of [models.Class] items stored in memory.
	Items []models.Class
}

// NewInMemoryClassRepository creates a new instance of
// [NewInMemoryClassRepository]
func NewInMemoryClassRepository(itemCount int) InMemoryClassRepository {
	var items []models.Class

  // Generate classes in-memory to use with repo methods.
	for i := 0; i < itemCount; i++ {
		id := cuid.New()

    desc := fmt.Sprintf(`This is class %v`, i)
		items = append(items, models.Class{
			Id:          id,
			Name:        fmt.Sprintf(`class-%v`, i),
			DisplayName: fmt.Sprintf(`Class %v`, i),
			Description: &desc,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

  // Return the new repository to the caller
  return InMemoryClassRepository{Items: items}
}

func (r *InMemoryClassRepository) GetAll(pq utils.PaginationQuery) ([]models.Class, error) {
	var items []models.Class

	offset := pq.Offset()
	for i := offset; i < (offset + pq.PerPage); i++ {
		item := r.Items[i]
		items = append(items, item)
	}

	return items, nil
}

func (r *InMemoryClassRepository) Get(id string) (*models.Class, error) {
	var found *models.Class

	for _, v := range r.Items {
		if v.Id == id {
			found = &v
			break
		}
	}

	return found, nil
}

func (r *InMemoryClassRepository) Update(class *models.Class) (*models.Class, error) {
	var found *models.Class
	for _, v := range r.Items {
		if v.Id == found.Id {
			if class.Name != v.Name {
				return nil, errors.New("you cannot update Name after creation")
			}

			v.DisplayName = class.DisplayName
			v.Description = class.Description
			v.UpdatedAt = time.Now()

			found = &v
		}
	}

	if found == nil {
		return nil, errors.New("no existing class found by that id")
	}

	return found, nil
}

func (r *InMemoryClassRepository) Create(class models.Class) (*models.Class, error) {
	model := models.Class{
		Id:          cuid.New(),
		Name:        class.Name,
		DisplayName: class.DisplayName,
		Description: class.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	r.Items = append(r.Items, model)

	return &model, nil
}

func (r *InMemoryClassRepository) Delete(class models.Class) error {
	var newItems []models.Class

	for _, c := range r.Items {
		if c.Id != class.Id {
			newItems = append(newItems, c)
		}
	}

	r.Items = newItems

	return nil
}

func (r *InMemoryClassRepository) Count() (int, error) {
  return len(r.Items), nil
}

