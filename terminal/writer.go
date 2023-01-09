package terminal

import "io"

// Writer returns a new terminal writer with color support.
type Writer struct {
	out   io.Writer
	color bool
}

// NewWriter returns a new terminal writer with color support.
func NewWriter(out io.Writer) *Writer {
	return &Writer{
		out: out,
	}
}

// NewWriterColor returns a new terminal writer with color support.
func NewWriterColor(out io.Writer, color bool) *Writer {
	return &Writer{
		out:   out,
		color: color,
	}
}

// Color

// Color sets the color.
func (w *Writer) Color(color string) error {
	if !w.color {
		return nil
	}
	return w.write(color)
}

// ResetColor resets the color.
func (w *Writer) ResetColor() error {
	if !w.color {
		return nil
	}
	return w.write(FgReset)
}

// Colors

// Default resets the color to the default color.
func (w *Writer) Default() error {
	if !w.color {
		return nil
	}
	return w.write(FgDefault)
}

// Black sets the color to black.
func (w *Writer) Black() error {
	if !w.color {
		return nil
	}
	return w.write(FgBlack)
}

// Red sets the color to red.
func (w *Writer) Red() error {
	if !w.color {
		return nil
	}
	return w.write(FgRed)
}

// Green sets the color to green.
func (w *Writer) Green() error {
	if !w.color {
		return nil
	}
	return w.write(FgGreen)
}

// Yellow sets the color to yellow.
func (w *Writer) Yellow() error {
	if !w.color {
		return nil
	}
	return w.write(FgYellow)
}

// Blue sets the color to blue.
func (w *Writer) Blue() error {
	if !w.color {
		return nil
	}
	return w.write(FgBlue)
}

// Magenta sets the color to magenta.
func (w *Writer) Magenta() error {
	if !w.color {
		return nil
	}
	return w.write(FgMagenta)
}

// Cyan sets the color to cyan.
func (w *Writer) Cyan() error {
	if !w.color {
		return nil
	}
	return w.write(FgCyan)
}

// White sets the color to white.
func (w *Writer) White() error {
	if !w.color {
		return nil
	}
	return w.write(FgWhite)
}

// Gray sets the color to gray.
func (w *Writer) Gray() error {
	if !w.color {
		return nil
	}
	return w.write(FgGray)
}

// Write

// Write writes bytes to the writer.
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.out.Write(p)
}

// WriteLine writes a string with a new line to the writer.
func (w *Writer) WriteLine(s string) (n int, err error) {
	n, err = io.WriteString(w.out, s)
	if err != nil {
		return
	}

	n1, err := io.WriteString(w.out, "\n")
	n += n1
	return
}

// WriteString writes a string.
func (w *Writer) WriteString(s string) (n int, err error) {
	return io.WriteString(w.out, s)
}

// internal

func (w *Writer) write(s string) error {
	_, err := w.WriteString(s)
	return err
}
