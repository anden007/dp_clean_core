package misc

import (
	"io/fs"
	"net/http"
	"os"
)

func GetFileSystem(useOS bool, dir string, fsys fs.FS) http.FileSystem {
	if useOS {
		return http.FS(os.DirFS(dir))
	}
	fsys, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
