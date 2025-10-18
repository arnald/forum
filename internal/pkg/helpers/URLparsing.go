package helpers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type URLParams struct {
	request *http.Request
}

func NewURLParams(r *http.Request) *URLParams {
	return &URLParams{
		request: r,
	}
}

func (p *URLParams) GetQueryInt(key string) (int, error) {
	value := p.request.URL.Query().Get(key)
	if value == "" {
		return 0, fmt.Errorf("parameter %s is required", key)
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parameter %s must be a valid integer", key)
	}

	return result, nil
}

func (p *URLParams) GetQueryIntOr(key string, defaultValue int) int {
	value := p.request.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return result
}

func (p *URLParams) GetQueryString(key string) (string, error) {
	value := p.request.URL.Query().Get(key)
	if value == "" {
		return "", fmt.Errorf("parameter %s is required", key)
	}
	return value, nil
}

func (p *URLParams) GetQueryStringOr(key, defaultValue string) string {
	value := p.request.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (p *URLParams) GetQueryBool(key string) (bool, error) {
	value := p.request.URL.Query().Get(key)
	if value == "" {
		return false, fmt.Errorf("parameter %s is required", key)
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("parameter %s must be true or false", key)
	}
}

func (p *URLParams) GetQueryBoolOr(key string, defaultValue bool) bool {
	value := p.request.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return defaultValue
	}
}

func (p *URLParams) GetMultiple(key string) []string {
	return p.request.URL.Query()[key]
}

type PaginationParams struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (p *URLParams) GetPagination() PaginationParams {
	page := p.GetQueryIntOr("page", 1)
	limit := p.GetQueryIntOr("limit", 20)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func GetQueryInt(r *http.Request, key string) (int, error) {
	params := NewURLParams(r)
	return params.GetQueryInt(key)
}

func GetQueryIntOr(r *http.Request, key string, defaultValue int) int {
	params := NewURLParams(r)
	return params.GetQueryIntOr(key, defaultValue)
}

func GetQueryString(r *http.Request, key string) (string, error) {
	params := NewURLParams(r)
	return params.GetQueryString(key)
}

func GetQueryStringOr(r *http.Request, key, defaultValue string) string {
	params := NewURLParams(r)
	return params.GetQueryStringOr(key, defaultValue)
}

func GetQueryBool(r *http.Request, key string) (bool, error) {
	params := NewURLParams(r)
	return params.GetQueryBool(key)
}

func GetQueryBoolOr(r *http.Request, key string, defaultValue bool) bool {
	params := NewURLParams(r)
	return params.GetQueryBoolOr(key, defaultValue)
}

func GetPagination(r *http.Request) PaginationParams {
	params := NewURLParams(r)
	return params.GetPagination()
}
