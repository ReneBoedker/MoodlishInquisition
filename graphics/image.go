package graphics

import (
	"io"
)

type Image interface {
	Filetype() string
	ToHtml(w io.Writer)
	ToBase64(w io.Writer)
}
