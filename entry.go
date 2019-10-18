package footprint

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

// Entry represents an entry of a file into our List.
type Entry struct {
	// Perms is the permissions this entry has.
	Perms os.FileMode
	// Path is the absolute path of this entry
	Path string
	// Uid is the file's user owner, not the UID number.
	Uid string
	// Gid is the file's group owner, not the GID number.
	Gid string
}

// MakeEntry returns an entry parsed from a given file.
func MakeEntry(f *os.File) (Entry, error) {
	info, err := os.Stat(f.Name())
	if err != nil {
		return Entry{}, err
	}
	path, err := filepath.Abs(f.Name())
	if err != nil {
		return Entry{}, err
	}

	uid, err := user.LookupId(strconv.FormatUint(uint64(info.Sys().(*syscall.Stat_t).Uid), 10))
	if err != nil {
		return Entry{}, err
	}

	gid, err := user.LookupGroupId(strconv.FormatUint(uint64(info.Sys().(*syscall.Stat_t).Gid), 10))
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		Path:  path,
		Perms: info.Mode().Perm(),
		Uid:   uid.Username,
		Gid:   gid.Name,
	}, nil
}

// String returns the representation of the Entry as a tab-delimited format.
// This format is similar to that of the CRUX operating system's `.footprint`
// file format. An example of that is found here: https://github.com/mikek/crux-mike/blob/master/mc/.footprint
func (e *Entry) String() string {
	return fmt.Sprintf("%s\t%s/%s\t%s", e.Perms.String(), e.Uid, e.Gid, e.Path)
}
