package helpers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type URLParams struct {
	pathParts   []string
	queryParams map[string][]string
}

func NewURLParams(r *http.Request) *URLParams {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	cleanParts := make([]string, 0, len(pathParts))
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}
	return &URLParams{
		pathParts:   cleanParts,
		queryParams: r.URL.Query(),
	}
}

func (p *URLParams) GetLastPathInt() (int, error) {
	if len(p.pathParts) == 0 {
		return 0, errors.New("no path parts")
	}
	return strconv.Atoi(p.pathParts[len(p.pathParts)-1])
}

func (p *URLParams) GetQueryParam(index int) (int, error) {
	if index < 0 || index >= len(p.queryParams) {
		return 0, errors.New("invalid query param index")
	}
	return strconv.Atoi(p.pathParts[index])
}

func (p *URLParams) GetQueryIntOr(key string, defailtValue int) int {
	if values, exists := p.queryParams[key]; !exists || len(values) == 0 {
		return defailtValue
	}
	result, err := strconv.Atoi(p.queryParams[key][0])
	if err != nil {
		return defailtValue
	}
	return result
}

func (p *URLParams) GetQueryOr(key, defaultValue string) string {
	if values := p.queryParams[key]; len(values) > 0 {
		return values[0]
	}
	return defaultValue
}

func (p *URLParams) GetQueryBoolOr(key string, defaultValue bool) bool {
	if values, exists := p.queryParams[key]; !exists || len(values) == 0 {
		return defaultValue
	}

	switch p.queryParams[key][0] {
	case "true":
		return true
	case "false":
		return false
	default:
		return defaultValue
	}
}
