package query

type Filter struct {
	Key   string
	Value interface{}
	Where string
}

type Filters []Filter

func NewFilters() *Filters {
	return new(Filters)
}

func (f *Filters) Add(key string, value interface{}, opts ...interface{}) *Filters {
	var where = "="
	if len(opts) > 0 {
		where = opts[0].(string)
	}

	*f = append(*f, Filter{
		Key:   key,
		Value: value,
		Where: where,
	})

	return f
}
