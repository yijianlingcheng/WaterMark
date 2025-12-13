package pkg

import (
	"bytes"
	"compress/zlib"
	"io"
)

// 进行zlib压缩字符串.
func ZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)

	_, err := w.Write(src)
	if err != nil {
		return []byte("")
	}

	err = w.Close()
	if err != nil {
		return []byte("")
	}

	return in.Bytes()
}

// 进行zlib解压缩.
func ZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	// 打开失败返回空
	if err != nil {
		return []byte("")
	}

	// 复制失败返回空
	//nolint:gosec
	_, err = io.Copy(&out, r)
	if err != nil {
		return []byte("")
	}

	return out.Bytes()
}
