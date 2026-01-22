package media

import (
	"errors"
	"image"
)

func Process(src image.Image, specs []VariantSpec) ([]Variant, error) {
	if len(specs) == 0 {
		return nil, errors.New("no specs provided")
	}

	out := make([]Variant, 0, len(specs))
	for _, spec := range specs {
		if spec.Width <= 0 || spec.Height <= 0 {
			return nil, errors.New("invalid spec dimensions")
		}

		if spec.Quality <= 0 {
			spec.Quality = 85
		}

		var processed image.Image
		var err error
		switch spec.Mode {
		case FitCover:
			processed, err = cover(src, spec.Width, spec.Height)
		case FitContain:
			processed, err = contain(src, spec.Width, spec.Height, spec.Format)
		default:
			return nil, errors.New("unknown fit mode")
		}

		if err != nil {
			return nil, err
		}

		b, ct, err := encode(processed, spec.Format, spec.Quality)
		if err != nil {
			return nil, err
		}

		out = append(out, Variant{
			Name:        spec.Name,
			ContentType: ct,
			Bytes:       b,
		})
	}

	return out, nil
}
