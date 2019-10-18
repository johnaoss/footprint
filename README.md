# footprint

footprint is a package that allows a central form of managing files in an application.

The actual structure of this project is inspired by CRUX's port system, and the `.footprint` file within.

Currently I've only implemented adding files to the tracker and formatting them accordingly, the rest is TBD depending on exams.

## Installation

```
go get github.com/johnaoss/footprint
```

## Usage

```go

package main

import (
    "os"
    "fmt"

    "github.com/johnaoss/footprint"
)

func main() {
	// The zero value of the footprint is valid.
	fp := new(footprint.List)

	// We now create a temporary file, and parse it into an Entry
	file, _ := os.Create("tempfile")
	defer os.Remove("tempfile")

	entry, err := footprint.MakeEntry(file)
	if err != nil {
        fmt.Println(err)
        return
	}

	// Now we can add it to the list
    fp.Add(entry)

    // Verify output
    fmt.Println(fp.String())
}
```

## Roadmap

- [ ] Tracking files added
- [ ] Validating existing files
- [ ] Writing to disk
- [ ] Advanced file management
- [ ] Hard file verification (md5 hash?)


## License

It's MIT. See the LICENSE.md file for details.

