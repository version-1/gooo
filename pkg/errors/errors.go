package errors

import (
	"errors"
	"fmt"
	"runtime"
)

type Error struct {
	err   error
	stack *stack
}

func Wrap(err error) *Error {
	return &Error{
		err:   err,
		stack: captureStack(),
	}
}

func New(msg string) *Error {
	return &Error{
		err:   errors.New(msg),
		stack: captureStack(),
	}
}

func (e Error) StackTrace() string {
	return fmt.Sprintf("%+v", e.stack)
}

func (e Error) Error() string {
	return fmt.Sprintf("pkg/errors : %s", e.err)
}

func (e Error) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%s\n", e.Error())
			fmt.Fprintln(f, "")
			fmt.Fprintf(f, "%+v\n", e.stack)
			return
		} else {
			fmt.Fprintf(f, "%s", e.Error())
			return
		}
	case 's':
		fmt.Fprintf(f, "%s", e.Error())
	}
}

type stack []frame

func (st *stack) Format(f fmt.State, c rune) {
	switch c {
	case 'v', 's':
		for _, fr := range *st {
			output := fr.String()
			if output != "" {
				fmt.Fprintln(f, output)
			}
		}
	}
}

type frame struct {
	pc   uintptr
	line *int
	file *string
	name *string
}

func (f frame) counter() uintptr { return uintptr(f.pc) - 1 }

func (f *frame) collect() {
	fn := runtime.FuncForPC(f.counter())
	if fn == nil {
		return
	}

	name := fn.Name()
	f.name = &name
	file, line := fn.FileLine(f.counter())

	f.file = &file
	f.line = &line
}

func (f frame) String() string {
	if f.file == nil {
		f.collect()
	}

	if f.file == nil {
		return ""
	}

	return fmt.Sprintf("%s. method: %s. line: %d", f.File(), f.FuncName(), f.Line())
}

func (f *frame) File() string {
	if f.file != nil {
		return *f.file
	}
	f.collect()

	return *f.file
}

func (f *frame) Line() int {
	if f.line != nil {
		return *f.line
	}
	f.collect()

	return *f.line
}

func (f *frame) FuncName() string {
	n := func(s string) string {
		for i := len(s) - 1; i > 0; i-- {
			if s[i] == '.' {
				return s[i+1:]
			}
		}
		return s
	}
	if f.name != nil {
		return n(*f.name)
	}

	f.collect()

	return n(*f.name)
}

func captureStack() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := make([]frame, n)
	for _, pc := range pcs[0:n] {
		frames = append(frames, frame{pc: pc})
	}
	st := stack(frames)
	return &st
}
