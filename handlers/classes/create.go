package classes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/classes/bus"
	"github.com/h4n-openschool/classes/models"
	repos "github.com/h4n-openschool/classes/repos/classes"
	"github.com/h4n-openschool/classes/utils"
	"github.com/rabbitmq/amqp091-go"
)

type CreateClassDto struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type CreateClassResponse struct {
	Class models.Class
}

// GetClasses godoc
// @Summary Get a paginated list of classes
// @Schemes
// @Tags classes
// @Accept json CreateClassDto
// @Produce json
// @Success 200 CreateClassResponse
// @Router /classes [get]
func CreateClass(repo repos.ClassRepository) gin.HandlerFunc {
	b := bus.GetOrCreateBus()

	return func(ctx *gin.Context) {
		body := CreateClassDto{}
		if err := ctx.BindJSON(&body); err != nil {
			utils.ErrValidation(ctx, err)
      return
		}

		in := models.Class{
			Name:        body.Name,
			DisplayName: body.DisplayName,
			Description: body.Description,
		}

		c, q, err := b.ChannelQueue()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		class, err := repo.Create(in)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ev := bus.Event[*models.Class]{
			Data: class,
			Metadata: bus.EventMeta{
				Type:     "create",
				Resource: "class",
			},
		}

		j, _ := json.Marshal(ev)

		err = c.Publish(
			"",
			q.Name,
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        j,
			},
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"class": class,
		})
	}
}
