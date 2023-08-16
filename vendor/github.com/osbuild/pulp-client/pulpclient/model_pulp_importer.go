/*
Pulp 3 API

Fetch, Upload, Organize, and Distribute Software Packages

API version: v3
Contact: pulp-list@redhat.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package pulpclient

import (
	"encoding/json"
)

// checks if the PulpImporter type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PulpImporter{}

// PulpImporter Serializer for PulpImporters.
type PulpImporter struct {
	// Unique name of the Importer.
	Name string `json:"name"`
	// Mapping of repo names in an export file to the repo names in Pulp. For example, if the export has a repo named 'foo' and the repo to import content into was 'bar', the mapping would be \"{'foo': 'bar'}\".
	RepoMapping *map[string]string `json:"repo_mapping,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PulpImporter PulpImporter

// NewPulpImporter instantiates a new PulpImporter object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPulpImporter(name string) *PulpImporter {
	this := PulpImporter{}
	this.Name = name
	return &this
}

// NewPulpImporterWithDefaults instantiates a new PulpImporter object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPulpImporterWithDefaults() *PulpImporter {
	this := PulpImporter{}
	return &this
}

// GetName returns the Name field value
func (o *PulpImporter) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *PulpImporter) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *PulpImporter) SetName(v string) {
	o.Name = v
}

// GetRepoMapping returns the RepoMapping field value if set, zero value otherwise.
func (o *PulpImporter) GetRepoMapping() map[string]string {
	if o == nil || IsNil(o.RepoMapping) {
		var ret map[string]string
		return ret
	}
	return *o.RepoMapping
}

// GetRepoMappingOk returns a tuple with the RepoMapping field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PulpImporter) GetRepoMappingOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.RepoMapping) {
		return nil, false
	}
	return o.RepoMapping, true
}

// HasRepoMapping returns a boolean if a field has been set.
func (o *PulpImporter) HasRepoMapping() bool {
	if o != nil && !IsNil(o.RepoMapping) {
		return true
	}

	return false
}

// SetRepoMapping gets a reference to the given map[string]string and assigns it to the RepoMapping field.
func (o *PulpImporter) SetRepoMapping(v map[string]string) {
	o.RepoMapping = &v
}

func (o PulpImporter) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PulpImporter) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	if !IsNil(o.RepoMapping) {
		toSerialize["repo_mapping"] = o.RepoMapping
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *PulpImporter) UnmarshalJSON(bytes []byte) (err error) {
	varPulpImporter := _PulpImporter{}

	if err = json.Unmarshal(bytes, &varPulpImporter); err == nil {
		*o = PulpImporter(varPulpImporter)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "repo_mapping")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePulpImporter struct {
	value *PulpImporter
	isSet bool
}

func (v NullablePulpImporter) Get() *PulpImporter {
	return v.value
}

func (v *NullablePulpImporter) Set(val *PulpImporter) {
	v.value = val
	v.isSet = true
}

func (v NullablePulpImporter) IsSet() bool {
	return v.isSet
}

func (v *NullablePulpImporter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePulpImporter(val *PulpImporter) *NullablePulpImporter {
	return &NullablePulpImporter{value: val, isSet: true}
}

func (v NullablePulpImporter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePulpImporter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


