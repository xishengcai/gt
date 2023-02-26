package gt

import "errors"

const (
	JSON DecodeFormat = "json"
	YAML DecodeFormat = "yaml"
)

var (
	DecoderTypeNotSupport = errors.New("decoder type not support")
)

type DecodeFormat string
