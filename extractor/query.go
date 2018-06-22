package extractor

import (
	"net/http"
	"github.com/pkg/errors"
)

type KeyExtractor interface {
	GetKey(r *http.Request) (string, error)
}

type QueryKeyExtractor struct {
	key string
}

func (e *QueryKeyExtractor) GetKey(r *http.Request) (string, error) {
	keys, ok := r.URL.Query()[e.key]
	if !ok || len(keys) != 1 {
		return "", errors.New("query key not exists")
	}

	return keys[0], nil
}

func NewQueryKeyExtractor(key string) KeyExtractor {
	return &QueryKeyExtractor{key}
}