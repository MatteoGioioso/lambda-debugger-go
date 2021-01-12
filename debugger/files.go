package debugger

import "io/ioutil"

type files map[string]file
type file struct {
	code string
}

func NewFiles() *files {
	fs := make(files, 0)
	return &fs
}

func (f files) Add(fileUrl string) error {
	if _, ok := f[fileUrl]; ok {
		return nil
	}

	fl := file{}
	fileBinary, err := ioutil.ReadFile(fileUrl)
	if err != nil {
		return err
	}

	fl.code = string(fileBinary)

	f[fileUrl] = fl

	return nil
}
