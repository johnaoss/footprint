package footprint_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/johnaoss/footprint"
)

var (
	dir = os.TempDir()
)

// TestAdd tests if we can properly add a file to the List.
func TestAdd(t *testing.T) {
	// Initialize new footprint file
	fp := new(footprint.List)

	t.Log("Initial Footprint:", fp.String())

	// Detect if after adding a new file, we can see that change reflected in
	// the footprint.
	exFile := newfile(t)
	defer os.Remove(exFile.Name())

	en, err := footprint.MakeEntry(exFile)
	if err != nil {
		t.Errorf("failed to make entry: %v", err)
		return
	}

	fp.Add(en)
	if fp.Len() != 1 {
		t.Errorf("length should be 1, instead was: %d", fp.Len())
	}

	t.Log("Updated Footprint:", strings.Trim(fp.String(), "\n"))
}

func TestParse(t *testing.T) {
	fp := &footprint.List{}
	entryFile := newfile(t)
	defer os.Remove(entryFile.Name())

	entry, err := footprint.MakeEntry(entryFile)
	if err != nil {
		t.Errorf("failed to make entry: %v", err)
		return
	}
	fp.Add(entry)

	f := newfile(t)
	defer os.Remove(f.Name())

	info, err := f.Stat()
	if err != nil {
		t.Error("failed to get file info")
		return
	}

	t.Log("Old Footprint:", strings.Trim(fp.String(), "\n"))

	if err := ioutil.WriteFile(f.Name(), []byte(fp.String()), info.Mode().Perm()); err != nil {
		t.Error(err)
		return
	}

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Error(err)
		return
	}

	newList, err := footprint.Parse(bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("failed to parse footprint: %v", err)
	}

	if newList.Len() != fp.Len() {
		t.Errorf("lists are of differing sizes, expected: %d, given: %d", fp.Len(), newList.Len())
	}

	t.Log("New footprint:", strings.Trim(newList.String(), "\n"))
}

func ExampleList() {
	// The zero value of the footprint is valid.
	fp := new(footprint.List)

	// We now create a temporary file, and parse it into an Entry
	file, _ := os.Create("tempfile")
	defer os.Remove("tempfile")

	entry, err := footprint.MakeEntry(file)
	if err != nil {
		panic(err)
	}

	// Now we can add it to the list
	fp.Add(entry)
}

// newfile is a helper that creates a temporary file.
func newfile(t *testing.T) *os.File {
	t.Helper()
	f, err := ioutil.TempFile(dir, "file")
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
	}

	t.Log("Created file:", f.Name())
	return f
}
