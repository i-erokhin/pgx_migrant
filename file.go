package migrations

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"text/template"
)

var fileRe = regexp.MustCompile(`^(\d+)(-\w+)?\.sql$`)

type File struct {
	Path     string
	Basename string
	Number   int
}

// path берется из файловой системы на этапе создания файл не читается с диска
func NewFile(pa string) (file *File, err error) {
	base := path.Base(pa)
	parts := fileRe.FindStringSubmatch(base)
	if parts == nil {
		err = fmt.Errorf("bad path: %s", pa)
		return
	}
	number, err := strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("bad path (migration number not found at begin): %s", pa)
		return
	}
	if number == 0 {
		err = fmt.Errorf("number catn`t be null: %s", pa)
		return
	}
	file = &File{
		Path:     pa,
		Basename: base,
		Number:   number,
	}
	return
}

func (f *File) GetSQL(data interface{}) (sql string, err error) {
	t, err := template.ParseFiles(f.Path)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return
	}
	sql = buf.String()
	return
}

