package caddy

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap/zapcore"
)

type Logging struct {
	Sink       *StandardLibLog
	Logs       map[string]*CustomLog
	writerKeys []string
}

type StandardLibLog struct {
	WriterRaw json.RawMessage `json:"writer,omitempty" caddy:"namespace=caddy.logging.writers inline_key=output"`
	write     io.WriteCloser
}

type CustomLog struct {
	WriteRaw   json.RawMessage `json:"writer,omitempty" caddy:"namespace=caddy.logging.writers inline_key=output"`
	EncoderRaw json.RawMessage `json:"encoder,omitempty" caddy:"namespace=caddy.logging.encoders inline_key=format"`
	Level      string          `json:"level,omitempty"`
	Sampling   *LogSampling    `json:"sampling,omitempty"`
	Include    []string        `json:"include,omitempty"`
	Exclude    []string        `json:"encoding,omitempty"`

	writerOpener WriterOpener
	writer       io.WriteCloser
	encoder      zapcore.Encoder
	levelEnable  zapcore.LevelEnabler
	core         zapcore.Core
}

type LogSampling struct {
	Interval   time.Duration
	First      int
	Thereafter int
}

type WriterOpener interface {
	fmt.Stringer
}