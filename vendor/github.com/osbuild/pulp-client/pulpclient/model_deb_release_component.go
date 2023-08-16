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

// checks if the DebReleaseComponent type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DebReleaseComponent{}

// DebReleaseComponent A Serializer for ReleaseComponent.
type DebReleaseComponent struct {
	// A URI of a repository the new content unit should be associated with.
	Repository *string `json:"repository,omitempty"`
	// Name of the component.
	Component string `json:"component"`
	// Name of the distribution.
	Distribution string `json:"distribution"`
	Codename string `json:"codename"`
	Suite string `json:"suite"`
	AdditionalProperties map[string]interface{}
}

type _DebReleaseComponent DebReleaseComponent

// NewDebReleaseComponent instantiates a new DebReleaseComponent object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDebReleaseComponent(component string, distribution string, codename string, suite string) *DebReleaseComponent {
	this := DebReleaseComponent{}
	this.Component = component
	this.Distribution = distribution
	this.Codename = codename
	this.Suite = suite
	return &this
}

// NewDebReleaseComponentWithDefaults instantiates a new DebReleaseComponent object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDebReleaseComponentWithDefaults() *DebReleaseComponent {
	this := DebReleaseComponent{}
	return &this
}

// GetRepository returns the Repository field value if set, zero value otherwise.
func (o *DebReleaseComponent) GetRepository() string {
	if o == nil || IsNil(o.Repository) {
		var ret string
		return ret
	}
	return *o.Repository
}

// GetRepositoryOk returns a tuple with the Repository field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DebReleaseComponent) GetRepositoryOk() (*string, bool) {
	if o == nil || IsNil(o.Repository) {
		return nil, false
	}
	return o.Repository, true
}

// HasRepository returns a boolean if a field has been set.
func (o *DebReleaseComponent) HasRepository() bool {
	if o != nil && !IsNil(o.Repository) {
		return true
	}

	return false
}

// SetRepository gets a reference to the given string and assigns it to the Repository field.
func (o *DebReleaseComponent) SetRepository(v string) {
	o.Repository = &v
}

// GetComponent returns the Component field value
func (o *DebReleaseComponent) GetComponent() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Component
}

// GetComponentOk returns a tuple with the Component field value
// and a boolean to check if the value has been set.
func (o *DebReleaseComponent) GetComponentOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Component, true
}

// SetComponent sets field value
func (o *DebReleaseComponent) SetComponent(v string) {
	o.Component = v
}

// GetDistribution returns the Distribution field value
func (o *DebReleaseComponent) GetDistribution() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Distribution
}

// GetDistributionOk returns a tuple with the Distribution field value
// and a boolean to check if the value has been set.
func (o *DebReleaseComponent) GetDistributionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Distribution, true
}

// SetDistribution sets field value
func (o *DebReleaseComponent) SetDistribution(v string) {
	o.Distribution = v
}

// GetCodename returns the Codename field value
func (o *DebReleaseComponent) GetCodename() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Codename
}

// GetCodenameOk returns a tuple with the Codename field value
// and a boolean to check if the value has been set.
func (o *DebReleaseComponent) GetCodenameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Codename, true
}

// SetCodename sets field value
func (o *DebReleaseComponent) SetCodename(v string) {
	o.Codename = v
}

// GetSuite returns the Suite field value
func (o *DebReleaseComponent) GetSuite() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Suite
}

// GetSuiteOk returns a tuple with the Suite field value
// and a boolean to check if the value has been set.
func (o *DebReleaseComponent) GetSuiteOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Suite, true
}

// SetSuite sets field value
func (o *DebReleaseComponent) SetSuite(v string) {
	o.Suite = v
}

func (o DebReleaseComponent) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DebReleaseComponent) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Repository) {
		toSerialize["repository"] = o.Repository
	}
	toSerialize["component"] = o.Component
	toSerialize["distribution"] = o.Distribution
	toSerialize["codename"] = o.Codename
	toSerialize["suite"] = o.Suite

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *DebReleaseComponent) UnmarshalJSON(bytes []byte) (err error) {
	varDebReleaseComponent := _DebReleaseComponent{}

	if err = json.Unmarshal(bytes, &varDebReleaseComponent); err == nil {
		*o = DebReleaseComponent(varDebReleaseComponent)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "repository")
		delete(additionalProperties, "component")
		delete(additionalProperties, "distribution")
		delete(additionalProperties, "codename")
		delete(additionalProperties, "suite")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDebReleaseComponent struct {
	value *DebReleaseComponent
	isSet bool
}

func (v NullableDebReleaseComponent) Get() *DebReleaseComponent {
	return v.value
}

func (v *NullableDebReleaseComponent) Set(val *DebReleaseComponent) {
	v.value = val
	v.isSet = true
}

func (v NullableDebReleaseComponent) IsSet() bool {
	return v.isSet
}

func (v *NullableDebReleaseComponent) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDebReleaseComponent(val *DebReleaseComponent) *NullableDebReleaseComponent {
	return &NullableDebReleaseComponent{value: val, isSet: true}
}

func (v NullableDebReleaseComponent) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDebReleaseComponent) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


