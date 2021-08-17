package query

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/utils"
)

type Builder struct {
	Filters       *Filters
	limit         int64
	offset        int64
	usePagination bool
}

func New() *Builder {
	return &Builder{
		Filters: NewFilters(),
	}
}

func (qb *Builder) SetLimitFromCtx(ctx *gin.Context) error {
	limit, err := utils.GetQueryLimit(ctx)
	if err != nil {
		return err
	}

	qb.limit = limit

	return nil
}

func (qb *Builder) SetOffsetFromCtx(ctx *gin.Context) error {
	page, err := utils.GetQueryPage(ctx)
	if err != nil {
		return err
	}

	offset := (page - 1) * qb.limit

	qb.offset = offset

	return nil
}

func (qb *Builder) SetLimit(limit int64) {
	qb.limit = limit
}

func (qb *Builder) SetOffset(offset int64) {
	qb.offset = offset
}

// returns a copy of the original query builder. chain and discard
func (qb *Builder) WithPagination(offset int64, limit int64) *Builder {
	nqb := *qb
	nqb.usePagination = true
	nqb.SetOffset(offset)
	nqb.SetLimit(limit)
	return &nqb
}

func (qb *Builder) WithPaginationFromCtx(ctx *gin.Context) *Builder {
	limit, _ := utils.GetQueryLimit(ctx)

	page, _ := utils.GetQueryPage(ctx)
	offset := (page - 1) * limit

	nqb := qb.WithPagination(offset, limit)

	return nqb
}

func (qb *Builder) Run(query string) (string, []interface{}) {
	where, params := qb.buildWhere()

	var offsetFilter, limitFilter string

	if qb.usePagination {
		// add offset
		params = append(params, qb.offset)
		offsetFilter = fmt.Sprintf("offset $%d", len(params))

		// add limit
		params = append(params, qb.limit)
		limitFilter = fmt.Sprintf("limit $%d", len(params))
	}

	query = strings.Replace(query, FiltersIdentifier, where, 1)
	query = strings.Replace(query, OffsetIdentifier, offsetFilter, 1)
	query = strings.Replace(query, LimitIdentifier, limitFilter, 1)

	return query, params
}

func (qb *Builder) buildWhere() (string, []interface{}) {
	var where string
	var params []interface{}

	if len(*qb.Filters) > 0 {
		for _, filter := range *qb.Filters {
			if where != "" {
				where += " and "
			} else {
				where += "where "
			}

			if filter.Where == "raw" {
				where += fmt.Sprintf("%s", filter.Value)
				continue
			}

			switch filter.Value.(type) {
			case []string:
				params = append(params, filter.Value)

				where += fmt.Sprintf("%s %s ANY($%d)", filter.Key, filter.Where, len(params))
			default:
				params = append(params, filter.Value)

				where += fmt.Sprintf("%s %s $%d", filter.Key, filter.Where, len(params))
			}
		}
	}

	return where, params
}
