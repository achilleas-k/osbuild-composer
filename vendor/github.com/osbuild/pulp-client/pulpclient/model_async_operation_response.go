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

// checks if the AsyncOperationResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AsyncOperationResponse{}

// AsyncOperationResponse Serializer for asynchronous operations.
type AsyncOperationResponse struct {
	// The href of the task.
	Task string `json:"task"`
	AdditionalProperties map[string]interface{}
}

type _AsyncOperationResponse AsyncOperationResponse

// NewAsyncOperationResponse instantiates a new AsyncOperationResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAsyncOperationResponse(task string) *AsyncOperationResponse {
	this := AsyncOperationResponse{}
	this.Task = task
	return &this
}

// NewAsyncOperationResponseWithDefaults instantiates a new AsyncOperationResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAsyncOperationResponseWithDefaults() *AsyncOperationResponse {
	this := AsyncOperationResponse{}
	return &this
}

// GetTask returns the Task field value
func (o *AsyncOperationResponse) GetTask() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Task
}

// GetTaskOk returns a tuple with the Task field value
// and a boolean to check if the value has been set.
func (o *AsyncOperationResponse) GetTaskOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Task, true
}

// SetTask sets field value
func (o *AsyncOperationResponse) SetTask(v string) {
	o.Task = v
}

func (o AsyncOperationResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AsyncOperationResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["task"] = o.Task

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *AsyncOperationResponse) UnmarshalJSON(bytes []byte) (err error) {
	varAsyncOperationResponse := _AsyncOperationResponse{}

	if err = json.Unmarshal(bytes, &varAsyncOperationResponse); err == nil {
		*o = AsyncOperationResponse(varAsyncOperationResponse)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "task")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableAsyncOperationResponse struct {
	value *AsyncOperationResponse
	isSet bool
}

func (v NullableAsyncOperationResponse) Get() *AsyncOperationResponse {
	return v.value
}

func (v *NullableAsyncOperationResponse) Set(val *AsyncOperationResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableAsyncOperationResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableAsyncOperationResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAsyncOperationResponse(val *AsyncOperationResponse) *NullableAsyncOperationResponse {
	return &NullableAsyncOperationResponse{value: val, isSet: true}
}

func (v NullableAsyncOperationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAsyncOperationResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


