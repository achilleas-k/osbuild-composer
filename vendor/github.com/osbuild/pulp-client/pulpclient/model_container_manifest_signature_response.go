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
	"time"
)

// checks if the ContainerManifestSignatureResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ContainerManifestSignatureResponse{}

// ContainerManifestSignatureResponse Serializer for image manifest signatures.
type ContainerManifestSignatureResponse struct {
	PulpHref *string `json:"pulp_href,omitempty"`
	// Timestamp of creation.
	PulpCreated *time.Time `json:"pulp_created,omitempty"`
	// Signature name in the format of `digest_algo:manifest_digest@random_32_chars`
	Name string `json:"name"`
	// sha256 digest of the signature blob
	Digest string `json:"digest"`
	// Container signature type, e.g. 'atomic'
	Type string `json:"type"`
	// Signing key ID
	KeyId string `json:"key_id"`
	// Timestamp of a signature
	Timestamp int64 `json:"timestamp"`
	// Signature creator
	Creator string `json:"creator"`
	// Manifest that is signed
	SignedManifest string `json:"signed_manifest"`
	AdditionalProperties map[string]interface{}
}

type _ContainerManifestSignatureResponse ContainerManifestSignatureResponse

// NewContainerManifestSignatureResponse instantiates a new ContainerManifestSignatureResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewContainerManifestSignatureResponse(name string, digest string, type_ string, keyId string, timestamp int64, creator string, signedManifest string) *ContainerManifestSignatureResponse {
	this := ContainerManifestSignatureResponse{}
	this.Name = name
	this.Digest = digest
	this.Type = type_
	this.KeyId = keyId
	this.Timestamp = timestamp
	this.Creator = creator
	this.SignedManifest = signedManifest
	return &this
}

// NewContainerManifestSignatureResponseWithDefaults instantiates a new ContainerManifestSignatureResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewContainerManifestSignatureResponseWithDefaults() *ContainerManifestSignatureResponse {
	this := ContainerManifestSignatureResponse{}
	return &this
}

// GetPulpHref returns the PulpHref field value if set, zero value otherwise.
func (o *ContainerManifestSignatureResponse) GetPulpHref() string {
	if o == nil || IsNil(o.PulpHref) {
		var ret string
		return ret
	}
	return *o.PulpHref
}

// GetPulpHrefOk returns a tuple with the PulpHref field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetPulpHrefOk() (*string, bool) {
	if o == nil || IsNil(o.PulpHref) {
		return nil, false
	}
	return o.PulpHref, true
}

// HasPulpHref returns a boolean if a field has been set.
func (o *ContainerManifestSignatureResponse) HasPulpHref() bool {
	if o != nil && !IsNil(o.PulpHref) {
		return true
	}

	return false
}

// SetPulpHref gets a reference to the given string and assigns it to the PulpHref field.
func (o *ContainerManifestSignatureResponse) SetPulpHref(v string) {
	o.PulpHref = &v
}

// GetPulpCreated returns the PulpCreated field value if set, zero value otherwise.
func (o *ContainerManifestSignatureResponse) GetPulpCreated() time.Time {
	if o == nil || IsNil(o.PulpCreated) {
		var ret time.Time
		return ret
	}
	return *o.PulpCreated
}

// GetPulpCreatedOk returns a tuple with the PulpCreated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetPulpCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.PulpCreated) {
		return nil, false
	}
	return o.PulpCreated, true
}

// HasPulpCreated returns a boolean if a field has been set.
func (o *ContainerManifestSignatureResponse) HasPulpCreated() bool {
	if o != nil && !IsNil(o.PulpCreated) {
		return true
	}

	return false
}

// SetPulpCreated gets a reference to the given time.Time and assigns it to the PulpCreated field.
func (o *ContainerManifestSignatureResponse) SetPulpCreated(v time.Time) {
	o.PulpCreated = &v
}

// GetName returns the Name field value
func (o *ContainerManifestSignatureResponse) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *ContainerManifestSignatureResponse) SetName(v string) {
	o.Name = v
}

