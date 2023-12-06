package server

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
	"strconv"
)

func etag(filename string) (string, error) {
	h := md5.New()
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return strconv.Quote(base64.StdEncoding.EncodeToString(h.Sum(nil))), nil
}
