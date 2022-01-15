package utils

import "os"

func ProjectRoot() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}

func AbsolutPath(relativePath string) string {
	return ProjectRoot() + relativePath
}
