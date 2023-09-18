package contracts

type ValidationError struct {
	Errors map[string][]string
	Global []string
}

func (v ValidationError) GetFieldErrors(field string) []string {
	return v.Errors[field]
}

func (v ValidationError) GetFirstFieldError(field string) string {
	if len(v.Errors[field]) > 0 {
		return v.Errors[field][0]
	}
	return ""
}

func (v ValidationError) HasErrors(field string) bool {
	return len(v.Errors[field]) > 0
}

func (v ValidationError) GetGlobalErrors() []string {
	return v.Global
}

func (v ValidationError) GetFirstGlobalError() string {
	if len(v.Global) > 0 {
		return v.Global[0]
	}
	return ""
}

func (v ValidationError) HasAnyErrors() bool {
	return len(v.Errors) > 0
}

type ValidationErrorHandler interface {
	Handle(err error) ValidationError
}
