package utils

import (
	"io/ioutil"
	"os"
	"path"
)

func ClearFolder() {
	//Clear folder output after upload to S3

	dir, err := ioutil.ReadDir("output")
	if err != nil {
		panic(err)
	}
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"output", d.Name()}...))
	}
}
