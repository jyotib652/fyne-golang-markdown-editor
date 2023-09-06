package main

import (
	"testing"

	"fyne.io/fyne/test"
)

func Test_makeUI(t *testing.T) {
	var testCfg config

	edit, preview := testCfg.makeUI()

	test.Type(edit, "Hello") // we're using fyne's testing package. Here, we're simulating user clicking edit and writing "Hello" to the file.

	if preview.String() != "Hello" {
		t.Error("Failed -- did not find expected value in preview")
	}
}
