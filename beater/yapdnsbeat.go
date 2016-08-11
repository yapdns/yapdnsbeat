package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"

	cfg "github.com/yapdns/yapdnsbeat/config"
	"github.com/yapdns/yapdnsbeat/crawler"
	"github.com/yapdns/yapdnsbeat/input"
)

// YapdnsBeat is a beater object. Contains all objects needed to run the beat
type YapdnsBeat struct {
	YbConfig *cfg.Config
	// Channel from harvesters to spooler
	publisherChan chan []*input.FileEvent
	spooler       *Spooler
	registrar     *crawler.Registrar
	crawler       *crawler.Crawler
	pub           logPublisher
	done          chan struct{}
}

// New creates a new YapdnsBeat pointer instance.
func New() *YapdnsBeat {
	return &YapdnsBeat{}
}

// Config setups up the filebeat configuration by fetch all additional config files
func (yb *YapdnsBeat) Config(b *beat.Beat) error {

	// Load Base config
	err := b.RawConfig.Unpack(&yb.YbConfig)

	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	// Check if optional config_dir is set to fetch additional prospector config files
	yb.YbConfig.FetchConfigs()

	return nil
}

// Setup applies the minimum required setup to a new YapdnsBeat instance for use.
func (yb *YapdnsBeat) Setup(b *beat.Beat) error {
	yb.done = make(chan struct{})

	// ClientId = yb.YbConfig.ClientId
	// ClientSecret = yb.YbConfig.ClientSecret

	// jsonEvent :=
	// // send POST request to validate client
	// resp, err = http.Post(, "application/json", bytes.NewBuffer(jsonEvent))

	// if err != nil {
	// 	logp.Err("Failed to send POST request to %v", h.config.ApiEndpoint)
	// }
	return nil
}

// Run allows the beater to be run as a beat.
func (yb *YapdnsBeat) Run(b *beat.Beat) error {

	var err error

	// Init channels
	yb.publisherChan = make(chan []*input.FileEvent, 1)

	// Setup registrar to persist state
	yb.registrar, err = crawler.NewRegistrar(yb.YbConfig.YapdnsBeat.RegistryFile)
	if err != nil {
		logp.Err("Could not init registrar: %v", err)
		return err
	}

	yb.crawler = &crawler.Crawler{
		Registrar: yb.registrar,
	}

	// Load the previous log file locations now, for use in prospector
	yb.registrar.LoadState()

	// Init and Start spooler: Harvesters dump events into the spooler.
	yb.spooler = NewSpooler(yb.YbConfig.YapdnsBeat, yb.publisherChan)

	if err != nil {
		logp.Err("Could not init spooler: %v", err)
		return err
	}

	yb.registrar.Start()
	yb.spooler.Start()

	err = yb.crawler.Start(yb.YbConfig.YapdnsBeat.Prospectors, yb.spooler.Channel)
	if err != nil {
		return err
	}

	// Publishes event to output
	yb.pub = newPublisher(yb.YbConfig.YapdnsBeat.PublishAsync,
		yb.publisherChan, yb.registrar.Channel, b.Publisher.Connect())
	yb.pub.Start()

	// Blocks progressing
	<-yb.done

	return nil
}

// Cleanup removes any temporary files, data, or other items that were created by the Beat.
func (yb *YapdnsBeat) Cleanup(b *beat.Beat) error {
	return nil
}

// Stop is called on exit to stop the crawling, spooling and registration processes.
func (yb *YapdnsBeat) Stop() {

	logp.Info("Stopping filebeat")

	// Stop crawler -> stop prospectors -> stop harvesters
	yb.crawler.Stop()

	// Stopping spooler will flush items
	yb.spooler.Stop()

	// stopping publisher (might potentially drop items)
	yb.pub.Stop()

	// Stopping registrar will write last state
	yb.registrar.Stop()

	// Stop YapdnsBeat
	close(yb.done)
}
