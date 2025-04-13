package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func CreateSlugName(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Remove special characters (keep only letters, numbers, spaces, and dashes)
	re := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove duplicate hyphens
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

func GetTimeStampHash() string {
	data := fmt.Sprintf("%d", time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(encoded[:22], "=")
}
