package binding

import (
	"github.com/go-chi/chi"
	"net/http"
)

// PathExtractor extract value from the chi router.
type PathExtractor struct{}

// Extract value from the chi router.
func (p PathExtractor) Extract(req *http.Request, valueOfTag string) ([]string, error) {
	str := chi.URLParam(req, valueOfTag)
	if str == "" {
		return nil, nil
	}

	return []string{str}, nil
}

// Tag return the tag name of this extractor.
func (p PathExtractor) Tag() string {
	return "path"
}
