package models

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

func safeDeref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}

func ParseTimestamp(input string) (int, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, errors.New("invalid timestamp format, expected M:SS")
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %w", err)
	}

	if seconds < 0 || seconds > 59 {
		return 0, errors.New("seconds must be between 0 and 59")
	}

	total := minutes*60 + seconds
	return total, nil
}

func SecondsToMMSS(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	timeString := fmt.Sprintf("%02d:%02d", minutes, secs)
	if strings.HasPrefix(timeString, "0") {
		return timeString[1:]
	}
	return timeString
}
