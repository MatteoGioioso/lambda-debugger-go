package debugger

import (
	"fmt"
	"github.com/go-delve/delve/service/api"
	"lambda-debugger-go/utils"
)

var nonPrimitiveTypes = []string{"struct", "slice", "interface", "func", "map"}

type variable struct {
	name      string
	kind      string
	value     interface{}
	pointers  []string
	hasParent bool
}

type variables map[string]variable

func NewVariables() *variables {
	vs := make(variables, 0)
	return &vs
}

func (v variables) Add(vr api.Variable, currentLine int) {
	if utils.SliceContainsString(nonPrimitiveTypes, vr.Kind.String()) {
		v.addNonPrimitive(vr, vr.RealType)
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
		name:  key,
		value: vr.Value,
		kind:  vr.Kind.String(),
	}

	var pointers []string

	for _, child := range vr.Children {
		kind := child.Kind.String()

		if utils.SliceContainsString(nonPrimitiveTypes, kind) {
			pointers = append(pointers, child.RealType)
			v.addNonPrimitive(child, child.RealType)
		} else {
			pointers = append(pointers, child.Name)
			rawVar := variable{
				name:      child.Name,
				value:     child.Value,
				kind:      kind,
				hasParent: true,
			}
			v[child.Name] = rawVar
		}
	}

	rootVar.pointers = pointers
	v[key] = rootVar
}

func (v variables) addSlice(vr api.Variable, key string) {
	if len(vr.Children) == 0 {
		return
	}

	rootVar := variable{
		name:  key,
		value: vr.Value,
		kind:  vr.Kind.String(),
	}

	var pointers []string
	var realLength int
	maxLen, _ := utils.GetMaxArrayValues()
	if vr.Len > int64(maxLen) {
		realLength = 0
	} else {
		realLength = int(vr.Len)
	}

	for i := 0; i < realLength; i++ {
		child := vr.Children[i]
		kind := child.Kind.String()
		key := child.Name
		if key == "" {
			key = fmt.Sprintf("%v[%d]", child.RealType, i)
		}

		if utils.SliceContainsString(nonPrimitiveTypes, kind) {
			pointers = append(pointers, key)
			v.addNonPrimitive(child, key)
		} else {
			pointers = append(pointers, key)
			rawVar := variable{
				name:      child.Name,
				value:     child.Value,
				kind:      kind,
				hasParent: true,
			}
			v[key] = rawVar
		}
	}

	rootVar.pointers = pointers
	v[key] = rootVar
}
