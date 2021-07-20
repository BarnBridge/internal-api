package utils

import (
	"fmt"
)

type Filters []QueryFilter

func NewFilters() *Filters {
	return new(Filters)
}

func (f *Filters) Add(key string, value interface{}, opts ...interface{}) *Filters {
	var where = "="
	if len(opts) > 0 {
		where = opts[0].(string)
	}

	*f = append(*f, QueryFilter{
		Key:   key,
		Value: value,
		Where: where,
	})

	return f
}

type QueryFilter struct {
	Key   string
	Value interface{}
	Where string
}

func BuildQueryWithFilter(query string, filters *Filters, limit *int64, offset *int64) (string, []interface{}) {
	var where string
	var offsetFilter, limitFilter string
	var params []interface{}
	if len(*filters) > 0 {

		for _, filter := range *filters {
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

	if offset != nil {
		params = append(params, offset)

		offsetFilter = fmt.Sprintf("offset $%d", len(params))
	}

	if limit != nil {
		params = append(params, limit)

		limitFilter = fmt.Sprintf("limit $%d", len(params))
	}

	return fmt.Sprintf(query, where, offsetFilter, limitFilter), params
}
