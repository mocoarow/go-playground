package main_test

import "testing"

func Test1(t *testing.T) {
	t.Parallel()
	if 1*2 != 2 {
		t.Errorf("2 is not equal to 2")
	}
}
