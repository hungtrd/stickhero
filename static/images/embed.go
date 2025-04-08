package images

import (
	_ "embed"
)

var (
	//go:embed gophers.png
	Gophers_png []byte

	//go:embed gophers.png
	Spritesheet_png []byte
)
