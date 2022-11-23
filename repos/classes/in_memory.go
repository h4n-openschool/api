package classes

import (
	"errors"
	"fmt"
	"time"

	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/utils"
	"github.com/lucsky/cuid"
)

var (
	ClassDoesNotExist = errors.New("no existing class found by that id")
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
			BaseMetadata: models.BaseMetadata{
				Id:        id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        fmt.Sprintf(`class-%v`, i),
			DisplayName: fmt.Sprintf(`Class %v`, i),
			Description: &desc,
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

	for k, v := range r.Items {
		if v.Id == class.Id {
			if class.Name != "" && class.Name != v.Name {
				return nil, errors.New("you cannot update Name after creation")
			}

			if v.DisplayName != "" {
				v.DisplayName = class.DisplayName
			}

			if v.Description != nil {
				v.Description = class.Description
			}

			v.UpdatedAt = time.Now()

			found = &v

			r.Items[k] = *found

			break
		}
	}

	if found == nil {
		return nil, ClassDoesNotExist
	}

	return found, nil
}

func (r *InMemoryClassRepository) Create(class models.Class) (*models.Class, error) {
	model := models.Class{
		BaseMetadata: models.BaseMetadata{
			Id:        cuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        class.Name,
		DisplayName: class.DisplayName,
		Description: class.Description,
	}

	r.Items = append(r.Items, model)

	return &model, nil
}

func (r *InMemoryClassRepository) Delete(class models.Class) error {
	var newItems []models.Class

	var found *models.Class
	for _, c := range r.Items {
		if c.Id == class.Id {
			found = &c
			break
		}
	}
	if found == nil {
		return ClassDoesNotExist
	}

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
