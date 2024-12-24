package assert

import (
	"os"
	"reflect"
	"slices"
	"testing"
)

func Strings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func Ints(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Bools(t testing.TB, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func StringSlices(t testing.TB, got, want []string) {
	t.Helper()

	slices.Sort(got)
	slices.Sort(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func Errors(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func Nils(t testing.TB, got interface{}) {
	t.Helper()

	if got == nil {
		t.Errorf("got %q want nil", got)
	}
}

func FileExists(t testing.TB, path string) {
	_, err := os.ReadFile(path)

	if err != nil {
		t.Errorf("File %s doesn't exist: Error %s", path, err)
	}
}
