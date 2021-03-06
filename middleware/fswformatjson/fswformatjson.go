package fswformatjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/webnice/lv2/middleware"

	s "github.com/webnice/lv2/sender"
)

// New Create new package implementation and return interface
func New(filename ...string) middleware.FsWriter {
	var (
		fsw *impl
		fnm string
		tmp []string
	)

	fsw = new(impl)
	tmp = strings.Split(os.Args[0], string(os.PathSeparator))
	if len(tmp) > 0 {
		fsw.SetFilename(tmp[len(tmp)-1] + `.log`)
	} else {
		fsw.SetFilename(os.Args[0] + `.log`)
	}
	for _, fnm = range filename {
		fsw.SetFilename(fnm)
	}

	return fsw
}

// SetFilename Set filename
func (fsw *impl) SetFilename(filename string) middleware.FsWriter {
	fsw.Lock()
	defer fsw.Unlock()
	fsw.Filename = filename

	return fsw
}

// SetFilemode Set filemode
func (fsw *impl) SetFilemode(filemode os.FileMode) middleware.FsWriter {
	fsw.Filemode = filemode
	return fsw
}

// SetFormat Set template line formatting
func (fsw *impl) SetFormat(f string) middleware.FsWriter { return fsw }

// WriteMessage Запись среза байт в файл
func (fsw *impl) WriteMessage(msg s.Message) (n int, err error) {
	var buf *bytes.Buffer

	buf = bytes.NewBufferString(``)
	if err = json.NewEncoder(buf).Encode(msg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "encode json error: %s", err)
		return
	}
	n, err = fsw.Write(buf.Bytes())

	return
}

// Write Запись среза байт как есть
func (fsw *impl) Write(buf []byte) (n int, err error) {
	var out *os.File

	fsw.Lock()
	defer fsw.Unlock()
	if out, err = os.OpenFile(fsw.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0644)); err != nil {
		err = fmt.Errorf("open file %q error: %s", fsw.Filename, err)
		return
	}
	defer func() { _ = out.Close() }()
	n, err = out.Write(buf)

	return
}
