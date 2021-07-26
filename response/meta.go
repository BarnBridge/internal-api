package response

type ResponseMeta map[string]interface{}

func (m ResponseMeta) Set(key string, value interface{}) ResponseMeta {
	m[key] = value

	return m
}

func Meta() ResponseMeta {
	return make(ResponseMeta)
}
