package osbuild2

import (
	"encoding/json"
	"fmt"
	"io"
)

type PipelineResult struct {
	Name   string        `json:"name"`
	Build  string        `json:"string"`
	Runner string        `json:"runner"`
	Stages []StageResult `json:"stages"`
}

type StageResult struct {
	Type     string          `json:"name"`
	Options  json.RawMessage `json:"options"`
	Success  bool            `json:"success"`
	Output   string          `json:"output"`
	Metadata StageMetadata   `json:"metadata"`
}

// StageMetadata specify the metadata of a given stage-type.
type StageMetadata interface {
	isStageMetadata()
}

type rawStageResult struct {
	Type     string          `json:"name"`
	Options  json.RawMessage `json:"options"`
	Success  bool            `json:"success"`
	Output   string          `json:"output"`
	Metadata json.RawMessage `json:"metadata"`
}

type Result struct {
	TreeID    string           `json:"tree_id"`
	OutputID  string           `json:"output_id"`
	Pipelines []PipelineResult `json:"pipelines"`
	Success   bool             `json:"success"`
}

func (result *StageResult) UnmarshalJSON(data []byte) error {
	var rawStageResult rawStageResult
	err := json.Unmarshal(data, &rawStageResult)
	if err != nil {
		return err
	}
	var metadata StageMetadata
	switch rawStageResult.Type {
	case "org.osbuild.rpm":
		metadata = new(RPMStageMetadata)
		err = json.Unmarshal(rawStageResult.Metadata, metadata)
		if err != nil {
			return err
		}
	default:
		metadata = nil
	}

	result.Type = rawStageResult.Type
	result.Options = rawStageResult.Options
	result.Success = rawStageResult.Success
	result.Output = rawStageResult.Output
	result.Metadata = metadata

	return nil
}

func (cr *Result) Write(writer io.Writer) error {
	if len(cr.Pipelines) == 0 {
		fmt.Fprintf(writer, "The compose result is empty.\n")
		return nil
	}

	fmt.Fprintf(writer, "Pipelines:\n")
	for _, pipeline := range cr.Pipelines {
		fmt.Fprintf(writer, "Pipeline: %s\n", pipeline.Name)
		fmt.Fprintf(writer, "Stages:\n")
		for _, stage := range pipeline.Stages {
			fmt.Fprintf(writer, "Stage: %s\n", stage.Type)
			enc := json.NewEncoder(writer)
			enc.SetIndent("", "  ")
			err := enc.Encode(stage.Options)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, "\nOutput:\n%s\n", stage.Output)
		}
	}

	return nil
}

func (cr *Result) Succeeded() bool {
	return cr.Success
}
