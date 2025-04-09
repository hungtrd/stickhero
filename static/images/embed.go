package images

import (
	_ "embed"
)

var (
	//go:embed gophers.png
	Gophers_png []byte

	//go:embed background.jpg
	Background_png []byte
)
