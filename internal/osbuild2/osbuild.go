// Package osbuild provides primitives for representing and (un)marshalling
// OSBuild (schema v1) types.
package osbuild2

// A Manifest represents an OSBuild source and pipeline manifest
type Manifest struct {
	Pipelines []Pipeline `json:"pipeline"`
	Sources   Sources    `json:"sources"`
}

// A Pipeline represents an OSBuild pipeline
type Pipeline struct {
	Name string
	// The build environment which can run this pipeline
	Build string

	Runner string

	// Sequence of stages that produce the filesystem tree, which is the
	// payload of the produced image.
	Stages []*Stage `json:"stages,omitempty"`
}

// SetBuild sets the pipeline and runner for generating the build environment
// for a pipeline.
func (p *Pipeline) SetBuild(build string) {
	p.Build = build
}

// AddStage appends a stage to the list of stages of a pipeline. The stages
// will be executed in the order they are appended.
func (p *Pipeline) AddStage(stage *Stage) {
	p.Stages = append(p.Stages, stage)
}
