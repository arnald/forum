package helpers

import (
	"log"
	"strconv"
	"strings"

	"github.com/arnald/forum/cmd/client/domain"
)

const (
	hexMinLength = 4
	hexMaxLength = 7
)

func PrepareCategories(categories []domain.Category) []domain.Category {
	for i := range categories {
		categories[i].Color = normalizeColor(categories[i].Color)
	}
	return categories
}

func normalizeColor(color string) string {
	// Add # prefix if missing
	if !strings.HasPrefix(color, "#") {
		color = "#" + color
	}

	// Validate hex format
	if !isValidHexColor(color) {
		log.Printf("Invalid color format: %s", color)
		return "#00c6ff"
	}

	return strings.ToUpper(color)
}

func isValidHexColor(s string) bool {
	if len(s) != hexMinLength && len(s) != hexMaxLength {
		return false
	}

	if s[0] != '#' {
		return false
	}

	hex := s[1:]
	_, err := strconv.ParseUint(hex, 16, 64)
	return err == nil
}
