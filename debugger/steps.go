package debugger

import (
	"fmt"
	"github.com/go-delve/delve/service/api"
	"path"
)

type steps map[string]step
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
type variables map[string]variable
type variable struct {
	name  string
	kind  string
	value interface{}
}

func NewSteps() *steps {
	s := make(steps, 0)
	return &s
}

func (s steps) AddStep(variables []variable, stackTrace []api.Stackframe) *step {
	st := step{}
	_, fileName := path.Split(stackTrace[0].File)
	st.file = stackTrace[0].File
	st.meta.currentPosition.line = stackTrace[0].Line
	st.meta.name = fmt.Sprintf(
		"%v at %v:%v",
		stackTrace[0].Function.Name(),
		fileName,
		stackTrace[0].Line,
	)

	if st.variables == nil {
		st.variables = map[string]variable{}
	}

	for _, val := range variables {
		st.variables[val.name] = val
	}

	s[st.meta.name] = st

	return &st
}
