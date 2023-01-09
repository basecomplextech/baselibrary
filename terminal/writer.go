package terminal

import "io"

// Writer returns a new terminal writer with color support.
type Writer struct {
	Out      io.Writer
	Colorize bool
}

// NewWriter returns a new terminal writer with color support.
func NewWriter(out io.Writer) *Writer {
	return &Writer{
		Out: out,
	}
}

// NewWriterColor returns a new terminal writer with color support.
func NewWriterColor(out io.Writer, color bool) *Writer {
	return &Writer{
		Out:      out,
		Colorize: color,
	}
}

// Color

// Color sets the color.
func (w *Writer) Color(color string) {
	w.WriteString(color)
}

// Default resets the color to the default color.
func (w *Writer) Default() {
	w.WriteString(FgDefault)
}

// Black sets the color to black.
func (w *Writer) Black() {
	w.WriteString(FgBlack)
}

// Red sets the color to red.
func (w *Writer) Red() {
	w.WriteString(FgRed)
}

// Green sets the color to green.
func (w *Writer) Green() {
	w.WriteString(FgGreen)
}

// Yellow sets the color to yellow.
func (w *Writer) Yellow() {
	w.WriteString(FgYellow)
}

// Blue sets the color to blue.
func (w *Writer) Blue() {
	w.WriteString(FgBlue)
}

// Magenta sets the color to magenta.
func (w *Writer) Magenta() {
	w.WriteString(FgMagenta)
}

// Cyan sets the color to cyan.
func (w *Writer) Cyan() {
	w.WriteString(FgCyan)
}

// White sets the color to white.
func (w *Writer) White() {
	w.WriteString(FgWhite)
}

// Write

// Write writes bytes to the writer.
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Out.Write(p)
}

// WriteLine writes a string with a new line to the writer.
func (w *Writer) WriteLine(s string) (n int, err error) {
	n, err = io.WriteString(w.Out, s)
	if err != nil {
		return
	}

	n1, err := io.WriteString(w.Out, "\n")
	n += n1
	return
}

// WriteString writes a string.
func (w *Writer) WriteString(s string) (n int, err error) {
	return io.WriteString(w.Out, s)
}
