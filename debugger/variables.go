package debugger

import (
	"fmt"
	"github.com/go-delve/delve/service/api"
	"lambda-debugger-go/utils"
	"strings"
)

var nonPrimitiveTypes = []string{"struct", "slice", "interface", "func", "map", "array"}

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
	v.addByType(vr, "", currentLine)
}

func (v variables) addByType(vr api.Variable, key string, currentLine int) {
	switch vr.Kind.String() {
	case "struct":
		v.addStruct(vr, key, currentLine)
	case "slice":
		v.addSlice(vr, key, currentLine)
	case "map":
		v.addMap(vr, key, currentLine)
	case "array":
		v.addSlice(vr, key, currentLine)
	default:
		v.addPrimitive(vr, key, currentLine)
	}
}

func (v variables) addPrimitive(vr api.Variable, key string, currentLine int) {
	if int(vr.DeclLine) == currentLine {
		return
	}

	primitive := variable{
		name:  vr.Name,
		kind:  vr.Kind.String(),
		value: vr.Value,
	}

	if key != "" {
		primitive.hasParent = true
	} else {
		key = vr.Name
	}

	v[key] = primitive
}

func (v variables) addStruct(vr api.Variable, initKey string, currentLine int) {
	if len(vr.Children) == 0 {
		return
	}

	rootVar, key := v.initVariable(vr, initKey)

	var pointers []string

	for _, child := range vr.Children {
		kind := child.Kind.String()

		if utils.SliceContainsString(nonPrimitiveTypes, kind) {
			pointers = append(pointers, child.Name)
			v.addByType(child, child.Name, currentLine)
		} else {
			pointers = append(pointers, child.Name)
			v.addByType(child, child.Name, currentLine)
		}
	}

	rootVar.pointers = pointers
	v[key] = rootVar
}

func (v variables) addSlice(vr api.Variable, initKey string, currentLine int) {
	defer func() {
		// Current workaround for delve bug: sometimes it returns incorrect length of slice
		if err := recover(); err != nil {
			errorConcrete := err.(error)
			if !strings.Contains(errorConcrete.Error(), "runtime error: index out of range") {
				panic(err)
			}
		}
	}()

	if len(vr.Children) == 0 {
		return
	}

	rootVar, key := v.initVariable(vr, initKey)
	realLength := v.getLength(vr)

	var pointers []string

	for i := 0; i < realLength; i++ {
		child := vr.Children[i]
		kind := child.Kind.String()
		key := child.Name
		if key == "" {
			key = fmt.Sprintf("%v[%d]", child.RealType, i)
		}

		if utils.SliceContainsString(nonPrimitiveTypes, kind) {
			pointers = append(pointers, key)
			v.addByType(child, key, currentLine)
		} else {
			pointers = append(pointers, key)
			v.addByType(child, key, currentLine)
		}
	}

	rootVar.pointers = pointers
	v[key] = rootVar
}

func (v variables) addMap(vr api.Variable, initKey string, currentLine int) {
	if len(vr.Children) == 0 {
		return
	}

	rootVar, key := v.initVariable(vr, initKey)
	realLength := v.getLength(vr)

	var pointers []string
	var childKey string

	// We need to double the length of the map because the keys are included
	// in the children array
	for i := 0; i < realLength*2; i++ {
		// In Map's children the first value is the key of the map
		// hence even indexes are going to be keys and odd values
		if i%2 == 0 {
			childKey = vr.Children[i].Value
		} else {
			child := vr.Children[i]
			kind := child.Kind.String()

			if utils.SliceContainsString(nonPrimitiveTypes, kind) {
				pointers = append(pointers, childKey)
				v.addByType(child, childKey, currentLine)
			} else {
				pointers = append(pointers, childKey)
				v.addByType(child, childKey, currentLine)
			}
		}
	}

	rootVar.pointers = pointers
	v[key] = rootVar
}

func (v variables) getLength(vr api.Variable) int {
	var realLength int
	maxLen, _ := utils.GetMaxArrayValues()
	if vr.Len > int64(maxLen) {
		realLength = 0
	} else {
		realLength = int(vr.Len)
	}

	return realLength
}

func (v variables) initVariable(vr api.Variable, key string) (variable, string) {
	rootVar := variable{
		value: vr.Value,
		kind:  vr.Kind.String(),
	}

	if key == "" {
		key = vr.Name
		rootVar.name = key
	} else {
		rootVar.name = key
		rootVar.hasParent = true
	}

	return rootVar, key
}
