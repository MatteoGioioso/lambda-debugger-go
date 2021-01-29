package collector

import (
	"encoding/json"
	"io/ioutil"
	"lambda-debugger-go/debugger"
	"os"
	"path/filepath"
	"strings"
)

type collector struct{
	outputPath string
	filePath string
}

func New() *collector {
	return &collector{
		outputPath: os.Getenv("LAMBDA_DEBUGGER_OUTPUT_PATH"),
		filePath:   os.Getenv("LAMBDA_DEBUGGER_FILE_PATH"),
	}
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

	htmlBytes, err := ioutil.ReadFile(filepath.Join(c.filePath, "index.html"))
	if err != nil {
		return err
	}

	htmlWithSteps := strings.
		Replace(string(htmlBytes), "//---DEBUG.JSON---//", string(stepsDTOBytes), 1)
	htmlWithStepsAndFiles := strings.
		Replace(htmlWithSteps, "//---FILES.JSON---//", string(filesDTOBytes), 1)
	htmlFileOutPath, err := filepath.Abs(filepath.Join(c.outputPath, "index.html"))

	if err := ioutil.WriteFile(htmlFileOutPath, []byte(htmlWithStepsAndFiles), 0700); err != nil {
		return err
	}

	return nil
}
