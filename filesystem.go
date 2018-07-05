package migrations

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"sort"
)

type Filesystem struct {
	Path      string
	Files     []*File
	MaxNumber int
}

func NewFilesystem(pa string) (fs *Filesystem, err error) {
	fs = &Filesystem{
		Path: pa,
	}

	files, err := ioutil.ReadDir(pa)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		var file *File
		file, err = NewFile(path.Join(pa, f.Name()))
		if err != nil {
			return
		}
		fs.Files = append(fs.Files, file)
		if fs.MaxNumber < file.Number {
			fs.MaxNumber = file.Number
		}
	}
	sort.Slice(fs.Files, func(i, j int) bool {
		return fs.Files[i].Number < fs.Files[j].Number
	})
	return
}

func (f *Filesystem) PrettyJSON() string {
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
