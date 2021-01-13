package debugger

import (
	"fmt"
	"github.com/go-delve/delve/service/api"
	"path"
)

type steps []step
type step struct {
	meta      meta
	file      string
	variables variables
}
type meta struct {
	currentPosition currentPosition
	name            string
}
type currentPosition struct {
	line   int
	column int
}

func NewSteps() *steps {
	s := make(steps, 0)
	return &s
}

func (s *steps) AddStep(variables []api.Variable, stackTrace []api.Stackframe) *step {
	st := step{}
	_, fileName := path.Split(stackTrace[0].File)
	currentLine := stackTrace[0].Line
	st.file = stackTrace[0].File
	st.meta.currentPosition.line = currentLine
	st.meta.name = fmt.Sprintf(
		"%v at %v:%v",
		stackTrace[0].Function.Name(),
		fileName,
		currentLine,
	)

	vars := NewVariables()

	for _, val := range variables {
		vars.Add(val, currentLine)
	}

	st.variables = *vars

	*s = append(*s, st)

	return &st
}
