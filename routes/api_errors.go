package routes

type InvalidParamsError struct {
	Message string
}

func (e *InvalidParamsError) Error() string {
	return e.Message
}

func (h *Handlers) InvalidParams(message string) error {
	return &InvalidParamsError{
		Message: message,
	}
}
