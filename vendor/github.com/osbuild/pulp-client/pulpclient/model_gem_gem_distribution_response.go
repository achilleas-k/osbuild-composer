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

// checks if the GemGemDistributionResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GemGemDistributionResponse{}

// GemGemDistributionResponse A Serializer for GemDistribution.
type GemGemDistributionResponse struct {
	PulpHref *string `json:"pulp_href,omitempty"`
	// Timestamp of creation.
	PulpCreated *time.Time `json:"pulp_created,omitempty"`
	// The base (relative) path component of the published url. Avoid paths that                     overlap with other distribution base paths (e.g. \"foo\" and \"foo/bar\")
	BasePath string `json:"base_path"`
	// The URL for accessing the publication as defined by this distribution.
	BaseUrl *string `json:"base_url,omitempty"`
	// An optional content-guard.
	ContentGuard NullableString `json:"content_guard,omitempty"`
	// Whether this distribution should be shown in the content app.
	Hidden *bool `json:"hidden,omitempty"`
	PulpLabels *map[string]string `json:"pulp_labels,omitempty"`
	// A unique name. Ex, `rawhide` and `stable`.
	Name string `json:"name"`
	// The latest RepositoryVersion for this Repository will be served.
	Repository NullableString `json:"repository,omitempty"`
	// Publication to be served
	Publication NullableString `json:"publication,omitempty"`
	// Remote that can be used to fetch content when using pull-through caching.
	Remote NullableString `json:"remote,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _GemGemDistributionResponse GemGemDistributionResponse

// NewGemGemDistributionResponse instantiates a new GemGemDistributionResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGemGemDistributionResponse(basePath string, name string) *GemGemDistributionResponse {
	this := GemGemDistributionResponse{}
	this.BasePath = basePath
	var hidden bool = false
	this.Hidden = &hidden
	this.Name = name
	return &this
}

// NewGemGemDistributionResponseWithDefaults instantiates a new GemGemDistributionResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGemGemDistributionResponseWithDefaults() *GemGemDistributionResponse {
	this := GemGemDistributionResponse{}
	var hidden bool = false
	this.Hidden = &hidden
	return &this
}

// GetPulpHref returns the PulpHref field value if set, zero value otherwise.
func (o *GemGemDistributionResponse) GetPulpHref() string {
	if o == nil || IsNil(o.PulpHref) {
		var ret string
		return ret
	}
	return *o.PulpHref
}

// GetPulpHrefOk returns a tuple with the PulpHref field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GemGemDistributionResponse) GetPulpHrefOk() (*string, bool) {
	if o == nil || IsNil(o.PulpHref) {
		return nil, false
	}
	return o.PulpHref, true
}

// HasPulpHref returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasPulpHref() bool {
	if o != nil && !IsNil(o.PulpHref) {
		return true
	}

	return false
}

// SetPulpHref gets a reference to the given string and assigns it to the PulpHref field.
func (o *GemGemDistributionResponse) SetPulpHref(v string) {
	o.PulpHref = &v
}

// GetPulpCreated returns the PulpCreated field value if set, zero value otherwise.
func (o *GemGemDistributionResponse) GetPulpCreated() time.Time {
	if o == nil || IsNil(o.PulpCreated) {
		var ret time.Time
		return ret
	}
	return *o.PulpCreated
}

// GetPulpCreatedOk returns a tuple with the PulpCreated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GemGemDistributionResponse) GetPulpCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.PulpCreated) {
		return nil, false
	}
	return o.PulpCreated, true
}

// HasPulpCreated returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasPulpCreated() bool {
	if o != nil && !IsNil(o.PulpCreated) {
		return true
	}

	return false
}

// SetPulpCreated gets a reference to the given time.Time and assigns it to the PulpCreated field.
func (o *GemGemDistributionResponse) SetPulpCreated(v time.Time) {
	o.PulpCreated = &v
}

// GetBasePath returns the BasePath field value
func (o *GemGemDistributionResponse) GetBasePath() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BasePath
}

// GetBasePathOk returns a tuple with the BasePath field value
// and a boolean to check if the value has been set.
func (o *GemGemDistributionResponse) GetBasePathOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BasePath, true
}

// SetBasePath sets field value
func (o *GemGemDistributionResponse) SetBasePath(v string) {
	o.BasePath = v
}

// GetBaseUrl returns the BaseUrl field value if set, zero value otherwise.
func (o *GemGemDistributionResponse) GetBaseUrl() string {
	if o == nil || IsNil(o.BaseUrl) {
		var ret string
		return ret
	}
	return *o.BaseUrl
}

// GetBaseUrlOk returns a tuple with the BaseUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GemGemDistributionResponse) GetBaseUrlOk() (*string, bool) {
	if o == nil || IsNil(o.BaseUrl) {
		return nil, false
	}
	return o.BaseUrl, true
}

// HasBaseUrl returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasBaseUrl() bool {
	if o != nil && !IsNil(o.BaseUrl) {
		return true
	}

	return false
}

// SetBaseUrl gets a reference to the given string and assigns it to the BaseUrl field.
func (o *GemGemDistributionResponse) SetBaseUrl(v string) {
	o.BaseUrl = &v
}

// GetContentGuard returns the ContentGuard field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *GemGemDistributionResponse) GetContentGuard() string {
	if o == nil || IsNil(o.ContentGuard.Get()) {
		var ret string
		return ret
	}
	return *o.ContentGuard.Get()
}

// GetContentGuardOk returns a tuple with the ContentGuard field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GemGemDistributionResponse) GetContentGuardOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ContentGuard.Get(), o.ContentGuard.IsSet()
}

// HasContentGuard returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasContentGuard() bool {
	if o != nil && o.ContentGuard.IsSet() {
		return true
	}

	return false
}

// SetContentGuard gets a reference to the given NullableString and assigns it to the ContentGuard field.
func (o *GemGemDistributionResponse) SetContentGuard(v string) {
	o.ContentGuard.Set(&v)
}
// SetContentGuardNil sets the value for ContentGuard to be an explicit nil
func (o *GemGemDistributionResponse) SetContentGuardNil() {
	o.ContentGuard.Set(nil)
}

// UnsetContentGuard ensures that no value is present for ContentGuard, not even an explicit nil
func (o *GemGemDistributionResponse) UnsetContentGuard() {
	o.ContentGuard.Unset()
}

// GetHidden returns the Hidden field value if set, zero value otherwise.
func (o *GemGemDistributionResponse) GetHidden() bool {
	if o == nil || IsNil(o.Hidden) {
		var ret bool
		return ret
	}
	return *o.Hidden
}

// GetHiddenOk returns a tuple with the Hidden field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GemGemDistributionResponse) GetHiddenOk() (*bool, bool) {
	if o == nil || IsNil(o.Hidden) {
		return nil, false
	}
	return o.Hidden, true
}

// HasHidden returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasHidden() bool {
	if o != nil && !IsNil(o.Hidden) {
		return true
	}

	return false
}

// SetHidden gets a reference to the given bool and assigns it to the Hidden field.
func (o *GemGemDistributionResponse) SetHidden(v bool) {
	o.Hidden = &v
}

// GetPulpLabels returns the PulpLabels field value if set, zero value otherwise.
func (o *GemGemDistributionResponse) GetPulpLabels() map[string]string {
	if o == nil || IsNil(o.PulpLabels) {
		var ret map[string]string
		return ret
	}
	return *o.PulpLabels
}

// GetPulpLabelsOk returns a tuple with the PulpLabels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GemGemDistributionResponse) GetPulpLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.PulpLabels) {
		return nil, false
	}
	return o.PulpLabels, true
}

// HasPulpLabels returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasPulpLabels() bool {
	if o != nil && !IsNil(o.PulpLabels) {
		return true
	}

	return false
}

// SetPulpLabels gets a reference to the given map[string]string and assigns it to the PulpLabels field.
func (o *GemGemDistributionResponse) SetPulpLabels(v map[string]string) {
	o.PulpLabels = &v
}

// GetName returns the Name field value
func (o *GemGemDistributionResponse) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *GemGemDistributionResponse) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *GemGemDistributionResponse) SetName(v string) {
	o.Name = v
}

// GetRepository returns the Repository field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *GemGemDistributionResponse) GetRepository() string {
	if o == nil || IsNil(o.Repository.Get()) {
		var ret string
		return ret
	}
	return *o.Repository.Get()
}

// GetRepositoryOk returns a tuple with the Repository field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GemGemDistributionResponse) GetRepositoryOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Repository.Get(), o.Repository.IsSet()
}

// HasRepository returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasRepository() bool {
	if o != nil && o.Repository.IsSet() {
		return true
	}

	return false
}

// SetRepository gets a reference to the given NullableString and assigns it to the Repository field.
func (o *GemGemDistributionResponse) SetRepository(v string) {
	o.Repository.Set(&v)
}
// SetRepositoryNil sets the value for Repository to be an explicit nil
func (o *GemGemDistributionResponse) SetRepositoryNil() {
	o.Repository.Set(nil)
}

// UnsetRepository ensures that no value is present for Repository, not even an explicit nil
func (o *GemGemDistributionResponse) UnsetRepository() {
	o.Repository.Unset()
}

// GetPublication returns the Publication field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *GemGemDistributionResponse) GetPublication() string {
	if o == nil || IsNil(o.Publication.Get()) {
		var ret string
		return ret
	}
	return *o.Publication.Get()
}

// GetPublicationOk returns a tuple with the Publication field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GemGemDistributionResponse) GetPublicationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Publication.Get(), o.Publication.IsSet()
}

// HasPublication returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasPublication() bool {
	if o != nil && o.Publication.IsSet() {
		return true
	}

	return false
}

// SetPublication gets a reference to the given NullableString and assigns it to the Publication field.
func (o *GemGemDistributionResponse) SetPublication(v string) {
	o.Publication.Set(&v)
}
// SetPublicationNil sets the value for Publication to be an explicit nil
func (o *GemGemDistributionResponse) SetPublicationNil() {
	o.Publication.Set(nil)
}

// UnsetPublication ensures that no value is present for Publication, not even an explicit nil
func (o *GemGemDistributionResponse) UnsetPublication() {
	o.Publication.Unset()
}

// GetRemote returns the Remote field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *GemGemDistributionResponse) GetRemote() string {
	if o == nil || IsNil(o.Remote.Get()) {
		var ret string
		return ret
	}
	return *o.Remote.Get()
}

// GetRemoteOk returns a tuple with the Remote field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GemGemDistributionResponse) GetRemoteOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Remote.Get(), o.Remote.IsSet()
}

// HasRemote returns a boolean if a field has been set.
func (o *GemGemDistributionResponse) HasRemote() bool {
	if o != nil && o.Remote.IsSet() {
		return true
	}

	return false
}

// SetRemote gets a reference to the given NullableString and assigns it to the Remote field.
func (o *GemGemDistributionResponse) SetRemote(v string) {
	o.Remote.Set(&v)
}
// SetRemoteNil sets the value for Remote to be an explicit nil
func (o *GemGemDistributionResponse) SetRemoteNil() {
	o.Remote.Set(nil)
}

// UnsetRemote ensures that no value is present for Remote, not even an explicit nil
func (o *GemGemDistributionResponse) UnsetRemote() {
	o.Remote.Unset()
}

func (o GemGemDistributionResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GemGemDistributionResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.PulpHref) {
		toSerialize["pulp_href"] = o.PulpHref
	}
	if !IsNil(o.PulpCreated) {
		toSerialize["pulp_created"] = o.PulpCreated
	}
	toSerialize["base_path"] = o.BasePath
	if !IsNil(o.BaseUrl) {
		toSerialize["base_url"] = o.BaseUrl
	}
	if o.ContentGuard.IsSet() {
		toSerialize["content_guard"] = o.ContentGuard.Get()
	}
	if !IsNil(o.Hidden) {
		toSerialize["hidden"] = o.Hidden
	}
	if !IsNil(o.PulpLabels) {
		toSerialize["pulp_labels"] = o.PulpLabels
	}
	toSerialize["name"] = o.Name
	if o.Repository.IsSet() {
		toSerialize["repository"] = o.Repository.Get()
	}
	if o.Publication.IsSet() {
		toSerialize["publication"] = o.Publication.Get()
	}
	if o.Remote.IsSet() {
		toSerialize["remote"] = o.Remote.Get()
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *GemGemDistributionResponse) UnmarshalJSON(bytes []byte) (err error) {
	varGemGemDistributionResponse := _GemGemDistributionResponse{}

	if err = json.Unmarshal(bytes, &varGemGemDistributionResponse); err == nil {
		*o = GemGemDistributionResponse(varGemGemDistributionResponse)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "pulp_href")
		delete(additionalProperties, "pulp_created")
		delete(additionalProperties, "base_path")
		delete(additionalProperties, "base_url")
		delete(additionalProperties, "content_guard")
		delete(additionalProperties, "hidden")
		delete(additionalProperties, "pulp_labels")
		delete(additionalProperties, "name")
		delete(additionalProperties, "repository")
		delete(additionalProperties, "publication")
		delete(additionalProperties, "remote")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableGemGemDistributionResponse struct {
	value *GemGemDistributionResponse
	isSet bool
}

func (v NullableGemGemDistributionResponse) Get() *GemGemDistributionResponse {
	return v.value
}

func (v *NullableGemGemDistributionResponse) Set(val *GemGemDistributionResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGemGemDistributionResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGemGemDistributionResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGemGemDistributionResponse(val *GemGemDistributionResponse) *NullableGemGemDistributionResponse {
	return &NullableGemGemDistributionResponse{value: val, isSet: true}
}

func (v NullableGemGemDistributionResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGemGemDistributionResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


