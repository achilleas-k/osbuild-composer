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

// checks if the PatchedPulpExporter type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PatchedPulpExporter{}

// PatchedPulpExporter Serializer for pulp exporters.
type PatchedPulpExporter struct {
	// Unique name of the file system exporter.
	Name *string `json:"name,omitempty"`
	// File system directory to store exported tar.gzs.
	Path *string `json:"path,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
	// Last attempted export for this PulpExporter
	LastExport NullableString `json:"last_export,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PatchedPulpExporter PatchedPulpExporter

// NewPatchedPulpExporter instantiates a new PatchedPulpExporter object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchedPulpExporter() *PatchedPulpExporter {
	this := PatchedPulpExporter{}
	return &this
}

// NewPatchedPulpExporterWithDefaults instantiates a new PatchedPulpExporter object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchedPulpExporterWithDefaults() *PatchedPulpExporter {
	this := PatchedPulpExporter{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PatchedPulpExporter) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedPulpExporter) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PatchedPulpExporter) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PatchedPulpExporter) SetName(v string) {
	o.Name = &v
}

// GetPath returns the Path field value if set, zero value otherwise.
func (o *PatchedPulpExporter) GetPath() string {
	if o == nil || IsNil(o.Path) {
		var ret string
		return ret
	}
	return *o.Path
}

// GetPathOk returns a tuple with the Path field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedPulpExporter) GetPathOk() (*string, bool) {
	if o == nil || IsNil(o.Path) {
		return nil, false
	}
	return o.Path, true
}

// HasPath returns a boolean if a field has been set.
func (o *PatchedPulpExporter) HasPath() bool {
	if o != nil && !IsNil(o.Path) {
		return true
	}

	return false
}

// SetPath gets a reference to the given string and assigns it to the Path field.
func (o *PatchedPulpExporter) SetPath(v string) {
	o.Path = &v
}

// GetRepositories returns the Repositories field value if set, zero value otherwise.
func (o *PatchedPulpExporter) GetRepositories() []string {
	if o == nil || IsNil(o.Repositories) {
		var ret []string
		return ret
	}
	return o.Repositories
}

// GetRepositoriesOk returns a tuple with the Repositories field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedPulpExporter) GetRepositoriesOk() ([]string, bool) {
	if o == nil || IsNil(o.Repositories) {
		return nil, false
	}
	return o.Repositories, true
}

// HasRepositories returns a boolean if a field has been set.
func (o *PatchedPulpExporter) HasRepositories() bool {
	if o != nil && !IsNil(o.Repositories) {
		return true
	}

	return false
}

// SetRepositories gets a reference to the given []string and assigns it to the Repositories field.
func (o *PatchedPulpExporter) SetRepositories(v []string) {
	o.Repositories = v
}

// GetLastExport returns the LastExport field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PatchedPulpExporter) GetLastExport() string {
	if o == nil || IsNil(o.LastExport.Get()) {
		var ret string
		return ret
	}
	return *o.LastExport.Get()
}

// GetLastExportOk returns a tuple with the LastExport field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PatchedPulpExporter) GetLastExportOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.LastExport.Get(), o.LastExport.IsSet()
}

// HasLastExport returns a boolean if a field has been set.
func (o *PatchedPulpExporter) HasLastExport() bool {
	if o != nil && o.LastExport.IsSet() {
		return true
	}

	return false
}

// SetLastExport gets a reference to the given NullableString and assigns it to the LastExport field.
func (o *PatchedPulpExporter) SetLastExport(v string) {
	o.LastExport.Set(&v)
}
// SetLastExportNil sets the value for LastExport to be an explicit nil
func (o *PatchedPulpExporter) SetLastExportNil() {
	o.LastExport.Set(nil)
}

// UnsetLastExport ensures that no value is present for LastExport, not even an explicit nil
func (o *PatchedPulpExporter) UnsetLastExport() {
	o.LastExport.Unset()
}

func (o PatchedPulpExporter) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PatchedPulpExporter) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Path) {
		toSerialize["path"] = o.Path
	}
	if !IsNil(o.Repositories) {
		toSerialize["repositories"] = o.Repositories
	}
	if o.LastExport.IsSet() {
		toSerialize["last_export"] = o.LastExport.Get()
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *PatchedPulpExporter) UnmarshalJSON(bytes []byte) (err error) {
	varPatchedPulpExporter := _PatchedPulpExporter{}

	if err = json.Unmarshal(bytes, &varPatchedPulpExporter); err == nil {
		*o = PatchedPulpExporter(varPatchedPulpExporter)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "path")
		delete(additionalProperties, "repositories")
		delete(additionalProperties, "last_export")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePatchedPulpExporter struct {
	value *PatchedPulpExporter
	isSet bool
}

func (v NullablePatchedPulpExporter) Get() *PatchedPulpExporter {
	return v.value
}

func (v *NullablePatchedPulpExporter) Set(val *PatchedPulpExporter) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedPulpExporter) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedPulpExporter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedPulpExporter(val *PatchedPulpExporter) *NullablePatchedPulpExporter {
	return &NullablePatchedPulpExporter{value: val, isSet: true}
}

func (v NullablePatchedPulpExporter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedPulpExporter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


