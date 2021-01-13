package debugger

import (
	"github.com/go-delve/delve/service/api"
	"lambda-debugger-go/utils"
)

var nonPrimitiveTypes = []string{"struct", "slice", "interface", "func"}

type variable struct {
	name     string
	kind     string
	value    interface{}
	pointers []string
}

type variables map[string]variable

func NewVariables() *variables {
	vs := make(variables, 0)
	return &vs
}

func (v variables) Add(vr api.Variable, currentLine int) {
	if utils.SliceContainsString(nonPrimitiveTypes, vr.Kind.String()) {
		v.addNonPrimitive(vr, vr.Name)
	} else {
		v.addPrimitive(vr, currentLine)
	}
}

func (v variables) addPrimitive(vr api.Variable, currentLine int) {
	if int(vr.DeclLine) == currentLine {
		return
	}

	v[vr.Name] = variable{
		name:  vr.Name,
		kind:  vr.Kind.String(),
		value: vr.Value,
	}
}

func (v variables) addNonPrimitive(vr api.Variable, key string) {
	switch vr.Kind.String() {
	case "struct":
		v.addStruct(vr, key)
	case "slice":
		v.addSlice(vr, key)
	}
}

func (v variables) addStruct(vr api.Variable, key string) {
	if len(vr.Children) == 0 {
		return
	}

	rootVar := variable{
		name:  vr.Name,
		value: vr.Value,
		kind:  vr.Kind.String(),
	}

	var pointers []string

	for _, child := range vr.Children {
		kind := child.Kind.String()

		if utils.SliceContainsString(nonPrimitiveTypes, kind) {
			pointers = append(pointers, child.Name)
			v.addNonPrimitive(child, child.Name)
		} else {
			pointers = append(pointers, child.Name)
			rawVar := variable{
				name:  child.Name,
				value: child.Value,
				kind:  kind,
			}
			v[child.Name] = rawVar
		}
	}

	rootVar.pointers = pointers
	v[key] = rootVar
}

func (v variables) addSlice(vr api.Variable, key string) {

}
