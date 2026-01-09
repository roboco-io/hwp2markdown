package ir

// ImageBlock represents an image reference in the document.
type ImageBlock struct {
	ID       string `json:"id"`                  // internal image ID
	Path     string `json:"path,omitempty"`      // extracted file path
	OrigName string `json:"orig_name,omitempty"` // original filename
	Alt      string `json:"alt,omitempty"`       // alt text
	Caption  string `json:"caption,omitempty"`   // image caption
	Width    int    `json:"width,omitempty"`     // width in pixels
	Height   int    `json:"height,omitempty"`    // height in pixels
	Format   string `json:"format,omitempty"`    // png, jpg, gif, bmp, etc.
	Data     []byte `json:"-"`                   // raw image data (not serialized)
}

// NewImage creates a new image block with the given ID.
func NewImage(id string) *ImageBlock {
	return &ImageBlock{
		ID: id,
	}
}

// SetDimensions sets the width and height of the image.
func (img *ImageBlock) SetDimensions(width, height int) {
	img.Width = width
	img.Height = height
}

// HasData returns true if the image has raw data loaded.
func (img *ImageBlock) HasData() bool {
	return len(img.Data) > 0
}
