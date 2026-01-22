package media

type VariantSpec struct {
	Name    string // "small" | "medium" | "large"
	Width   int
	Height  int
	Mode    FitMode // Cover (crop) vs Contain (letterbox)
	Format  Format  // JPEG/PNG/WebP
	Quality int     // for JPEG/WebP
}

type Variant struct {
	Name        string
	ContentType string
	Bytes       []byte
}

type FitMode string

const (
	// Cover scales the image up/down and crops overflow to exactly match WxH.
	// Good for thumbnails/profile crops.
	FitCover FitMode = "cover"

	// Contain scales the image to fit within WxH without cropping.
	// (You may pad/letterbox depending on implementation.)
	FitContain FitMode = "contain"
)

type Size string

const (
	Small  Size = "small"
	Medium Size = "medium"
	Large  Size = "large"
)

type ImageType string

const (
	Profile   ImageType = "profile"
	Thumbnail ImageType = "thumbnail"
)

type Format string

const (
	FormatJPEG Format = "jpeg"
	FormatPNG  Format = "png"
	FormatWEBP Format = "webp"
)

func (f Format) ContentType() string {
	switch f {
	case FormatJPEG:
		return "image/jpeg"
	case FormatPNG:
		return "image/png"
	case FormatWEBP:
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
