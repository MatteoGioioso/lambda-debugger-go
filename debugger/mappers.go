package debugger

type StepsDTO []StepDTO
type StepDTO struct {
	Meta struct {
		CurrentPosition struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"currentPosition"`
		Name string `json:"name"`
	} `json:"meta"`
	File      string `json:"file"`
	Variables map[string]struct {
		Name  string      `json:"name"`
		Kind  string      `json:"kind"`
		Value interface{} `json:"value"`
	} `json:"variables"`
}

func ToStepsDTO(ss steps) StepsDTO {
	stepsDTO := make(StepsDTO, 0)
	for _, stepValue := range ss {
		stepDTO := StepDTO{
			Meta: struct {
				CurrentPosition struct {
					Line   int `json:"line"`
					Column int `json:"column"`
				} `json:"currentPosition"`
				Name string `json:"name"`
			}{
				CurrentPosition: struct {
					Line   int `json:"line"`
					Column int `json:"column"`
				}{
					Line:   stepValue.meta.currentPosition.line,
					Column: stepValue.meta.currentPosition.column,
				},
				Name: stepValue.meta.name,
			},
			File:      stepValue.file,
			Variables: make(map[string]struct {
				Name  string      `json:"name"`
				Kind  string      `json:"kind"`
				Value interface{} `json:"value"`
			}, 0),
		}

		for varKey, varValue := range stepValue.variables {
			stepDTO.Variables[varKey] = struct {
				Name  string      `json:"name"`
				Kind  string      `json:"kind"`
				Value interface{} `json:"value"`
			}{Name: varValue.name, Kind: varValue.kind, Value: varValue.value}
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
