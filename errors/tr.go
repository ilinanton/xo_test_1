package errors

type errorTR struct {
	s string
}

func (e *errorTR) Error() string {
	return e.s
}
