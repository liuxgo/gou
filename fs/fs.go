package fs

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/yaoapp/gou/connector"
	"github.com/yaoapp/gou/fs/system"
	"github.com/yaoapp/kun/exception"
)

// FileSystems Register filesystems
var FileSystems = map[string]FileSystem{
	"system": system.New(),
}

// RootFileSystems high-level filesystem
var RootFileSystems = map[string]FileSystem{}

// RegisterConnector register a fileSystem via connector
func RegisterConnector(c connector.Connector) error {
	// if c.Is(connector.DATABASE) {
	// 	FileSystems[c.ID()] = system.New() // xun.New(Connector)

	// } else if c.Is(connector.REDIS) {
	// 	FileSystems[c.ID()] = system.New() // redis.New(Connector)

	// } else if c.Is(connector.MONGO) {
	// 	FileSystems[c.ID()] = system.New() // mongo.New(Connector)
	// }
	return fmt.Errorf("connector %s does not support", c.ID())
}

// Register a FileSystem
func Register(id string, fs FileSystem) FileSystem {
	FileSystems[id] = fs
	return FileSystems[id]
}

// RootRegister Register a root FileSystem
func RootRegister(id string, fs FileSystem) FileSystem {
	RootFileSystems[id] = fs
	return RootFileSystems[id]
}

// Get pick a filesystem via the given name
func Get(name string) (FileSystem, error) {
	if fs, has := FileSystems[name]; has {
		return fs, nil
	}
	return nil, fmt.Errorf("%s does not registered", name)
}

// MustGet pick a filesystem via the given name
func MustGet(name string) FileSystem {
	fs, err := Get(name)
	if err != nil {
		exception.New(err.Error(), 400).Throw()
		return nil
	}
	return fs
}

// RootGet pick a filesystem via the given name (root first)
func RootGet(name string) (FileSystem, error) {

	if fs, has := RootFileSystems[name]; has {
		return fs, nil
	}

	if fs, has := FileSystems[name]; has {
		return fs, nil
	}

	return nil, fmt.Errorf("%s does not registered", name)
}

// MustRootGet pick a filesystem via the given name
func MustRootGet(name string) FileSystem {
	fs, err := RootGet(name)
	if err != nil {
		exception.New(err.Error(), 400).Throw()
		return nil
	}
	return fs
}

// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile reads the whole file, it does not treat an EOF from Read as an error to be reported.
func ReadFile(xfs FileSystem, file string) ([]byte, error) {
	return xfs.ReadFile(file)
}

// WriteFile writes data to the named file, creating it if necessary.
//
//	If the file does not exist, WriteFile creates it with permissions perm (before umask); otherwise WriteFile truncates it before writing, without changing permissions.
func WriteFile(xfs FileSystem, file string, data []byte, perm uint32) (int, error) {
	return xfs.WriteFile(file, data, perm)
}

// ReadDir reads the named directory, returning all its directory entries sorted by filename.
// If an error occurs reading the directory, ReadDir returns the entries it was able to read before the error, along with the error.
func ReadDir(xfs FileSystem, dir string, recursive bool) ([]string, error) {
	return xfs.ReadDir(dir, recursive)
}

// Mkdir creates a new directory with the specified name and permission bits (before umask).
// If there is an error, it will be of type *PathError.
func Mkdir(xfs FileSystem, dir string, perm uint32) error {
	return xfs.Mkdir(dir, perm)
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
// The permission bits perm (before umask) are used for all directories that MkdirAll creates. If path is already a directory, MkdirAll does nothing and returns nil.
func MkdirAll(xfs FileSystem, dir string, perm uint32) error {
	return xfs.MkdirAll(dir, perm)
}

// MkdirTemp creates a new temporary directory in the directory dir and returns the pathname of the new directory.
// The new directory's name is generated by adding a random string to the end of pattern.
// If pattern includes a "*", the random string replaces the last "*" instead. If dir is the empty string, MkdirTemp uses the default directory for temporary files, as returned by TempDir.
// Multiple programs or goroutines calling MkdirTemp simultaneously will not choose the same directory. It is the caller's responsibility to remove the directory when it is no longer needed.
func MkdirTemp(xfs FileSystem, dir string, pattern string) (string, error) {
	return xfs.MkdirTemp(dir, pattern)
}

// Chmod changes the mode of the named file to mode. If the file is a symbolic link, it changes the mode of the link's target. If there is an error, it will be of type *PathError.
// A different subset of the mode bits are used, depending on the operating system.
// On Unix, the mode's permission bits, ModeSetuid, ModeSetgid, and ModeSticky are used.
// On Windows, only the 0200 bit (owner writable) of mode is used; it controls whether the file's read-only attribute is set or cleared. The other bits are currently unused.
// For compatibility with Go 1.12 and earlier, use a non-zero mode. Use mode 0400 for a read-only file and 0600 for a readable+writable file.
// On Plan 9, the mode's permission bits, ModeAppend, ModeExclusive, and ModeTemporary are used.
func Chmod(xfs FileSystem, name string, mode uint32) error {
	return xfs.Chmod(name, mode)
}

// Remove removes the named file or (empty) directory. If there is an error, it will be of type *PathError.
func Remove(xfs FileSystem, name string) error {
	return xfs.Remove(name)
}

// RemoveAll removes path and any children it contains. It removes everything it can but returns the first error it encounters. If the path does not exist, RemoveAll returns nil (no error). If there is an error, it will be of type *PathError.
func RemoveAll(xfs FileSystem, name string) error {
	return xfs.RemoveAll(name)
}

// Move move from src to dst
func Move(xfs FileSystem, name string, dst string) error {
	return xfs.Move(name, dst)
}

// Copy copy from src to dst
func Copy(xfs FileSystem, name string, dst string) error {
	return xfs.Copy(name, dst)
}

// Exists returns a boolean indicating whether the error is known to report that a file or directory already exists.
// It is satisfied by ErrExist as well as some syscall errors.
func Exists(xfs FileSystem, name string) (bool, error) {
	return xfs.Exists(name)
}

// Size return the length in bytes for regular files; system-dependent for others
func Size(xfs FileSystem, name string) (int, error) {
	return xfs.Size(name)
}

// Mode return the file mode bits
func Mode(xfs FileSystem, name string) (uint32, error) {
	return xfs.Mode(name)
}

// ModTime return the file modification time
func ModTime(xfs FileSystem, name string) (time.Time, error) {
	return xfs.ModTime(name)
}

// IsDir check the given path is dir
func IsDir(xfs FileSystem, name string) bool {
	return xfs.IsDir(name)
}

// IsFile check the given path is file
func IsFile(xfs FileSystem, name string) bool {
	return xfs.IsFile(name)
}

// IsLink check the given path is symbolic link
func IsLink(xfs FileSystem, name string) bool {
	return xfs.IsLink(name)
}

// MimeType return the MimeType
func MimeType(xfs FileSystem, name string) (string, error) {
	return xfs.MimeType(name)
}

// BaseName return the base name
func BaseName(name string) string {
	return filepath.Base(name)
}

// DirName return the dir name
func DirName(name string) string {
	return filepath.Dir(name)
}

// ExtName return the extension name
func ExtName(name string) string {
	return strings.TrimPrefix(filepath.Ext(name), ".")
}
