package onelogin

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

// GetFixture returns the HCL example to be used in an acceptance test
func GetFixture(name string, t *testing.T) string {
	_, filename, _, _ := runtime.Caller(0)
	exPath := filepath.Dir(filename)
	p := path.Join(exPath, "../examples", name)
	file, err := os.Open(p)
	if err != nil {
		t.Fatalf("failed to load fixture for acceptance test")
	}
	rawFile, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to read fixture for acceptance test")
	}
	tfConfig := string(rawFile)
	return tfConfig
}
