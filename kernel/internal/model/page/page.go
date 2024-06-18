package page

type Information[T any] struct {
	Size          int   `json:"size"`
	Page          int   `json:"page"`
	TotalPage     int64 `json:"totalPage"`
	TotalElements int64 `json:"totalElements"`
	Data          []T   `json:"data"`
}

func FromSlice[T any](s []T, page, size int, totalPage, totalElements int64) Information[T] {
	if s == nil {
		s = make([]T, 0)
	}
	return Information[T]{
		Size:          size,
		Page:          page,
		TotalPage:     totalPage,
		TotalElements: totalElements,
		Data:          s,
	}
}
func FromSliceWithMap[I, O any](s []I, page, size int, totalPage, totalElements int64, mapFn func(I) O) Information[O] {
	var outputs []O
	for _, i := range s {
		outputs = append(outputs, mapFn(i))
	}
	return FromSlice(outputs, page, size, totalPage, totalElements)
}
