package customerrors

type TpTooLowError struct {
	Message string
}

func (e *TpTooLowError) Error() string {
	return e.Message
}

type TpTooHighError struct {
	Message string
}

func (e *TpTooHighError) Error() string {
	return e.Message
}

type SlTooLowError struct {
	Message string
}

func (e *SlTooLowError) Error() string {
	return e.Message
}

type SlTooHighError struct {
	Message string
}

func (e *SlTooHighError) Error() string {
	return e.Message
}
