package luchen

import (
	"io/fs"
	"net/http"
	"path/filepath"
)

type onlyFilesFS struct {
	root string
	fs   http.FileSystem
}

// Dir 返回 http.FileSystem 实现
// listDirectory = true，返回 http.Dir()
// listDirectory = false，返回 onlyFilesFS，不会显示目录
func Dir(root string, listDirectory bool) http.FileSystem {
	hfs := http.Dir(root)
	if listDirectory {
		return hfs
	}
	return &onlyFilesFS{fs: hfs}
}

// OnlyFilesFS 将 fs.FS 包装为 onlyFilesFS
// root 可以设置根路径，文件将会从根路径开始查找
func OnlyFilesFS(efs fs.FS, root string) http.FileSystem {
	hfs := http.FS(efs)
	return &onlyFilesFS{
		root: root,
		fs:   hfs,
	}
}

// Open conforms to http.Filesystem.
func (ofs onlyFilesFS) Open(name string) (http.File, error) {
	if ofs.root != "" {
		name = filepath.Join(ofs.root, name)
	}
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
