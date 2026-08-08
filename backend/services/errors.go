package services

type AppError struct {
	Status  int
	Message string
}

func (err *AppError) Error() string { return err.Message }
func NewError(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}
