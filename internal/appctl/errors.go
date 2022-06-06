package appctl

type appError string

const (
	ErrWrongState  appError = "wrong application state"
	ErrMainOmitted appError = "main function is omitted"
	ErrShutdown    appError = "application is in shutdown state"
	ErrTermTimeout appError = "termination timeout"
)

func (e appError) Error() string {
	return string(e)
}

type arrError []error

func (e arrError) Error() string {
	if len(e) == 0 {
		return "something went wrong"
	}
	s := "the following errors occurred:"
	for _, err := range e {
		s += "\n" + err.Error()
	}
	return s
}