// GetDigest returns the Digest field value
func (o *ContainerManifestSignatureResponse) GetDigest() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Digest
}

// GetDigestOk returns a tuple with the Digest field value
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetDigestOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Digest, true
}

// SetDigest sets field value
func (o *ContainerManifestSignatureResponse) SetDigest(v string) {
	o.Digest = v
}

// GetType returns the Type field value
func (o *ContainerManifestSignatureResponse) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *ContainerManifestSignatureResponse) SetType(v string) {
	o.Type = v
}

// GetKeyId returns the KeyId field value
func (o *ContainerManifestSignatureResponse) GetKeyId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.KeyId
}

// GetKeyIdOk returns a tuple with the KeyId field value
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetKeyIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.KeyId, true
}

// SetKeyId sets field value
func (o *ContainerManifestSignatureResponse) SetKeyId(v string) {
	o.KeyId = v
}

// GetTimestamp returns the Timestamp field value
func (o *ContainerManifestSignatureResponse) GetTimestamp() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetTimestampOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Timestamp, true
}

// SetTimestamp sets field value
func (o *ContainerManifestSignatureResponse) SetTimestamp(v int64) {
	o.Timestamp = v
}

// GetCreator returns the Creator field value
func (o *ContainerManifestSignatureResponse) GetCreator() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Creator
}

// GetCreatorOk returns a tuple with the Creator field value
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetCreatorOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Creator, true
}

// SetCreator sets field value
func (o *ContainerManifestSignatureResponse) SetCreator(v string) {
	o.Creator = v
}

// GetSignedManifest returns the SignedManifest field value
func (o *ContainerManifestSignatureResponse) GetSignedManifest() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SignedManifest
}

// GetSignedManifestOk returns a tuple with the SignedManifest field value
// and a boolean to check if the value has been set.
func (o *ContainerManifestSignatureResponse) GetSignedManifestOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SignedManifest, true
}

// SetSignedManifest sets field value
func (o *ContainerManifestSignatureResponse) SetSignedManifest(v string) {
	o.SignedManifest = v
}

func (o ContainerManifestSignatureResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ContainerManifestSignatureResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.PulpHref) {
		toSerialize["pulp_href"] = o.PulpHref
	}
	if !IsNil(o.PulpCreated) {
		toSerialize["pulp_created"] = o.PulpCreated
	}
	toSerialize["name"] = o.Name
	toSerialize["digest"] = o.Digest
	toSerialize["type"] = o.Type
	toSerialize["key_id"] = o.KeyId
	toSerialize["timestamp"] = o.Timestamp
	toSerialize["creator"] = o.Creator
	toSerialize["signed_manifest"] = o.SignedManifest

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ContainerManifestSignatureResponse) UnmarshalJSON(bytes []byte) (err error) {
	varContainerManifestSignatureResponse := _ContainerManifestSignatureResponse{}

	if err = json.Unmarshal(bytes, &varContainerManifestSignatureResponse); err == nil {
		*o = ContainerManifestSignatureResponse(varContainerManifestSignatureResponse)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "pulp_href")
		delete(additionalProperties, "pulp_created")
		delete(additionalProperties, "name")
		delete(additionalProperties, "digest")
		delete(additionalProperties, "type")
		delete(additionalProperties, "key_id")
		delete(additionalProperties, "timestamp")
		delete(additionalProperties, "creator")
		delete(additionalProperties, "signed_manifest")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableContainerManifestSignatureResponse struct {
	value *ContainerManifestSignatureResponse
	isSet bool
}

func (v NullableContainerManifestSignatureResponse) Get() *ContainerManifestSignatureResponse {
	return v.value
}

func (v *NullableContainerManifestSignatureResponse) Set(val *ContainerManifestSignatureResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableContainerManifestSignatureResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableContainerManifestSignatureResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableContainerManifestSignatureResponse(val *ContainerManifestSignatureResponse) *NullableContainerManifestSignatureResponse {
	return &NullableContainerManifestSignatureResponse{value: val, isSet: true}
}

func (v NullableContainerManifestSignatureResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableContainerManifestSignatureResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


