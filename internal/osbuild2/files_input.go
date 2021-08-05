package osbuild2

import (
	"encoding/json"
	"fmt"
)

const InputTypeFiles string = "org.osbuild.files"

type SourceFilesInputReferences interface {
	isSourceFilesReferences()
	isReferences()
}

type PlainFilesReferences []string

func (PlainFilesReferences) isSourceFilesReferences() {}
func (PlainFilesReferences) isReferences()            {}

type SourceObjectReferences map[string]FilesSourceOptions

func (SourceObjectReferences) isSourceFilesReferences() {}
func (SourceObjectReferences) isReferences()            {}

type FilesSourceOptions struct {
	// Additional metadata to forward to the stage
	Metadata map[string]interface{} `json:"metadata"`
}

func NewSourceFilesInput(references SourceFilesInputReferences) *Input {
	return &Input{
		Type:       InputTypeFiles,
		Origin:     InputOriginSource,
		References: references,
	}
}

type PipelineFilesInputReferences map[string]FilesPipelineOptions

func (PipelineFilesInputReferences) isReferences() {}

type FilesPipelineOptions struct {
	// File to access within the pipeline
	File string `json:"file"`

	// Additional metadata to forward to the stage
	Metadata map[string]interface{} `json:"metadata"`
}

func NewPipelineFilesInput(references FilesInputReferencesPipeline) *Input {
	return &Input{
		Type:       InputTypeFiles,
		Origin:     InputOriginPipeline,
		References: references,
	}
}

// func NewFilesInput(references FilesInputReferences) *Input {

// 	switch t := references.(type) {
// 	case *FilesInputReferencesPipeline:
// 		input.Origin = InputOriginPipeline
// 	default:
// 		panic(fmt.Sprintf("unknown FilesInputReferences type: %v", t))
// 	}

// 	input.References = references

// 	return *input
// }

type rawFilesInput struct {
	Type       string
	Origin     string
	References json.RawMessage `json:"references"`
}

// func (f *FilesInput) UnmarshalJSON(data []byte) error {
// var rawFilesInput rawFilesInput
// if err := json.Unmarshal(data, &rawFilesInput); err != nil {
// 	return err
// }

// var ref FilesInputReferences
// switch rawFilesInput.Origin {
// case InputOriginPipeline:
// 	ref = &FilesInputReferencesPipeline{}
// default:
// 	return fmt.Errorf("FilesInput: unknown input origin: %s", rawFilesInput.Origin)
// }

// if err := json.Unmarshal(rawFilesInput.References, ref); err != nil {
// 	return err
// }

// f.Type = rawFilesInput.Type
// f.Origin = rawFilesInput.Origin
// f.References = ref

// return nil
// }

// SUPPORTED FILE INPUT REFERENCES

type FilesInputReferences interface {
	isFilesInputReferences()
	isInputReference()
}

// The expected JSON structure is:
// `"name:<pipeline_name>": {"file": "<filename>"}`
type FilesInputReferencesPipeline map[string]FileReference

func (*FilesInputReferencesPipeline) isFilesInputReferences() {}
func (*FilesInputReferencesPipeline) isInputReference()       {}

type FileReference struct {
	File string `json:"file"`
}

func NewFilesInputReferencesPipeline(pipeline, filename string) FilesInputReferences {
	ref := &FilesInputReferencesPipeline{
		fmt.Sprintf("name:%s", pipeline): {File: filename},
	}
	return ref
}

// TODO: define FilesInputReferences for "sources"
