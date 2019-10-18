// Package footprint maintains a footprint of all the entries used in an
// application.
//
// This allows centralized management of entries, and eventually will allow a method
// to remove all entries associated with this program from a user's disk. Essentially
// acting as a complete uninstall, instead of leaving leftover entries.
package footprint

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

// List represents a list of entries that a program may use.
// This allows the central handling of all temporary entries a program may create.
// The zero value of this is safe to use.
type List struct {
	// mutex to allow concurrent access
	mu sync.RWMutex

	// entries is a list of all the entries currently being watched in the footprint
	entries []Entry
}

// Parse reads in a List from a given reader.
func Parse(r io.Reader) (*List, error) {
	buf := bufio.NewReader(r)

	var entries []Entry

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed to read line: %w", err)
			}
			break
		}

		parts := bytes.Split(line, []byte{'\t'})
		if len(parts) != 3 {
			return nil, fmt.Errorf("line improperly seperated: %s", line)
		}

		perms, err := parsePerms(parts[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse permission bits: %w", err)
		}

		owners := bytes.Split(parts[1], []byte{'/'})

		entries = append(entries, Entry{
			Perms: perms,
			Path:  string(parts[2]),
			Uid:   string(owners[0]),
			Gid:   string(owners[1]),
		})
	}

	return &List{entries: entries}, nil
}

// Len returns the current number of entries stored within the List.
func (l *List) Len() int {
	l.mu.RLock()
	num := len(l.entries)
	l.mu.RUnlock()
	return num
}

// Add registers a file with the List.
func (l *List) Add(e Entry) {
	l.mu.Lock()
	l.entries = append(l.entries, e)
	l.mu.Unlock()
}

// Create wraps os.Create in order to add the file as an entry in the
// middle of the creation step.
func (l *List) Create(name string) (*os.File, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	entry, err := MakeEntry(file)
	if err != nil {
		return nil, err
	}
	l.Add(entry)

	return file, nil
}

// Validate verifies that all entries have information accurate to what is stored
// in the List.
// Currently unimplemented
func (l *List) Validate() error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return nil
}

// Remove deletes an entry from the List.
// Currently unimplemented.
func (l *List) Remove(name string) error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return nil
}

// String writes out a representation of the List to a string.
func (l *List) String() string {
	if l.Len() == 0 {
		return ""
	}
	buf := new(bytes.Buffer)
	if err := l.write(buf); err != nil {
		return err.Error()
	}
	return buf.String()
}

// write outputs the representation of the file to an io.Writer.
func (l *List) write(w io.Writer) error {
	b := new(bytes.Buffer)
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, elem := range l.entries {
		b.WriteString(elem.String())
		b.WriteByte('\n')
		if _, err := w.Write(b.Bytes()); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		b.Reset()
	}

	return nil
}

// parsePerms parses the Unix permission bits from a string.
func parsePerms(s []byte) (os.FileMode, error) {
	if len(s) != len("drwxr-xr-x") {
		return 0, fmt.Errorf("invalid len of given string")
	}

	var sum int
	for i := 1; i < len(s); i += 3 {
		var count int
		for j := 0; j < 3; j++ {
			switch s[i+j] {
			case 'r':
				count += 4
			case 'w':
				count += 2
			case 'x':
				count += 1
			case '-':
			default:
				return 0, fmt.Errorf("invalid char: %d", s[i+j])
			}
		}
		sum += count * (1 << (6 - (i - 1)))
	}

	return os.FileMode(sum), nil
}
