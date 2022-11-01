package repos

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/h4n-openschool/classes/models"
	"github.com/lucsky/cuid"
)

type InMemoryClassRepository struct {
	Items []models.Class
}

func NewInMemoryClassRepository(itemCount int) InMemoryClassRepository {
	var items []models.Class

	for i := 0; i < itemCount; i++ {
		id := cuid.New()
		log.Println(id)
		items = append(items, models.Class{
			Id:          id,
			Name:        fmt.Sprintf(`class-%v`, i),
			DisplayName: fmt.Sprintf(`Class %v`, i),
			Description: fmt.Sprintf(`This is class %v`, i),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return InMemoryClassRepository{Items: items}
}

func (r *InMemoryClassRepository) GetAll(perPage int, page int) ([]models.Class, error) {
	var items []models.Class

	offset := (page - 1) * perPage
	if page == 1 {
		offset = 0
	}

	for i := offset; i < (offset + perPage); i++ {
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
