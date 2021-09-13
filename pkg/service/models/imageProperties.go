package models

type Params struct {
	Width,
	Height int
	Crop,
	ConvToGrayscale bool
}

type ImageProperties struct {
	ClientUUID,
	ResourcePath,
	ResourceName string
	Params Params
}
