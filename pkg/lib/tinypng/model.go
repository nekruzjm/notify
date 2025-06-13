package tinypng

type Size struct {
	Format string
	Width  int
	Height int
}

type Response struct {
	Output struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Url    string `json:"url"`
	} `json:"output"`
}

type Request struct {
	Store  Store  `json:"store"`
	Resize Resize `json:"resize"`
}

type Store struct {
	Service            string `json:"service"`
	AwsAccessKeyId     string `json:"aws_access_key_id"`
	AwsSecretAccessKey string `json:"aws_secret_access_key"`
	Region             string `json:"region"`
	Path               string `json:"path"`
}

type Resize struct {
	Method string `json:"method"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Formats
const (
	_format075x   = "0.75x/"
	_format1x     = "1x/"
	_format15x    = "1.5x/"
	_format2x     = "2x/"
	_format3x     = "3x/"
	_format4x     = "4x/"
	_formatNotify = "notify/"
)

// Width/Height constants for each format
const (
	_size075xWidth    = 246
	_size075xHeight   = 144
	_size1xWidth      = 328
	_size1xHeight     = 190
	_size15xWidth     = 492
	_size15xHeight    = 285
	_size2xWidth      = 656
	_size2xHeight     = 380
	_size3xWidth      = 984
	_size3xHeight     = 570
	_size4xWidth      = 1312
	_size4xHeight     = 760
	_sizeNotifyWidth  = 1312
	_sizeNotifyHeight = 760
)

func Sizes() [7]Size {
	return [7]Size{
		{Format: _format075x, Width: _size075xWidth, Height: _size075xHeight},
		{Format: _format1x, Width: _size1xWidth, Height: _size1xHeight},
		{Format: _format15x, Width: _size15xWidth, Height: _size15xHeight},
		{Format: _format2x, Width: _size2xWidth, Height: _size2xHeight},
		{Format: _format3x, Width: _size3xWidth, Height: _size3xHeight},
		{Format: _format4x, Width: _size4xWidth, Height: _size4xHeight},
		{Format: _formatNotify, Width: _sizeNotifyWidth, Height: _sizeNotifyHeight},
	}
}
