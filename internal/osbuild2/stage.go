package osbuild2

import (
	"encoding/json"
	"fmt"
)

// Single stage of a pipeline executing one step
type Stage struct {
	// Well-known name in reverse domain-name notation, uniquely identifying
	// the stage type.
	Type string `json:"type"`
	// Stage-type specific options fully determining the operations of the
	// stage.
	Inputs  Inputs       `json:"inputs"`
	Options StageOptions `json:"options"`
}

// Collection of Inputs for a Stage
type Inputs map[string]Input

// Single Input for a Stage
type Input interface {
	isInput()
}

// Fields shared between all Input types (should be embedded in each instance)
type inputCommon struct {
	Type string `json:"type"`
	// Origin should be either 'org.osbuild.source' or 'org.osbuild.pipeline'
	// TODO: Enum?
	Origin string `json:"origin"`

	References map[string]Reference `json:"references"`
}

type StageInput interface {
	isStageInput()
}

// NOTE: [Array of string] or [object]
type Reference interface{}

// StageOptions specify the operations of a given stage-type.
type StageOptions interface {
	isStageOptions()
}

type InputOptions interface {
}

type rawStage struct {
	Name    string          `json:"name"`
	Options json.RawMessage `json:"options"`
}

// UnmarshalJSON unmarshals JSON into a Stage object. Each type of stage has
// a custom unmarshaller for its options, selected based on the stage name.
func (stage *Stage) UnmarshalJSON(data []byte) error {
	var rawStage rawStage
	err := json.Unmarshal(data, &rawStage)
	if err != nil {
		return err
	}
	var options StageOptions
	switch rawStage.Name {
	case "org.osbuild.fix-bls":
		// TODO: verify that we can unmarshall this also if "options" is omitted
		options = new(FixBLSStageOptions)
	case "org.osbuild.fstab":
		options = new(FSTabStageOptions)
	case "org.osbuild.grub2":
		options = new(GRUB2StageOptions)
	case "org.osbuild.locale":
		options = new(LocaleStageOptions)
	case "org.osbuild.selinux":
		options = new(SELinuxStageOptions)
	case "org.osbuild.hostname":
		options = new(HostnameStageOptions)
	case "org.osbuild.users":
		options = new(UsersStageOptions)
	case "org.osbuild.groups":
		options = new(GroupsStageOptions)
	case "org.osbuild.timezone":
		options = new(TimezoneStageOptions)
	case "org.osbuild.chrony":
		options = new(ChronyStageOptions)
	case "org.osbuild.keymap":
		options = new(KeymapStageOptions)
	case "org.osbuild.firewall":
		options = new(FirewallStageOptions)
	case "org.osbuild.rhsm":
		options = new(RHSMStageOptions)
	case "org.osbuild.rpm":
		options = new(RPMStageOptions)
	case "org.osbuild.rpm-ostree":
		options = new(RPMOSTreeStageOptions)
	case "org.osbuild.systemd":
		options = new(SystemdStageOptions)
	case "org.osbuild.script":
		options = new(ScriptStageOptions)
	default:
		return fmt.Errorf("unexpected stage name: %s", rawStage.Name)
	}
	err = json.Unmarshal(rawStage.Options, options)
	if err != nil {
		return err
	}

	stage.Type = rawStage.Name
	stage.Options = options

	return nil
}
