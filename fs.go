package luchen

import (
	"io/fs"
	"net/http"
)

type onlyFilesFS struct {
	fs http.FileSystem
}

// Dir 返回 http.FileSystem 实现
// listDirectory = true，返回 http.Dir()
// listDirectory = false，返回 onlyFilesFS，不会显示目录
func Dir(root string, listDirectory bool) http.FileSystem {
	hfs := http.Dir(root)
	if listDirectory {
		return hfs
	}
	return &onlyFilesFS{hfs}
}

// Open conforms to http.Filesystem.
func (ofs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := ofs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, fs.ErrNotExist
	}
	return f, nil
}
