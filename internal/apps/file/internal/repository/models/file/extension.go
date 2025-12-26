package file

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Extension string

const (
	ExtensionJPEG Extension = "jpeg"
	ExtensionPNG  Extension = "png"
	ExtensionGIF  Extension = "gif"
	ExtensionBMP  Extension = "bmp"
	ExtensionTIFF Extension = "tiff"
	ExtensionICO  Extension = "ico"
	ExtensionWEBP Extension = "webp"
	ExtensionSVG  Extension = "svg"
	ExtensionHEIC Extension = "heic"
	ExtensionHEIF Extension = "heif"
)

func ExtensionFromString(s string) (Extension, error) {
	switch s {
	case "jpeg":
		return ExtensionJPEG, nil
	case "png":
		return ExtensionPNG, nil
	case "gif":
		return ExtensionGIF, nil
	case "bmp":
		return ExtensionBMP, nil
	case "tiff":
		return ExtensionTIFF, nil
	case "ico":
		return ExtensionICO, nil
	case "webp":
		return ExtensionWEBP, nil
	case "svg":
		return ExtensionSVG, nil
	case "heic":
		return ExtensionHEIC, nil
	case "heif":
		return ExtensionHEIF, nil
	default:
		return "", errors.New("invalid extension")
	}
}

func ExtensionFromFileName(fileName string) (Extension, error) {
	return ExtensionFromString(strings.TrimPrefix(filepath.Ext(fileName), "."))
}

func (e Extension) String() string {
	return string(e)
}
