package exceptions

// MyError simplify the error handling stage. All native errors must be this type
type MyError struct {
	Msg string `json:"message,omitempty"`
	Err error  `json:"err,omitempty"`
}

// Error implements the Error type
func (e MyError) Error() string {
	return e.Err.Error()
}

// newError creates a new MyError and returns it
func newError(msg string) MyError {
	return MyError{
		Msg: msg,
	}
}

// AddErr saves an error into MyError and returns it
func (e MyError) AddErr(err error) MyError {
	e.Err = err
	return e
}
