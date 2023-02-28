package gt

import (
	"errors"
)

const (
	JSON DecodeFormat = "json"
	YAML DecodeFormat = "yaml"
	BODY DecodeFormat = "body"
)

var (
	DecoderTypeNotSupport = errors.New("decoder type not support")
)

type DecodeFormat string
