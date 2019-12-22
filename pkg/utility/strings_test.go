package utility

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCombineStrings(t *testing.T) {
	expected := "hoge:moge"
	actual := CombineStrings([]string{"hoge:", "moge"}, WithLength(0), WithCapacity(64))

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
}
