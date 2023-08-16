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

// checks if the NamespaceLink type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NamespaceLink{}

// NamespaceLink Provides backwards compatible interface for links with the legacy GalaxyNG API.
type NamespaceLink struct {
	Url string `json:"url"`
	Name string `json:"name"`
	AdditionalProperties map[string]interface{}
}

type _NamespaceLink NamespaceLink

// NewNamespaceLink instantiates a new NamespaceLink object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNamespaceLink(url string, name string) *NamespaceLink {
	this := NamespaceLink{}
	this.Url = url
	this.Name = name
	return &this
}

// NewNamespaceLinkWithDefaults instantiates a new NamespaceLink object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNamespaceLinkWithDefaults() *NamespaceLink {
	this := NamespaceLink{}
	return &this
}

// GetUrl returns the Url field value
func (o *NamespaceLink) GetUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Url
}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
func (o *NamespaceLink) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Url, true
}

// SetUrl sets field value
func (o *NamespaceLink) SetUrl(v string) {
	o.Url = v
}

// GetName returns the Name field value
func (o *NamespaceLink) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *NamespaceLink) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *NamespaceLink) SetName(v string) {
	o.Name = v
}

func (o NamespaceLink) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NamespaceLink) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["url"] = o.Url
	toSerialize["name"] = o.Name

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *NamespaceLink) UnmarshalJSON(bytes []byte) (err error) {
	varNamespaceLink := _NamespaceLink{}

	if err = json.Unmarshal(bytes, &varNamespaceLink); err == nil {
		*o = NamespaceLink(varNamespaceLink)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "url")
		delete(additionalProperties, "name")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableNamespaceLink struct {
	value *NamespaceLink
	isSet bool
}

func (v NullableNamespaceLink) Get() *NamespaceLink {
	return v.value
}

func (v *NullableNamespaceLink) Set(val *NamespaceLink) {
	v.value = val
	v.isSet = true
}

func (v NullableNamespaceLink) IsSet() bool {
	return v.isSet
}

func (v *NullableNamespaceLink) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNamespaceLink(val *NamespaceLink) *NullableNamespaceLink {
	return &NullableNamespaceLink{value: val, isSet: true}
}

func (v NullableNamespaceLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNamespaceLink) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


