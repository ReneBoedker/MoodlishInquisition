package moodle

import (
	"io"
)

// Question is the common interface of all question types
type Question interface {
	ToXml(io.Writer)
	MoodleName() string
}
