package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"Page 1, Limit 10", 1, 10, 0},
		{"Page 2, Limit 10", 2, 10, 10},
		{"Page 3, Limit 20", 3, 20, 40},
		{"Page 0, Limit 10", 0, 10, 0},
		{"Page 1, Limit 0", 1, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetOffset(tt.page, tt.limit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		total    int
		limit    int
		expected int
	}{
		{"100 items, 10 per page", 100, 10, 10},
		{"105 items, 10 per page", 105, 10, 11},
		{"0 items, 10 per page", 0, 10, 0},
		{"5 items, 10 per page", 5, 10, 1},
		{"1 item, 1 per page", 1, 1, 1},
		{"0 limit", 100, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTotalPages(tt.total, tt.limit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name          string
		inputPage     int
		inputLimit    int
		expectedPage  int
		expectedLimit int
	}{
		{"Valid pagination", 1, 20, 1, 20},
		{"Page 0 should default to 1", 0, 20, 1, 20},
		{"Negative page should default to 1", -1, 20, 1, 20},
		{"Limit 0 should default to 20", 1, 0, 1, 20},
		{"Negative limit should default to 20", 1, -1, 1, 20},
		{"Limit > 100 should default to 100", 1, 150, 1, 100},
		{"High valid limit", 1, 100, 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PaginationRequest{
				Page:  tt.inputPage,
				Limit: tt.inputLimit,
			}
			ValidatePagination(req)
			assert.Equal(t, tt.expectedPage, req.Page)
			assert.Equal(t, tt.expectedLimit, req.Limit)
		})
	}
}

func TestNewPaginatedResponse(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	page := 1
	limit := 10
	total := 30

	result := NewPaginatedResponse(data, page, limit, total)

	assert.NotNil(t, result)
	assert.Equal(t, data, result.Data)
	assert.Equal(t, page, result.Pagination.Page)
	assert.Equal(t, limit, result.Pagination.Limit)
	assert.Equal(t, total, result.Pagination.Total)
	assert.Equal(t, 3, result.Pagination.TotalPages)
}
