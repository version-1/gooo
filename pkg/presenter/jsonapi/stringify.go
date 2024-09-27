package jsonapi

func Stringify(v any) string {
	s, err := Escape(v)
	if err != nil {
		panic(err)
	}

	return s
}
