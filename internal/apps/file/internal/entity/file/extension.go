package file

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

func (e Extension) String() string {
	return string(e)
}
