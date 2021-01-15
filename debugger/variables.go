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
	v.addByType(vr, "" ,currentLine)
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
	 	name: vr.Name,
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

func (v variables) addStruct(vr api.Variable, key string, currentLine int) {
	if len(vr.Children) == 0 {
		return
	}

	if key == "" {
		key = vr.RealType
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
			v.addByType(child, child.RealType, currentLine)
		} else {
			pointers = append(pointers, child.Name)
			v.addByType(child, child.Name, currentLine)
		}
	}

	rootVar.pointers = pointers
	v[key] = rootVar
}

func (v variables) addSlice(vr api.Variable, key string, currentLine int) {
	if len(vr.Children) == 0 {
		return
	}

	if key == "" {
		// If is a custom type the name should be RealType
		// I think is better to show the full pathname instead of just the name
		if strings.Contains(vr.RealType, ".") {
			key = vr.RealType
		} else {
			key = vr.Name
		}
	}

	rootVar := variable{
		name:  key,
		value: vr.Value,
		kind:  vr.Kind.String(),
	}

	var pointers []string
	realLength := v.getLength(vr)

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

func (v variables) addMap(vr api.Variable, key string, currentLine int) {
	if len(vr.Children) == 0 {
		return
	}

	if key == "" {
		key = vr.Name
	}

	rootVar := variable{
		name:  key,
		value: vr.Value,
		kind:  vr.Kind.String(),
	}

	var pointers []string
	realLength := v.getLength(vr)

	var childKey string

	// We need to double the length of the map because the keys are included
	// in the children array
	for i := 0; i < realLength * 2; i++ {
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
