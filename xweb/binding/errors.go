package binding

import "fmt"

// ValidateParamsError is the error returned when a params is not either a pointer to a struct or a string.
type ValidateParamsError struct {
	error
}

type BindBodyError struct {
	error
	ContentType string
}

type ExtractError struct {
	error
	Tag string
}

type FieldSetterError struct {
	FieldSetterContext FieldSetterContext
	Message            string
}

func (f FieldSetterError) Error() string {
	return fmt.Sprintf("%s: %v", f.Message, f.FieldSetterContext)
}

type FieldSetterContext struct {
	Value            string `json:"value,omitempty"`
	FieldType        string `json:"field_type,omitempty"`
	ValueType        string `json:"value_type,omitempty"`
	Path             string `json:"path,omitempty"`
	ValueIndex       int    `json:"value_index,omitempty"`
	DecodingStrategy string `json:"decoding_strategy,omitempty"`
}
