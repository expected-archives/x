package binding

import (
	"encoding"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"io"
	"net/http"
	"reflect"
	"strings"
)

var ErrInvalidParam = errors.New("invalid params: only ptr to a struct or string are accepted")

var MaxBodySize = int64(256 * 1024)

// Binder allow to bind params from a request to a struct.
type Binder struct {
	structCaches map[reflect.Type]StructCache

	stringsExtractors []StringsParamExtractor
	valueExtractors   []ValueParamExtractor

	stringsTags []string
	valuesTags  []string
}

func NewBinder(stringsExtractors []StringsParamExtractor, valueExtractors []ValueParamExtractor) *Binder {
	b := Binder{
		stringsExtractors: stringsExtractors,
		valueExtractors:   valueExtractors,

		structCaches: make(map[reflect.Type]StructCache),
	}

	stringsTags := lo.Map(
		b.stringsExtractors, func(item StringsParamExtractor, index int) string {
			return item.Tag()
		})

	valuesTags := lo.Map(
		b.valueExtractors, func(item ValueParamExtractor, index int) string {
			return item.Tag()
		})

	b.stringsTags = stringsTags
	b.valuesTags = valuesTags

	return &b
}

func (b *Binder) Bind(request *http.Request, params any) []error {
	errors := make([]error, 0)

	dec := NewDecoder(
		request,
		b.stringsExtractors,
		b.valueExtractors,
	)

	if err := validateParam(params); err != nil {
		return []error{err}
	}

	paramsType := reflect.TypeOf(params)

	if err := b.bindBody(request, params); err != nil {
		errors = append(errors, err)
	}

	if reflect.ValueOf(params).Kind() != reflect.String {
		var (
			structCache StructCache
			ok          bool
		)

		if structCache, ok = b.structCaches[paramsType]; !ok {
			structCache = NewStructAnalyzer(b.stringsTags, b.valuesTags, paramsType).Cache()
			b.structCaches[paramsType] = structCache
		}

		if errs := dec.Decode(structCache, reflect.ValueOf(params)); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return errors

}

// validateParam validate if the param is valid.
// Accepted values :
// - pointer to a struct
// - string
func validateParam(param any) error {
	ref := reflect.ValueOf(param)

	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()

		if ref.Kind() == reflect.Struct {
			return nil
		}
	}

	if ref.Kind() == reflect.String {
		return nil
	}

	return ValidateParamsError{error: ErrInvalidParam}
}

// bindBody bind the body of the request to the params.
// it supports 3 types of content-type:
// - application/json
// - application/xml
// - text/plain
func (b *Binder) bindBody(r *http.Request, params any) error {
	if r.ContentLength == 0 {
		return nil
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		bytesBody, err := io.ReadAll(io.LimitReader(r.Body, MaxBodySize))
		if err != nil {
			return BindBodyError{
				error:       fmt.Errorf("unable to read body: %w", err),
				ContentType: "application/json",
			}
		}

		if err := json.Unmarshal(bytesBody, params); err != nil {
			return BindBodyError{
				error:       fmt.Errorf("unable to unmarchal json request: %w", err),
				ContentType: "application/json",
			}
		}
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/xml") {
		bytesBody, err := io.ReadAll(io.LimitReader(r.Body, MaxBodySize))
		if err != nil {
			return BindBodyError{
				error:       fmt.Errorf("unable to read body: %w", err),
				ContentType: "application/xml",
			}
		}

		if err := xml.Unmarshal(bytesBody, params); err != nil {
			return BindBodyError{
				error:       fmt.Errorf("unable to unmarchal xml request: %w", err),
				ContentType: "application/xml",
			}
		}
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
		if !reflect.TypeOf(params).AssignableTo(TextUnmarshaller) {
			return nil
		}

		bytesBody, err := io.ReadAll(io.LimitReader(r.Body, MaxBodySize))
		if err != nil {
			return fmt.Errorf("unable to read body: %w", err)
		}

		if err := params.(encoding.TextUnmarshaler).UnmarshalText(bytesBody); err != nil {
			return fmt.Errorf("unable to read text request: %w", err)
		}
	}

	return nil

}
