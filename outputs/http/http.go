package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/op"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
)

func init() {
	outputs.RegisterOutputPlugin("http", New)
}

type httpApi struct {
	config config
}

func New(config *common.Config, _ int) (outputs.Outputer, error) {
	c := &httpApi{config: defaultConfig}
	err := config.Unpack(&c.config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newHttpApi() *httpApi {
	return &httpApi{config{}}
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
func (h *httpApi) Close() error {
	return nil
}

func (h *httpApi) PublishEvent(
	s op.Signaler,
	opts outputs.Options,
	event common.MapStr,
) error {
	var jsonEvent []byte
	var resp *http.Response
	var err error

	jsonEvent, err = json.Marshal(event)

	if err != nil {
		logp.Err("Fail to convert the event to JSON (%v): %#v", err, event)
		op.SigCompleted(s)
		return err
	}

	req, err := http.NewRequest("POST", h.config.ApiEndpoint, bytes.NewBuffer(jsonEvent))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(h.config.ClientId, h.config.ClientSecretKey)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		logp.Err("Failed to send POST request to %v", h.config.ApiEndpoint)
		goto fail
	}

	defer resp.Body.Close()

	op.SigCompleted(s)
	return nil
fail:
	if opts.Guaranteed {
		logp.Critical("Unable to publish events to http: %v", err)
	}
	op.SigFailed(s, err)
	return err
}

func (h *httpApi) BulkPublish(
	s op.Signaler,
	opts outputs.Options,
	event []common.MapStr,
) error {
	var jsonEvent []byte
	var resp *http.Response
	var err error

	jsonEvent, err = json.Marshal(event)

	if err != nil {
		logp.Err("Fail to convert the event to JSON (%v): %#v", err, event)
		op.SigCompleted(s)
		return err
	}

	req, err := http.NewRequest("POST", h.config.BulkApiEndpoint, bytes.NewBuffer(jsonEvent))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(h.config.ClientId, h.config.ClientSecretKey)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		logp.Err("Failed to send POST request to %v", h.config.ApiEndpoint)
		goto fail
	}

	defer resp.Body.Close()

	op.SigCompleted(s)
	return nil
fail:
	if opts.Guaranteed {
		logp.Critical("Unable to publish events to http: %v", err)
	}
	op.SigFailed(s, err)
	return err
}
