package errors

import (
	"fmt"
	"runtime"
)

type Error struct {
	code  string
	msg   string
	stack *stack
}

func New(code, msg string) *Error {
	return &Error{
		code:  code,
		msg:   msg,
		stack: captureStack(),
	}
}

func (e Error) StackTrace() string {
	return fmt.Sprintf("%+v", e.stack)
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %s. msg: %s", e.code, e.msg)
}

func (e Error) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%+v\n", e.stack)
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(f, "%s", e.msg)
	case 'q':
		fmt.Fprintf(f, "%q", e.msg)
	}
}

type stack []frame

func (st *stack) Format(f fmt.State, c rune) {
	switch c {
	case 'v', 's':
		for _, fr := range *st {
			output := fr.String()
			if output != "" {
				fmt.Fprintln(f, fr.String())
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
		return ""
	}
	return fmt.Sprintf("%s:%s:%d", f.File(), f.FuncName(), f.Line())
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
	if f.name != nil {
		return *f.name
	}

	f.collect()

	return *f.name
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
