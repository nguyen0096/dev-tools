package iostreams

import (
	"fmt"
	"io"
	"os"
)

// IOStreams is an interface for interacting with the standard input, output, and error streams.
type IOStreams interface {
	Errorf(string, ...interface{})
	Printf(string, ...interface{})
}

type ioStreams struct {
	in     io.ReadCloser
	out    io.Writer
	errOut io.Writer
}

// New returns a new IOStreams instance.
func New() IOStreams {
	return &ioStreams{
		in:     os.Stdin,
		out:    os.Stdout,
		errOut: os.Stdout,
	}
}

func (c *ioStreams) Errorf(msg string, args ...any) {
	fmt.Fprintf(c.errOut, msg, args...)
	fmt.Fprintf(c.errOut, "\n")
}

func (c *ioStreams) Printf(msg string, args ...any) {
	fmt.Fprintf(c.out, msg, args...)
	fmt.Fprintf(c.out, "\n")
}
