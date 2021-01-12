package collector

import (
	"encoding/json"
	"io/ioutil"
	"lambda-debugger-go/debugger"
	"path/filepath"
	"strings"
)

type collector struct{}

func New() *collector {
	return &collector{}
}

func (c collector) InjectDebuggerOutputIntoHtml(steps debugger.StepsDTO, files debugger.FilesDTO) error {
	stepsDTOBytes, err := json.MarshalIndent(steps, "", "    ")
	if err != nil {
		return err
	}
	filesDTOBytes, err := json.MarshalIndent(files, "", "    ")
	if err != nil {
		return err
	}

	htmlFilePath, err := filepath.Abs("../index.html")
	if err != nil {
		return err
	}

	htmlBytes, err := ioutil.ReadFile(htmlFilePath)
	if err != nil {
		return err
	}

	htmlWithSteps := strings.
		Replace(string(htmlBytes), "//---DEBUG.JSON---//", string(stepsDTOBytes), 1)
	htmlWithStepsAndFiles := strings.
		Replace(htmlWithSteps, "//---FILES.JSON---//", string(filesDTOBytes), 1)
	htmlFileOutPath, err := filepath.Abs("../tmp/index.html")

	if err := ioutil.WriteFile(htmlFileOutPath, []byte(htmlWithStepsAndFiles), 0700); err != nil {
		return err
	}

	return nil
}
