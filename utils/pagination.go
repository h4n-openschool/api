package utils

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
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

func GeneratePaginationData(prefix string, total int, pq PaginationQuery) api.PaginationData {
	paginationData := api.PaginationData{
		Total:   total,
		Page:    pq.Page,
		PerPage: pq.PerPage,
	}

	paginationData.NextUrl = fmt.Sprintf("%s?page=%v&perPage=%v", prefix, getNextPage(pq.Page, total, pq.PerPage), pq.PerPage)
	paginationData.PrevUrl = fmt.Sprintf("%s?page=%v&perPage=%v", prefix, getPrevPage(pq.Page), pq.PerPage)

	paginationData.LastUrl = fmt.Sprintf("%s?page=%v&perPage=%v", prefix, getLastPage(total, pq.PerPage), pq.PerPage)
	paginationData.FirstUrl = fmt.Sprintf("%s?page=1&perPage=%v", prefix, pq.PerPage)

	return paginationData
}

func getLastPage(total int, perPage int) int {
	return total / perPage
}

func getNextPage(page int, total int, perPage int) int {
	next := page + 1
	last := getLastPage(total, perPage)

	if next >= last {
		return last
	}

	return next
}

func getPrevPage(page int) int {
	prev := page - 1

	if prev == 0 {
		return 1
	}

	return prev
}
