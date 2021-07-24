package query

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/utils"
)

type Builder struct {
	Filters       *Filters
	limit         *int64
	offset        *int64
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

	qb.limit = &limit

	return nil
}

func (qb *Builder) SetOffsetFromCtx(ctx *gin.Context) error {
	page, err := utils.GetQueryPage(ctx)
	if err != nil {
		return err
	}

	if qb.limit == nil {
		return errors.New("limit must be set first")
	}

	offset := (page - 1) * (*qb.limit)

	qb.offset = &offset

	return nil
}

func (qb *Builder) SetLimit(limit int64) {
	qb.limit = &limit
}

func (qb *Builder) SetOffset(offset int64) {
	qb.offset = &offset
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

func (qb *Builder) UsePagination(use bool) *Builder {
	qb.usePagination = use

	return qb
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
