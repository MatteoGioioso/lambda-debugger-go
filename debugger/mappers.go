package debugger

type StepsDTO []StepDTO
type StepDTO struct {
	Meta      Meta      `json:"meta"`
	File      string    `json:"file"`
	Variables Variables `json:"variables"`
}
type Meta struct {
	CurrentPosition CurrentPosition `json:"currentPosition"`
	Name            string          `json:"name"`
}
type CurrentPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
type Variables map[string]Variable
type Variable struct {
	Name      string      `json:"name"`
	Kind      string      `json:"kind"`
	Value     interface{} `json:"value"`
	Pointers  []string    `json:"pointers"`
	HasParent bool        `json:"hasParent"`
}

func ToStepsDTO(ss steps) StepsDTO {
	stepsDTO := make(StepsDTO, 0)
	for _, stepValue := range ss {
		stepDTO := StepDTO{
			Meta: Meta{
				CurrentPosition: CurrentPosition{
					Line:   stepValue.meta.currentPosition.line,
					Column: stepValue.meta.currentPosition.column,
				},
				Name: stepValue.meta.name,
			},
			File:      stepValue.file,
			Variables: make(Variables, 0),
		}

		for varKey, varValue := range stepValue.variables {
			stepDTO.Variables[varKey] = Variable{
				Name:     varValue.name,
				Kind:     varValue.kind,
				Value:    varValue.value,
				Pointers: varValue.pointers,
				HasParent: varValue.hasParent,
			}
		}

		stepsDTO = append(stepsDTO, stepDTO)
	}

	return stepsDTO
}

type FilesDTO map[string]FileDTO
type FileDTO struct {
	Code string `json:"code"`
}

func ToFilesDTO(fs files) FilesDTO {
	filesDTO := make(FilesDTO, 0)
	for key, value := range fs {
		filesDTO[key] = FileDTO{
			Code: value.code,
		}
	}

	return filesDTO
}
