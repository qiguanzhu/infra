/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/8/29 -- 11:17
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: file.go
*/

package xhttp

import (
	"io"
	"mime/multipart"
	"os"
)

type FormFiles map[string]FormFile

func (f FormFiles) AddFilePath(name, filename, filepath string) {
	if filename == "" {
		filename = name
	}
	f.Add(name, file{name: name, filename: filename, filepath: filepath})
}

func (f FormFiles) AddFileObject(name, filename string, file io.Reader) {
	if filename == "" {
		filename = name
	}
	f.Add(name, fileObject{name: name, filename: filename, reader: file})
}

func (f FormFiles) Add(name string, file FormFile) {
	f[name] = file
}

func (f FormFiles) Del(name string) {
	delete(f, name)
}

func (f FormFiles) Has(name string) bool {
	_, ok := f[name]
	return ok
}

type FormFile interface {
	WriteTo(writer *multipart.Writer) error
}

type file struct {
	name     string
	filename string
	filepath string
}

func (f file) WriteTo(writer *multipart.Writer) error {
	nFile, err := os.Open(f.filepath)
	if err != nil {
		return err
	}
	defer nFile.Close()
	nWriter, err := writer.CreateFormFile(f.name, f.filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(nWriter, nFile); err != nil {
		return err
	}
	return nil
}

type fileObject struct {
	name     string
	filename string
	reader   io.Reader
}

func (f fileObject) WriteTo(writer *multipart.Writer) error {
	if f.reader == nil {
		return nil
	}
	nWriter, err := writer.CreateFormFile(f.name, f.filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(nWriter, f.reader); err != nil {
		return err
	}
	return nil
}
