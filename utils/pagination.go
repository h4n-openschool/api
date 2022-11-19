package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPerPage = 10
	DefaultPage    = 1
)

type PaginationQuery struct {
	PerPage int `form:"perPage" json:"perPage"`
	Page    int `form:"page" json:"page"`
}

func NewPaginationQuery() PaginationQuery {
	return PaginationQuery{
		PerPage: DefaultPerPage,
		Page:    DefaultPage,
	}
}

func (pq *PaginationQuery) Read(ctx *gin.Context) {
	if err := ctx.ShouldBindQuery(&pq); err != nil {
		log.Println("binding pagination failed, using default values")
	}
	if pq.Page < 0 {
		log.Println("invalid page specified, using default value")
		pq.Page = DefaultPage
	}
	if pq.PerPage <= 0 {
		log.Println("invalid perpage specified, using default value")
		pq.PerPage = DefaultPerPage
	}
}

func (pq *PaginationQuery) ReadFromOptional(page *int, perPage *int) {
  if page != nil {
    pq.Page = *page
  }
  if perPage != nil {
    pq.PerPage = *perPage
  }
}

func (pq PaginationQuery) Offset() int {
	if pq.Page == 0 {
		pq.Page = 1
	}

	return pq.PerPage * (pq.Page - 1)
}
