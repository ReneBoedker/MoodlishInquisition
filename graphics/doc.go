// Package graphics enables inclusion of images in generated questions.
//
// It is primarily intended to be used for generating graphics using TikZ or
// pgfplots and converting them to a Moodle-ready format. But it also has
// experimental support for inclusion of arbitrary image files.
//
// Pictures can be converted to two different formats depending on the use case.
// For inclusion directly into question descriptions or answers, use
// CompileToHtml, which generates an SVG string that Moodle can render. For
// inclusion in the DropMarker type question, use CompileToBase64.
//
// The package relies on external tools to perform the compilation and
// conversion of TikZ graphics. More precisely, pdflatex and pdf2svg must be
// installed for the functions to succeed. If cropping of figures is requested,
// the package will call Inkscape.
package graphics
