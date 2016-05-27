package http

import (
	"encoding/json"
	"os"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
)

func init() {
	outputs.RegisterOutputPlugin("http", New)
}

type http struct {
	config config
}

func New(config *common.Config, _ int) (outputs.Outputer, error) {
	c := &http{config: defaultConfig}
	err := config.Unpack(&c.config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newHttp(pretty bool) *http {
	return &http{config{pretty, true}}
}

func writeBuffer(buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := os.Stdout.Write(buf[written:])
		if err != nil {
			return err
		}

		written += n
	}
	return nil
}

// Implement Outputer
func (c *http) Close() error {
	return nil
}

func (c *http) PublishEvent(
	s op.Signaler,
	opts outputs.Options,
	event common.MapStr,
) error {
	var jsonEvent []byte
	var err error

	if c.config.Pretty {
		jsonEvent, err = json.MarshalIndent(event, "", "  ")
	} else {
		jsonEvent, err = json.Marshal(event)
	}
	if err != nil {
		logp.Err("Fail to convert the event to JSON (%v): %#v", err, event)
		op.SigCompleted(s)
		return err
	}

	if c.config.Foo {
		a := []byte("foo")
		jsonEvent = append(a, jsonEvent...)
	}

	if err = writeBuffer(jsonEvent); err != nil {
		goto fail
	}
	if err = writeBuffer([]byte{'\n'}); err != nil {
		goto fail
	}

	op.SigCompleted(s)
	return nil
fail:
	if opts.Guaranteed {
		logp.Critical("Unable to publish events to http: %v", err)
	}
	op.SigFailed(s, err)
	return err
}
