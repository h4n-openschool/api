package teachers

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/utils"
	"github.com/lucsky/cuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	TeacherDoesNotExist = errors.New("no existing class found by that id")
)

// InMemoryTeacherRepository implements the [TeacherRepository] interface using an
// in-memory slice of [models.Teacher] items.
type InMemoryTeacherRepository struct {
	// Items is the slice of [models.Teacher] items stored in memory.
	Items []models.Teacher
}

// NewInMemoryTeacherRepository creates a new instance of
// [NewInMemoryTeacherRepository]
func NewInMemoryTeacherRepository(itemCount int) InMemoryTeacherRepository {
	var items []models.Teacher

	password, err := bcrypt.GenerateFromPassword([]byte("password"), 10)
	if err != nil {
		panic(err)
	}

	// Generate classes in-memory to use with repo methods.
	for i := 0; i < itemCount; i++ {
		id := cuid.New()

		items = append(items, models.Teacher{
			BaseMetadata: models.BaseMetadata{
				Id:        id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			FullName:     fmt.Sprintf("%v %v", faker.FirstName(), faker.LastName()),
			Email:        faker.Email(),
			PasswordHash: string(password),
		})
	}

  items = append(items, models.Teacher{
    BaseMetadata: models.BaseMetadata{
      Id: "clb3x2ugq0004txk80dyoemxa",
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
    },
    FullName: "John Doe",
    Email: "john.doe@school.edu",
    PasswordHash: string(password),
  })

	// Return the new repository to the caller
	return InMemoryTeacherRepository{Items: items}
}

func (r *InMemoryTeacherRepository) GetAll(pq utils.PaginationQuery) ([]models.Teacher, error) {
	var items []models.Teacher

	offset := pq.Offset()
	for i := offset; i < (offset + pq.PerPage); i++ {
		if len(r.Items) <= i {
			break
		}
		item := r.Items[i]
		items = append(items, item)
	}

	return items, nil
}

func (r *InMemoryTeacherRepository) Get(id string) (*models.Teacher, error) {
	var found *models.Teacher

	for _, v := range r.Items {
		if v.Id == id {
			found = &v
			break
		}
	}

	return found, nil
}

func (r *InMemoryTeacherRepository) GetByEmail(email string) (*models.Teacher, error) {
	var found *models.Teacher

	for _, v := range r.Items {
		if v.Email == email {
			found = &v
			break
		}
	}

	return found, nil
}

func (r *InMemoryTeacherRepository) Update(teacher *models.Teacher) (*models.Teacher, error) {
	var found *models.Teacher

	for k, v := range r.Items {
		if v.Id == teacher.Id {
			v.FullName = teacher.FullName
			v.Email = teacher.Email
			v.UpdatedAt = time.Now()

			found = &v

			r.Items[k] = *found

			break
		}
	}

	if found == nil {
		return nil, TeacherDoesNotExist
	}

	return found, nil
}

func (r *InMemoryTeacherRepository) Create(teacher models.Teacher) (*models.Teacher, error) {
	model := models.Teacher{
		BaseMetadata: models.BaseMetadata{
			Id:        cuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		FullName: teacher.FullName,
		Email:    teacher.Email,
	}

	r.Items = append(r.Items, model)

	return &model, nil
}

func (r *InMemoryTeacherRepository) Delete(teacher models.Teacher) error {
	var newItems []models.Teacher

	var found *models.Teacher
	for _, c := range r.Items {
		if c.Id == teacher.Id {
			found = &c
			break
		}
	}
	if found == nil {
		return TeacherDoesNotExist
	}

	for _, c := range r.Items {
		if c.Id != teacher.Id {
			newItems = append(newItems, c)
		}
	}

	r.Items = newItems

	return nil
}

func (r *InMemoryTeacherRepository) Count() (int, error) {
	return len(r.Items), nil
}
