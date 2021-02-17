package router

import (
	"io"
	"math/rand"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
	"time"
)

type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type FileStore struct {
	File        multipart.File
	FileHeader  *multipart.FileHeader
	MIMEHeader  textproto.MIMEHeader
	StoragePath string
}

func NewFileStore(file multipart.File, header *multipart.FileHeader) FileStore {
	return FileStore{
		File:       file,
		FileHeader: header,
		MIMEHeader: header.Header,
	}
}

func (f FileStore) SaveTo(path string) {
	f.SaveAs(path, f.HashName())
}

func (f FileStore) SaveAs(path, filename string) {
	if e := os.MkdirAll(path, os.ModePerm); e != nil {
		panic(e)
	}

	if "/" != string(path[len(path)-1]) {
		path += "/"
	}

	path += filename
	out, e := os.Create(path)
	if e != nil {
		panic(e)
	}

	defer out.Close()
	_, e = io.Copy(out, f.File)
	if e != nil {
		panic(e)
	}
}

func (f FileStore) GuessExtension() string {
	MIMETypes := NewMIMETypes()
	return MIMETypes.GuessExtension(f.MIMEType())
}

func (f FileStore) Extension() string {
	ext := f.GuessExtension()
	if ext == "undefined" {
		fn := strings.Split(f.FileName(), ".")
		ext = fn[len(fn)-1]
	}
	return ext
}

func (f FileStore) FileName() string {
	return f.FileHeader.Filename
}

func (f FileStore) Size() int64 {
	return f.FileHeader.Size
}

func (f FileStore) MIMEType() string {
	return f.MIMEHeader["Content-Type"][0]
}

func (f FileStore) Stream() File {
	return f.File
}

func (f FileStore) HashName() string {
	const letterBytes = "qwertyuiopasdfghjklzxcvbnm-QWERTYUIOPASDFGHJKLZXCVBNM_1234567890"
	n := 40
	bts := make([]byte, n)
	for i := range bts {
		rand.Seed(time.Now().UnixNano() + int64(n-i))
		bts[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bts) + "." + f.GuessExtension()
}
