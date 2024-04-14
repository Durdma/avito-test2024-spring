package models

type ErrService struct {
	Status int
	Error  string
}

func NewErrorService(status int, err string) ErrService {
	return ErrService{
		Status: status,
		Error:  err,
	}
}
