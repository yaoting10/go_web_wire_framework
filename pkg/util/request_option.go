package util

type RequestOption[T any] func(req T)

func ApplyOptions[T any](r T, options ...RequestOption[T]) {
	for _, op := range options {
		op(r)
	}
}
