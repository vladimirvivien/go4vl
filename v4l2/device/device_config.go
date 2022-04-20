package device

import (
	"github.com/vladimirvivien/go4vl/v4l2"
)

type Config struct {
	ioType v4l2.IOType
	pixFormat v4l2.PixFormat
	bufSize uint32
	fps uint32
}

type Option func(*Config)

func WithIOType(ioType v4l2.IOType) Option {
	return func(o *Config) {
		o.ioType = ioType
	}
}

func WithPixFormat(pixFmt v4l2.PixFormat) Option {
	return func(o *Config) {
		o.pixFormat = pixFmt
	}
}

func WithBufferSize(size uint32) Option {
	return func(o *Config) {
		o.bufSize = size
	}
}

func WithFPS(fps uint32) Option {
	return func(o *Config) {
		o.fps = fps
	}
}