package model

import "io"

type S3Object struct {
	Body          io.ReadCloser
	ContentLength int64
	ContentType   string
}
