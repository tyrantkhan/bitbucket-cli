package models

import (
	"fmt"
	"time"
)

// Pipeline represents a Bitbucket pipeline.
type Pipeline struct {
	UUID        string         `json:"uuid"`
	BuildNumber int            `json:"build_number"`
	State       *PipelineState `json:"state"`
	Target      PipelineTarget `json:"target"`
	CreatedOn   string         `json:"created_on"`
	CompletedOn string         `json:"completed_on"`
	Creator     User           `json:"creator"`
	Links       Links          `json:"links"`
}

// PipelineState represents the pipeline state.
type PipelineState struct {
	Name   string              `json:"name"` // PENDING, IN_PROGRESS, COMPLETED, etc.
	Result *PipelineStateValue `json:"result,omitempty"`
	Stage  *PipelineStateValue `json:"stage,omitempty"`
}

// PipelineStateValue holds a state result or stage name.
type PipelineStateValue struct {
	Name string `json:"name"` // SUCCESSFUL, FAILED, STOPPED, etc.
}

// PipelineTarget represents what triggered the pipeline.
type PipelineTarget struct {
	Type     string  `json:"type"` // pipeline_ref_target, pipeline_custom_target
	RefType  string  `json:"ref_type"`
	RefName  string  `json:"ref_name"`
	Selector *struct {
		Type    string `json:"type"`
		Pattern string `json:"pattern"`
	} `json:"selector,omitempty"`
	Commit Commit `json:"commit"`
}

// PipelineStep represents a step in a pipeline.
type PipelineStep struct {
	UUID        string         `json:"uuid"`
	Name        string         `json:"name"`
	State       *PipelineState `json:"state"`
	StartedOn   string         `json:"started_on"`
	CompletedOn string         `json:"completed_on"`
	LogLink     string         `json:"-"` // Computed from links
	Links       Links          `json:"links"`
}

// StatusText returns a human-readable status string for a pipeline.
func (p *Pipeline) StatusText() string {
	if p.State == nil {
		return "UNKNOWN"
	}
	if p.State.Result != nil {
		return p.State.Result.Name
	}
	if p.State.Stage != nil {
		return p.State.Stage.Name
	}
	return p.State.Name
}

// Duration returns the pipeline duration as a string.
func (p *Pipeline) Duration() string {
	if p.CreatedOn == "" {
		return ""
	}
	start, err := time.Parse(time.RFC3339Nano, p.CreatedOn)
	if err != nil {
		return ""
	}

	var end time.Time
	if p.CompletedOn != "" {
		end, err = time.Parse(time.RFC3339Nano, p.CompletedOn)
		if err != nil {
			return ""
		}
	} else {
		end = time.Now()
	}

	d := end.Sub(start)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
