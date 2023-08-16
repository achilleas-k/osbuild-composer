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

// checks if the CertguardRHSMCertGuardResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CertguardRHSMCertGuardResponse{}

// CertguardRHSMCertGuardResponse RHSM Content Guard Serializer.
type CertguardRHSMCertGuardResponse struct {
	PulpHref *string `json:"pulp_href,omitempty"`
	// Timestamp of creation.
	PulpCreated *time.Time `json:"pulp_created,omitempty"`
	// The unique name.
	Name string `json:"name"`
	// An optional description.
	Description NullableString `json:"description,omitempty"`
	// A Certificate Authority (CA) certificate (or a bundle thereof) used to verify client-certificate authenticity.
	CaCertificate string `json:"ca_certificate"`
	AdditionalProperties map[string]interface{}
}

type _CertguardRHSMCertGuardResponse CertguardRHSMCertGuardResponse

// NewCertguardRHSMCertGuardResponse instantiates a new CertguardRHSMCertGuardResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCertguardRHSMCertGuardResponse(name string, caCertificate string) *CertguardRHSMCertGuardResponse {
	this := CertguardRHSMCertGuardResponse{}
	this.Name = name
	this.CaCertificate = caCertificate
	return &this
}

// NewCertguardRHSMCertGuardResponseWithDefaults instantiates a new CertguardRHSMCertGuardResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCertguardRHSMCertGuardResponseWithDefaults() *CertguardRHSMCertGuardResponse {
	this := CertguardRHSMCertGuardResponse{}
	return &this
}

// GetPulpHref returns the PulpHref field value if set, zero value otherwise.
func (o *CertguardRHSMCertGuardResponse) GetPulpHref() string {
	if o == nil || IsNil(o.PulpHref) {
		var ret string
		return ret
	}
	return *o.PulpHref
}

// GetPulpHrefOk returns a tuple with the PulpHref field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CertguardRHSMCertGuardResponse) GetPulpHrefOk() (*string, bool) {
	if o == nil || IsNil(o.PulpHref) {
		return nil, false
	}
	return o.PulpHref, true
}

// HasPulpHref returns a boolean if a field has been set.
func (o *CertguardRHSMCertGuardResponse) HasPulpHref() bool {
	if o != nil && !IsNil(o.PulpHref) {
		return true
	}

	return false
}

// SetPulpHref gets a reference to the given string and assigns it to the PulpHref field.
func (o *CertguardRHSMCertGuardResponse) SetPulpHref(v string) {
	o.PulpHref = &v
}

// GetPulpCreated returns the PulpCreated field value if set, zero value otherwise.
func (o *CertguardRHSMCertGuardResponse) GetPulpCreated() time.Time {
	if o == nil || IsNil(o.PulpCreated) {
		var ret time.Time
		return ret
	}
	return *o.PulpCreated
}

// GetPulpCreatedOk returns a tuple with the PulpCreated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CertguardRHSMCertGuardResponse) GetPulpCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.PulpCreated) {
		return nil, false
	}
	return o.PulpCreated, true
}

// HasPulpCreated returns a boolean if a field has been set.
func (o *CertguardRHSMCertGuardResponse) HasPulpCreated() bool {
	if o != nil && !IsNil(o.PulpCreated) {
		return true
	}

	return false
}

// SetPulpCreated gets a reference to the given time.Time and assigns it to the PulpCreated field.
func (o *CertguardRHSMCertGuardResponse) SetPulpCreated(v time.Time) {
	o.PulpCreated = &v
}

// GetName returns the Name field value
func (o *CertguardRHSMCertGuardResponse) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *CertguardRHSMCertGuardResponse) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *CertguardRHSMCertGuardResponse) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CertguardRHSMCertGuardResponse) GetDescription() string {
	if o == nil || IsNil(o.Description.Get()) {
		var ret string
		return ret
	}
	return *o.Description.Get()
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertguardRHSMCertGuardResponse) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Description.Get(), o.Description.IsSet()
}

// HasDescription returns a boolean if a field has been set.
func (o *CertguardRHSMCertGuardResponse) HasDescription() bool {
	if o != nil && o.Description.IsSet() {
		return true
	}

	return false
}

// SetDescription gets a reference to the given NullableString and assigns it to the Description field.
func (o *CertguardRHSMCertGuardResponse) SetDescription(v string) {
	o.Description.Set(&v)
}
// SetDescriptionNil sets the value for Description to be an explicit nil
func (o *CertguardRHSMCertGuardResponse) SetDescriptionNil() {
	o.Description.Set(nil)
}

// UnsetDescription ensures that no value is present for Description, not even an explicit nil
func (o *CertguardRHSMCertGuardResponse) UnsetDescription() {
	o.Description.Unset()
}

// GetCaCertificate returns the CaCertificate field value
func (o *CertguardRHSMCertGuardResponse) GetCaCertificate() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CaCertificate
}

// GetCaCertificateOk returns a tuple with the CaCertificate field value
// and a boolean to check if the value has been set.
func (o *CertguardRHSMCertGuardResponse) GetCaCertificateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CaCertificate, true
}

// SetCaCertificate sets field value
func (o *CertguardRHSMCertGuardResponse) SetCaCertificate(v string) {
	o.CaCertificate = v
}

func (o CertguardRHSMCertGuardResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CertguardRHSMCertGuardResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.PulpHref) {
		toSerialize["pulp_href"] = o.PulpHref
	}
	if !IsNil(o.PulpCreated) {
		toSerialize["pulp_created"] = o.PulpCreated
	}
	toSerialize["name"] = o.Name
	if o.Description.IsSet() {
		toSerialize["description"] = o.Description.Get()
	}
	toSerialize["ca_certificate"] = o.CaCertificate

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CertguardRHSMCertGuardResponse) UnmarshalJSON(bytes []byte) (err error) {
	varCertguardRHSMCertGuardResponse := _CertguardRHSMCertGuardResponse{}

	if err = json.Unmarshal(bytes, &varCertguardRHSMCertGuardResponse); err == nil {
		*o = CertguardRHSMCertGuardResponse(varCertguardRHSMCertGuardResponse)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "pulp_href")
		delete(additionalProperties, "pulp_created")
		delete(additionalProperties, "name")
		delete(additionalProperties, "description")
		delete(additionalProperties, "ca_certificate")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCertguardRHSMCertGuardResponse struct {
	value *CertguardRHSMCertGuardResponse
	isSet bool
}

func (v NullableCertguardRHSMCertGuardResponse) Get() *CertguardRHSMCertGuardResponse {
	return v.value
}

func (v *NullableCertguardRHSMCertGuardResponse) Set(val *CertguardRHSMCertGuardResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableCertguardRHSMCertGuardResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableCertguardRHSMCertGuardResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCertguardRHSMCertGuardResponse(val *CertguardRHSMCertGuardResponse) *NullableCertguardRHSMCertGuardResponse {
	return &NullableCertguardRHSMCertGuardResponse{value: val, isSet: true}
}

func (v NullableCertguardRHSMCertGuardResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCertguardRHSMCertGuardResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


