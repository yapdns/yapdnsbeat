// +build !integration

package beater

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
	cfg "github.com/yapdns/yapdnsbeat/config"
)

func load(t *testing.T, in string) cfg.YapdnsBeatConfig {
	yaml, err := common.NewConfigWithYAML([]byte(in), "")
	if err != nil {
		t.Fatalf("Failed to parse config input: %v", err)
	}

	var config cfg.YapdnsBeatConfig
	err = yaml.Unpack(&config)
	if err != nil {
		t.Fatalf("Failed to unpack config: %v", err)
	}

	return config
}

func TestNewSpoolerDefaultConfig(t *testing.T) {
	config := load(t, "")

	// Read from empty yaml config
	spooler := NewSpooler(config, nil)

	assert.Equal(t, cfg.DefaultSpoolSize, spooler.spoolSize)
	assert.Equal(t, cfg.DefaultIdleTimeout, spooler.idleTimeout)
}

func TestNewSpoolerSpoolSize(t *testing.T) {
	spoolSize := uint64(19)
	config := cfg.YapdnsBeatConfig{SpoolSize: spoolSize}
	spooler := NewSpooler(config, nil)

	assert.Equal(t, spoolSize, spooler.spoolSize)
}

func TestNewSpoolerIdleTimeout(t *testing.T) {
	config := load(t, "idle_timeout: 10s")
	spooler := NewSpooler(config, nil)

	assert.Equal(t, time.Duration(10*time.Second), spooler.idleTimeout)
}
