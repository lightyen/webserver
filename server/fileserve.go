package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	WebRoot    string
	EnableGzip bool
	gzPool     = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
)

func getAcceptEncoding(h http.Header) string {
	return h.Get("Accept-Encoding")
}

func fileServe() gin.HandlerFunc {

	fs := static.LocalFile(WebRoot, false)
	serve := http.StripPrefix("/", http.FileServer(fs))

	return func(c *gin.Context) {
		if c.Request.URL.Path != "/" && fs.Exists("/", c.Request.URL.Path) {
			filename := filepath.Join(WebRoot, c.Request.URL.Path)
			if eTag, _ := etag(filename); eTag != "" {
				c.Header("Cache-Control", "max-age=30")
				c.Header("Etag", eTag)
			}

			if EnableGzip {
				if e := ParseAcceptEncoding(getAcceptEncoding(c.Request.Header)); e.Contains("gzip") {
					c.Header("Content-Encoding", "gzip")
					c.Header("Vary", "Accept-Encoding")

					gz := gzPool.Get().(*gzip.Writer)
					defer gzPool.Put(gz)

					gz.Reset(c.Writer)
					defer gz.Reset(io.Discard)

					c.Writer = &gzWriter{c.Writer, gz}
					defer gz.Close()
				}
			}

			serve.ServeHTTP(c.Writer, c.Request)
			return
		}

		fallback(true)(c)
	}
}

func staticPath(req *http.Request) bool {
	ext := filepath.Ext(req.URL.Path)
	if ext == "" {
		return false
	}
	return true
}

func fallback(useAny bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			return
		}

		if staticPath(c.Request) {
			return
		}

		if a := ParseAccept(c.Request.Header.Get("Accept")); a.Contains("text/html") || (useAny && a.Contains("*/*")) {
			filename := filepath.Join(WebRoot, "index.html")
			eTag, _ := etag(filename)
			im := c.Request.Header.Get("If-Match")
			if im != "" && im == eTag {
				c.Status(http.StatusNotModified)
				return
			}

			if eTag != "" {
				c.Header("Cache-Control", "no-cache, max-age=0, private, must-revalidate")
				c.Header("Etag", eTag)
			}

			if EnableGzip {
				if e := ParseAcceptEncoding(getAcceptEncoding(c.Request.Header)); e.Contains("gzip") {
					c.Header("Content-Encoding", "gzip")
					c.Header("Vary", "Accept-Encoding")

					gz := gzPool.Get().(*gzip.Writer)
					defer gzPool.Put(gz)

					gz.Reset(c.Writer)
					defer gz.Reset(io.Discard)

					c.Writer = &gzWriter{c.Writer, gz}
					defer gz.Close()
				}
			}

			c.File(filename)
		}
	}
}
