package mcfile

import "io/fs"

// mcfile.Contentity (even including those representing directories)
// should implement several methods:
// -- three for fs.File with the exact same method signatures -
// Stat() (FileInfo, error), Read([]byte) (int, error), Close() error
// -- four as zero-argument methods on Contentity that method calls on
// ContentityFS delegate can to: fs.FS.Open(path), fs.StatFS.Stat(),
// fs.ReadFileFS.ReadFile(path), fs.ReadDirFS.ReadDir(path).
// -- one method that is expected of an fs.File that is a directory:
// ReadDir(n int) ([]DirEntry, error)
// Then a Contentity can be treated like an fs.File, and the actual
// file system does not have to be accessed again.

/*
fs.ErrInvalid    "invalid argument"
fs.ErrPermission "permission denied"
fs.ErrExist      "file already exists"
fs.ErrNotExist   "file does not exist"
fs.ErrClosed     "file already closed"
*/

// ===================
//  IMPLEMENT fs.File
//   (copying os.File)
// ===================

// Stat is Stat.
func (p *Contentity) Stat() (FI fs.FileInfo, patherror error) {
	return nil, nil
}

// Read reads up to len(b) bytes from the File. It
// returns the number of bytes read and any error
// encountered. At EOF, Read returns (0,io.EOF).
func (p *Contentity) Read([]byte) (int, error) {
	return 0, nil
}

// Close closes the File, rendering it unusable for I/O.
// On files that support SetDeadline, any pending I/O
// operations are canceled and return at once with error.
// Close returns an error if it has already been called.
func (p *Contentity) Close() error {
	return nil
}

// ReadDir reads the contents of the directory associated
// with the file and returns a slice of DirEntry values in
// directory order. Subsequent calls on the same file will
// yield later DirEntry records in the directory.
// If n > 0,
// ReadDir returns at most n DirEntry records. If it returns
// an empty slice, it will return an error explaining why.
// At the end of a directory, the error is io.EOF.
// If n <= 0,
// ReadDir returns all the DirEntry records remaining in the
// directory. On succeeds, it returns a nil error (not io.EOF).
func (p *Contentity) ReadDir(n int) ([]fs.DirEntry, error) {
	return nil, nil
}

// ===================
//  Sort-of-IMPLEMENT
//   fs.(various)FS
// ===================

// Open returns Contentity or nil, and is
// fs.FS.Open(path)(fs.File,fs.patherror) implements fs.FS ;
// patherror is { Op: "open", Path: path, Err: "theReason" },
// e.g. ErrNotExist, or ErrInvalid if fails fs.ValidPath(name)
func (p *ContentityFS) Open(path string) (F fs.File, patherror error) {
	return nil, nil
}

// Stat is fs.FS.Stat(path)(fs.FileInfo,fs.patherror) implements fs.StatFS
func (p *ContentityFS) Stat(path string) (FI fs.FileInfo, patherror error) {
	return nil, nil
}

// ReadFile is fs.FS.ReadFile(path)([]byte,error) implements fs.ReadFileFS ;
// Note that success returns nil error, not io.EOF error.
func (p *ContentityFS) ReadFile(path string) ([]byte, error) {
	return nil, nil
}

/*
type DirEntry interface {
	// Name returns only the final element of the path (the base name).
	Name() string
	IsDir() bool
	Type() FileMode
	// Info returns a FileInfo, which may be from the time of the
	// original directory read or from the time of the call to Info.
	// If the file has been removed/renamed since the directory read,
	// Info may return an error such that errors.Is(err, ErrNotExist).
	// If the entry denotes a symbolic link, Info reports about the
	// link itself, not the link's target.
	Info() (FileInfo, error)
}

type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // nr of bytes for regular files; else system-dep
	Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}
*/
