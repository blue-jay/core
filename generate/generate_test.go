package generate_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/blue-jay/core/generate"
)

func fileAssertSame(t *testing.T, fileActual, fileExpected string) {
	// Actual output
	actual, err := ioutil.ReadFile(fileActual)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Expected output
	expected, err := ioutil.ReadFile(fileExpected)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Clean up test files
	os.Remove(fileActual)

	// Compare the files
	if string(actual) != string(expected) {
		t.Fatalf("\nactual: %v\nexpected: %v", string(actual), string(expected))
	}
}

// TestSingle ensures single file can be generated.
func TestSingle(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"
	file := "model/foo/foo.go"
	fileActual := filepath.Join(actualFolder, file)
	fileExpected := filepath.Join(expectedFolder, file)

	// Set the arguments
	args := []string{
		"single/default",
		"package:foo",
		"table:bar",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual, fileExpected)
}

// TestSingleMissing ensures single file fails on a missing value.
func TestSingleMissing(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	file := "model/foo/foo.go"
	fileActual := filepath.Join(actualFolder, file)

	// Set the arguments
	args := []string{
		"single/default",
		"package:foo",
		//"table:bar",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err == nil {
		t.Fatalf("%v", err)
	}
}

// TestSingleNoParse ensures single file can be generated without parsing.
func TestSingleNoParse(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"
	file := "view/foo/index.tmpl"
	fileActual := filepath.Join(actualFolder, file)
	fileExpected := filepath.Join(expectedFolder, file)

	// Set the arguments
	args := []string{
		"single/noparse",
		"model:foo",
	}

	// Clear out files from old tests
	os.Remove(fileActual)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual, fileExpected)
}

// TestSingle ensures a collection can be generated.
func TestCollection(t *testing.T) {
	// Set the variables
	templateFolder := "testdata/generate"
	actualFolder := "testdata/actual"
	expectedFolder := "testdata/expected"

	// Set the arguments
	args := []string{
		"collection/default",
		"model:modeltest",
		"package:packagetest",
		"view:viewtest",
	}

	// Info for file 1
	file1 := "model/packagetest/packagetest.go"
	fileActual1 := filepath.Join(actualFolder, file1)
	fileExpected1 := filepath.Join(expectedFolder, file1)
	os.Remove(fileActual1)

	// Info for file 2
	file2 := "view/viewtest/index.tmpl"
	fileActual2 := filepath.Join(actualFolder, file2)
	fileExpected2 := filepath.Join(expectedFolder, file2)
	os.Remove(fileActual2)

	// Generate the code
	err := generate.Run(args, actualFolder, templateFolder)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Ensure the files are the same
	fileAssertSame(t, fileActual1, fileExpected1)
	fileAssertSame(t, fileActual2, fileExpected2)
}
