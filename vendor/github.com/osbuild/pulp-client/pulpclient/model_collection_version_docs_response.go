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

// checks if the CollectionVersionDocsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CollectionVersionDocsResponse{}

// CollectionVersionDocsResponse A serializer to display the docs_blob of a CollectionVersion.
type CollectionVersionDocsResponse struct {
	DocsBlob map[string]interface{} `json:"docs_blob"`
	AdditionalProperties map[string]interface{}
}

type _CollectionVersionDocsResponse CollectionVersionDocsResponse

// NewCollectionVersionDocsResponse instantiates a new CollectionVersionDocsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCollectionVersionDocsResponse(docsBlob map[string]interface{}) *CollectionVersionDocsResponse {
	this := CollectionVersionDocsResponse{}
	this.DocsBlob = docsBlob
	return &this
}

// NewCollectionVersionDocsResponseWithDefaults instantiates a new CollectionVersionDocsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCollectionVersionDocsResponseWithDefaults() *CollectionVersionDocsResponse {
	this := CollectionVersionDocsResponse{}
	return &this
}

// GetDocsBlob returns the DocsBlob field value
func (o *CollectionVersionDocsResponse) GetDocsBlob() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.DocsBlob
}

// GetDocsBlobOk returns a tuple with the DocsBlob field value
// and a boolean to check if the value has been set.
func (o *CollectionVersionDocsResponse) GetDocsBlobOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.DocsBlob, true
}

// SetDocsBlob sets field value
func (o *CollectionVersionDocsResponse) SetDocsBlob(v map[string]interface{}) {
	o.DocsBlob = v
}

func (o CollectionVersionDocsResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CollectionVersionDocsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["docs_blob"] = o.DocsBlob

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CollectionVersionDocsResponse) UnmarshalJSON(bytes []byte) (err error) {
	varCollectionVersionDocsResponse := _CollectionVersionDocsResponse{}

	if err = json.Unmarshal(bytes, &varCollectionVersionDocsResponse); err == nil {
		*o = CollectionVersionDocsResponse(varCollectionVersionDocsResponse)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "docs_blob")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCollectionVersionDocsResponse struct {
	value *CollectionVersionDocsResponse
	isSet bool
}

func (v NullableCollectionVersionDocsResponse) Get() *CollectionVersionDocsResponse {
	return v.value
}

func (v *NullableCollectionVersionDocsResponse) Set(val *CollectionVersionDocsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableCollectionVersionDocsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableCollectionVersionDocsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCollectionVersionDocsResponse(val *CollectionVersionDocsResponse) *NullableCollectionVersionDocsResponse {
	return &NullableCollectionVersionDocsResponse{value: val, isSet: true}
}

func (v NullableCollectionVersionDocsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCollectionVersionDocsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

