package onelogin

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// fixtureUniqueToken is replaced in every fixture with a value unique to the
// run. Fixtures in examples/ are documentation as well as test input, so they
// carry readable names — but OneLogin enforces uniqueness on usernames and
// emails within a tenant, and a fixture that hard-codes them passes once and
// fails on every run afterwards with:
//
//	Validation failed: Email must be unique within <tenant>,
//	Username must be unique within <tenant>
//
// A token keeps the examples readable while letting the suite be re-run.
const fixtureUniqueToken = "acctest"

// GetFixture returns the HCL example to be used in an acceptance test, with
// fixtureUniqueToken replaced by a value unique to this run.
func GetFixture(name string, t *testing.T) string {
	t.Helper()
	return getFixtureWithSuffix(name, acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum), t)
}

// GetFixtures returns several fixtures sharing one suffix, for a test whose
// steps must refer to the same resources. Loading them separately would give
// each step a different suffix, and the second step would replace every
// resource rather than updating it.
func GetFixtures(names []string, t *testing.T) []string {
	t.Helper()
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	out := make([]string, 0, len(names))
	for _, name := range names {
		out = append(out, getFixtureWithSuffix(name, suffix, t))
	}
	return out
}

func getFixtureWithSuffix(name, suffix string, t *testing.T) string {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	p := path.Join(filepath.Dir(filename), "../examples", name)

	rawFile, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to load fixture %s for acceptance test: %v", name, err)
	}

	return strings.ReplaceAll(string(rawFile), fixtureUniqueToken, suffix)
}
