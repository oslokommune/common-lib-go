package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

func checkIfModuleRoot(path string) bool {
	_, err := os.ReadFile(filepath.Join(path, "go.mod"))
	return err == nil
}

func getModulePath() string {
	currPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	found := checkIfModuleRoot(currPath)
	for !found && len(currPath) > 1 {
		currPath = filepath.Dir(currPath)
		found = checkIfModuleRoot(currPath)
	}

	if !found {
		panic("Not a module")
	}

	return currPath
}

func LoadResourceFile(t *testing.T, filename string) []byte {
	basePath := getModulePath()
	resourcePath := filepath.Join(basePath, filename)
	content, err := os.ReadFile(resourcePath)
	if err != nil {
		t.Fatal(err)
	}

	return content
}
