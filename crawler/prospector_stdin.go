package crawler

import (
	"fmt"

	"github.com/yapdns/yapdnsbeat/harvester"
	"github.com/yapdns/yapdnsbeat/input"
)

type ProspectorStdin struct {
	Prospector *Prospector
	harvester  *harvester.Harvester
	started    bool
}

func NewProspectorStdin(p *Prospector) (*ProspectorStdin, error) {

	prospectorer := &ProspectorStdin{
		Prospector: p,
	}

	var err error

	prospectorer.harvester, err = p.createHarvester(input.FileState{Source: "-"})
	if err != nil {
		return nil, fmt.Errorf("Error initializing stdin harvester: %v", err)
	}

	return prospectorer, nil
}

func (p *ProspectorStdin) Init() {
	p.started = false
}

func (p *ProspectorStdin) Run() {

	// Make sure stdin harvester is only started once
	if !p.started {
		go p.harvester.Harvest()
		p.started = true
	}
}
