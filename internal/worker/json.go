package worker

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/osbuild/osbuild-composer/internal/distro"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
	"github.com/osbuild/osbuild-composer/internal/osbuild1"
	"github.com/osbuild/osbuild-composer/internal/osbuild2"
	"github.com/osbuild/osbuild-composer/internal/target"
)

//
// JSON-serializable types for the jobqueue
//

type OSBuildJob struct {
	Manifest        distro.Manifest  `json:"manifest"`
	Targets         []*target.Target `json:"targets,omitempty"`
	ImageName       string           `json:"image_name,omitempty"`
	StreamOptimized bool             `json:"stream_optimized,omitempty"`
}

type OSBuildJobResult struct {
	Success       bool           `json:"success"`
	OSBuildOutput osbuild.Result `json:"osbuild_output,omitempty"`
	TargetErrors  []string       `json:"target_errors,omitempty"`
	UploadStatus  string         `json:"upload_status"`
}

// Custom unmarshaller for the temporary dual implementations of osbuild.Result
func (jr *OSBuildJobResult) UnmarshalJSON(data []byte) error {
	var rawResult struct {
		Success       bool            `json:"success"`
		OSBuildOutput json.RawMessage `json:"osbuild_output,omitempty"`
		TargetErrors  []string        `json:"target_errors,omitempty"`
		UploadStatus  string          `json:"upload_status"`
	}

	if err := json.Unmarshal(data, &rawResult); err != nil {
		return err
	}

	// Add the rest of the fields to return even if the result unmarshalling fails
	jr.Success = rawResult.Success
	jr.TargetErrors = rawResult.TargetErrors
	jr.UploadStatus = rawResult.UploadStatus

	var br1 osbuild1.Result
	// TODO: Use decoder with DisallowUnknownFields
	if err := json.Unmarshal(rawResult.OSBuildOutput, &br1); err == nil && len(br1.Stages) != 0 {
		jr.OSBuildOutput = &br1
		return nil
	}

	var br2 osbuild2.Result
	if err := json.Unmarshal(rawResult.OSBuildOutput, &br2); err == nil && len(br2.Pipelines) != 0 {
		jr.OSBuildOutput = &br2
		return nil
	}
	return nil
}

type KojiInitJob struct {
	Server  string `json:"server"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Release string `json:"release"`
}

type KojiInitJobResult struct {
	BuildID   uint64 `json:"build_id"`
	Token     string `json:"token"`
	KojiError string `json:"koji_error"`
}

type OSBuildKojiJob struct {
	Manifest      distro.Manifest `json:"manifest"`
	ImageName     string          `json:"image_name"`
	KojiServer    string          `json:"koji_server"`
	KojiDirectory string          `json:"koji_directory"`
	KojiFilename  string          `json:"koji_filename"`
}

type OSBuildKojiJobResult struct {
	HostOS        string         `json:"host_os"`
	Arch          string         `json:"arch"`
	OSBuildOutput osbuild.Result `json:"osbuild_output"`
	ImageHash     string         `json:"image_hash"`
	ImageSize     uint64         `json:"image_size"`
	KojiError     string         `json:"koji_error"`
}

// Custom unmarshaller for the temporary dual implementations of osbuild.Result
func (jr *OSBuildKojiJobResult) UnmarshalJSON(data []byte) error {
	var rawResult struct {
		HostOS        string          `json:"host_os"`
		Arch          string          `json:"arch"`
		OSBuildOutput json.RawMessage `json:"osbuild_output"`
		ImageHash     string          `json:"image_hash"`
		ImageSize     uint64          `json:"image_size"`
		KojiError     string          `json:"koji_error"`
	}

	if err := json.Unmarshal(data, &rawResult); err != nil {
		return err
	}

	// Add the rest of the fields to return even if the result unmarshalling fails
	jr.HostOS = rawResult.HostOS
	jr.Arch = rawResult.Arch
	jr.ImageHash = rawResult.ImageHash
	jr.ImageSize = rawResult.ImageSize
	jr.KojiError = rawResult.KojiError

	var br1 osbuild1.Result
	// TODO: Use decoder with DisallowUnknownFields
	if err := json.Unmarshal(rawResult.OSBuildOutput, &br1); err == nil && len(br1.Stages) != 0 {
		jr.OSBuildOutput = &br1
		return nil
	}

	var br2 osbuild2.Result
	if err := json.Unmarshal(rawResult.OSBuildOutput, &br2); err == nil && len(br2.Pipelines) != 0 {
		jr.OSBuildOutput = &br2
		return nil
	}
	return nil
}

type KojiFinalizeJob struct {
	Server        string   `json:"server"`
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	Release       string   `json:"release"`
	KojiFilenames []string `json:"koji_filenames"`
	KojiDirectory string   `json:"koji_directory"`
	TaskID        uint64   `json:"task_id"` /* https://pagure.io/koji/issue/215 */
	StartTime     uint64   `json:"start_time"`
}

type KojiFinalizeJobResult struct {
	KojiError string `json:"koji_error"`
}

//
// JSON-serializable types for the HTTP API
//

type statusResponse struct {
	Status string `json:"status"`
}

type requestJobResponse struct {
	Id               uuid.UUID         `json:"id"`
	Location         string            `json:"location"`
	ArtifactLocation string            `json:"artifact_location"`
	Type             string            `json:"type"`
	Args             json.RawMessage   `json:"args,omitempty"`
	DynamicArgs      []json.RawMessage `json:"dynamic_args,omitempty"`
}

type getJobResponse struct {
	Canceled bool `json:"canceled"`
}

type updateJobRequest struct {
	Result json.RawMessage `json:"result"`
}

type updateJobResponse struct {
}
