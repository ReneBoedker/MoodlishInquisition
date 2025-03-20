[![Go Report Card](https://goreportcard.com/badge/github.com/ReneBoedker/MoodlishInquisition)](https://goreportcard.com/report/github.com/ReneBoedker/MoodlishInquisition)
[![Go Reference](https://pkg.go.dev/badge/github.com/ReneBoedker/MoodlishInquisition.svg)](https://pkg.go.dev/github.com/ReneBoedker/MoodlishInquisition)

# Moodlish Inquisition
Moodlish Inquisition is a package for generating questions and question banks that can be imported in [Moodle](https://moodle.com/).

The package will for instance allow you to
- Programmatically generate banks of questions
- Generate graphics using [TikZ](https://www.ctan.org/pkg/pgf), and include them in questions and answers
- Use various question types such as multiple-choice, numerical, drop markers etc.

## Installation
To use MoodlishInquisition, simply import the needed package in your go source code. Afterwards, run `go mod tidy` to fetch the packages.

In order to use the TikZ-features, you need a working LaTeX-distribution as well as either [pdftocairo](https://manpages.ubuntu.com/manpages/noble/en/man1/pdftocairo.1.html) or [pdf2svg](https://cityinthesky.co.uk/opensource/pdf2svg/).

## Question types
For a list of supported question types and examples, see the documentation of the `moodle` subpackage.
