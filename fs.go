package luchen

import (
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const indexPage = "index.html"

type onlyFilesFS struct {
	root string
	fs   fs.FS
}

// Dir 返回 http.FileSystem 实现
// listDirectory = true，返回 http.Dir()
// listDirectory = false，返回 onlyFilesFS，不会显示目录
func Dir(root string, listDirectory bool) fs.FS {
	dfs := os.DirFS(root)
	if listDirectory {
		return dfs
	}
	return OnlyFilesFS(dfs, listDirectory, "")
}

// OnlyFilesFS 将 fs.FS 包装为 onlyFilesFS
// root 可以设置根路径，文件将会从根路径开始查找
func OnlyFilesFS(fs fs.FS, listDirectory bool, root string) fs.FS {
	if listDirectory {
		return fs
	}
	ofs := &onlyFilesFS{
		root: root,
		fs:   fs,
	}
	return ofs
}

// Open conforms to http.Filesystem.
func (ofs onlyFilesFS) Open(name string) (fs.File, error) {
	fname := name
	if ofs.root != "" {
		fname = filepath.Join(ofs.root, name)
	}
	f, err := ofs.fs.Open(fname)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		// 访问目录自动定向到 index.html
		return ofs.Open(filepath.Join(name, indexPage))
	}
	return f, nil
}

// FileHandler 处理静态文件请求
// 参考 http.StripPrefix
func FileHandler(prefix string, fs fs.FS) http.Handler {
	h := http.FileServerFS(fs)
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			upath := r.URL.Path
			if !strings.HasPrefix(upath, "/") {
				upath = "/" + upath
				r.URL.Path = upath
			}
			if upath == "/" {
				// 这里是为了避免 http fileHandler 重定向，导致访问错误
				http.ServeFileFS(w, r, fs, path.Clean(upath))
				return
			}
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			r2.URL.RawPath = rp
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})
}
