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
 @Time    : 2024/11/5 -- 16:39
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: serializer.go
*/

package store

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

type Serializer interface {
	decode([]byte, interface{}) error
	encode(value interface{}) ([]byte, error)
}

type JsonSerializer struct {
	compress   bool
	escapeHTML bool
}

func NewJsonSerializer(compress, escapeHTML bool) *JsonSerializer {
	return &JsonSerializer{
		compress:   compress,
		escapeHTML: escapeHTML,
	}
}

func (j *JsonSerializer) decode(src []byte, dst interface{}) error {
	if j.compress {
		decompressor, err := gzip.NewReader(bytes.NewReader(src))
		if err != nil {
			return err
		}
		defer decompressor.Close()

		return json.NewDecoder(decompressor).Decode(dst)
	}

	return json.Unmarshal(src, dst)
}
func (j *JsonSerializer) encode(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	if j.compress {
		compressor := gzip.NewWriter(&buf)
		defer compressor.Close()

		encoder := json.NewEncoder(compressor)
		encoder.SetEscapeHTML(j.escapeHTML)
		if err := encoder.Encode(value); err != nil {
			return nil, err
		}

		if err := compressor.Flush(); err != nil {
			return nil, err
		}
	} else {
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(j.escapeHTML)
		if err := encoder.Encode(value); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
