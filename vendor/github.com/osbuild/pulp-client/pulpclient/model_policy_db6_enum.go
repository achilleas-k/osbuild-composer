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
	"fmt"
)

// PolicyDb6Enum * `immediate` - immediate * `When syncing, download all metadata and content now.` - When syncing, download all metadata and content now.
type PolicyDb6Enum string

// List of PolicyDb6Enum
const (
	POLICYDB6ENUM_IMMEDIATE PolicyDb6Enum = "immediate"
	POLICYDB6ENUM_WHEN_SYNCING_DOWNLOAD_ALL_METADATA_AND_CONTENT_NOW PolicyDb6Enum = "When syncing, download all metadata and content now."
)

// All allowed values of PolicyDb6Enum enum
var AllowedPolicyDb6EnumEnumValues = []PolicyDb6Enum{
	"immediate",
	"When syncing, download all metadata and content now.",
}

func (v *PolicyDb6Enum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := PolicyDb6Enum(value)
	for _, existing := range AllowedPolicyDb6EnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid PolicyDb6Enum", value)
}

// NewPolicyDb6EnumFromValue returns a pointer to a valid PolicyDb6Enum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewPolicyDb6EnumFromValue(v string) (*PolicyDb6Enum, error) {
	ev := PolicyDb6Enum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for PolicyDb6Enum: valid values are %v", v, AllowedPolicyDb6EnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v PolicyDb6Enum) IsValid() bool {
	for _, existing := range AllowedPolicyDb6EnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to PolicyDb6Enum value
func (v PolicyDb6Enum) Ptr() *PolicyDb6Enum {
	return &v
}

type NullablePolicyDb6Enum struct {
	value *PolicyDb6Enum
	isSet bool
}

func (v NullablePolicyDb6Enum) Get() *PolicyDb6Enum {
	return v.value
}

func (v *NullablePolicyDb6Enum) Set(val *PolicyDb6Enum) {
	v.value = val
	v.isSet = true
}

func (v NullablePolicyDb6Enum) IsSet() bool {
	return v.isSet
}

func (v *NullablePolicyDb6Enum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePolicyDb6Enum(val *PolicyDb6Enum) *NullablePolicyDb6Enum {
	return &NullablePolicyDb6Enum{value: val, isSet: true}
}

func (v NullablePolicyDb6Enum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePolicyDb6Enum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

