package random

import (
	"reflect"
	"testing"
)

func TestRandomString(t *testing.T) {
	gotRandomString := NewRandomString("p_", 7)

	wantRandomString := len("p_") + 7

	if !reflect.DeepEqual(len(gotRandomString), wantRandomString) {
		t.Fatalf("Test №1 - Not equeal got: %d && want: %d", len(gotRandomString), wantRandomString)
	}

	gotRandomString = NewRandomString("i_", 5)

	wantRandomString = len("i_") + 7

	if !reflect.DeepEqual(len(gotRandomString), wantRandomString) {
		t.Fatalf("Test №2 - Not equeal got: %d && want: %d", len(gotRandomString), wantRandomString)
	}

}
