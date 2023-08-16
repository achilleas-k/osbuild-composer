/*
Pulp 3 API

Fetch, Upload, Organize, and Distribute Software Packages

API version: v3
Contact: pulp-list@redhat.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package pulpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"reflect"
)


// PulpAnsibleDefaultApiV3NamespacesAPIService PulpAnsibleDefaultApiV3NamespacesAPI service
type PulpAnsibleDefaultApiV3NamespacesAPIService service

type PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest struct {
	ctx context.Context
	ApiService *PulpAnsibleDefaultApiV3NamespacesAPIService
	company *string
	companyContains *string
	companyIcontains *string
	companyIn *[]string
	companyStartswith *string
	limit *int32
	metadataSha256 *string
	metadataSha256In *[]string
	name *string
	nameContains *string
	nameIcontains *string
	nameIn *[]string
	nameStartswith *string
	offset *int32
	ordering *[]string
	pulpHrefIn *[]string
	pulpIdIn *[]string
	fields *[]string
	excludeFields *[]string
}

// Filter results where company matches value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) Company(company string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.company = &company
	return r
}

// Filter results where company contains value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) CompanyContains(companyContains string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.companyContains = &companyContains
	return r
}

// Filter results where company contains value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) CompanyIcontains(companyIcontains string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.companyIcontains = &companyIcontains
	return r
}

// Filter results where company is in a comma-separated list of values
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) CompanyIn(companyIn []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.companyIn = &companyIn
	return r
}

// Filter results where company starts with value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) CompanyStartswith(companyStartswith string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.companyStartswith = &companyStartswith
	return r
}

// Number of results to return per page.
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) Limit(limit int32) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.limit = &limit
	return r
}

// Filter results where metadata_sha256 matches value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) MetadataSha256(metadataSha256 string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.metadataSha256 = &metadataSha256
	return r
}

// Filter results where metadata_sha256 is in a comma-separated list of values
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) MetadataSha256In(metadataSha256In []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.metadataSha256In = &metadataSha256In
	return r
}

// Filter results where name matches value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) Name(name string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.name = &name
	return r
}

// Filter results where name contains value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) NameContains(nameContains string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.nameContains = &nameContains
	return r
}

// Filter results where name contains value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) NameIcontains(nameIcontains string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.nameIcontains = &nameIcontains
	return r
}

// Filter results where name is in a comma-separated list of values
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) NameIn(nameIn []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.nameIn = &nameIn
	return r
}

// Filter results where name starts with value
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) NameStartswith(nameStartswith string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.nameStartswith = &nameStartswith
	return r
}

// The initial index from which to return the results.
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) Offset(offset int32) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.offset = &offset
	return r
}

// Ordering  * &#x60;pulp_id&#x60; - Pulp id * &#x60;-pulp_id&#x60; - Pulp id (descending) * &#x60;pulp_created&#x60; - Pulp created * &#x60;-pulp_created&#x60; - Pulp created (descending) * &#x60;pulp_last_updated&#x60; - Pulp last updated * &#x60;-pulp_last_updated&#x60; - Pulp last updated (descending) * &#x60;pulp_type&#x60; - Pulp type * &#x60;-pulp_type&#x60; - Pulp type (descending) * &#x60;upstream_id&#x60; - Upstream id * &#x60;-upstream_id&#x60; - Upstream id (descending) * &#x60;timestamp_of_interest&#x60; - Timestamp of interest * &#x60;-timestamp_of_interest&#x60; - Timestamp of interest (descending) * &#x60;name&#x60; - Name * &#x60;-name&#x60; - Name (descending) * &#x60;company&#x60; - Company * &#x60;-company&#x60; - Company (descending) * &#x60;email&#x60; - Email * &#x60;-email&#x60; - Email (descending) * &#x60;description&#x60; - Description * &#x60;-description&#x60; - Description (descending) * &#x60;resources&#x60; - Resources * &#x60;-resources&#x60; - Resources (descending) * &#x60;links&#x60; - Links * &#x60;-links&#x60; - Links (descending) * &#x60;avatar_sha256&#x60; - Avatar sha256 * &#x60;-avatar_sha256&#x60; - Avatar sha256 (descending) * &#x60;metadata_sha256&#x60; - Metadata sha256 * &#x60;-metadata_sha256&#x60; - Metadata sha256 (descending) * &#x60;pk&#x60; - Pk * &#x60;-pk&#x60; - Pk (descending)
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) Ordering(ordering []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.ordering = &ordering
	return r
}

// Multiple values may be separated by commas.
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) PulpHrefIn(pulpHrefIn []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.pulpHrefIn = &pulpHrefIn
	return r
}

// Multiple values may be separated by commas.
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) PulpIdIn(pulpIdIn []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.pulpIdIn = &pulpIdIn
	return r
}

// A list of fields to include in the response.
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) Fields(fields []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.fields = &fields
	return r
}

// A list of fields to exclude from the response.
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) ExcludeFields(excludeFields []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	r.excludeFields = &excludeFields
	return r
}

func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) Execute() (*PaginatedansibleAnsibleNamespaceMetadataResponseList, *http.Response, error) {
	return r.ApiService.PulpAnsibleGalaxyDefaultApiV3NamespacesListExecute(r)
}

/*
PulpAnsibleGalaxyDefaultApiV3NamespacesList Method for PulpAnsibleGalaxyDefaultApiV3NamespacesList

Legacy v3 endpoint.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest

Deprecated
*/
func (a *PulpAnsibleDefaultApiV3NamespacesAPIService) PulpAnsibleGalaxyDefaultApiV3NamespacesList(ctx context.Context) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest {
	return PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return PaginatedansibleAnsibleNamespaceMetadataResponseList
// Deprecated
func (a *PulpAnsibleDefaultApiV3NamespacesAPIService) PulpAnsibleGalaxyDefaultApiV3NamespacesListExecute(r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesListRequest) (*PaginatedansibleAnsibleNamespaceMetadataResponseList, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *PaginatedansibleAnsibleNamespaceMetadataResponseList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PulpAnsibleDefaultApiV3NamespacesAPIService.PulpAnsibleGalaxyDefaultApiV3NamespacesList")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/pulp_ansible/galaxy/default/api/v3/namespaces/"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.company != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "company", r.company, "")
	}
	if r.companyContains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "company__contains", r.companyContains, "")
	}
	if r.companyIcontains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "company__icontains", r.companyIcontains, "")
	}
	if r.companyIn != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "company__in", r.companyIn, "csv")
	}
	if r.companyStartswith != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "company__startswith", r.companyStartswith, "")
	}
	if r.limit != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "limit", r.limit, "")
	}
	if r.metadataSha256 != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "metadata_sha256", r.metadataSha256, "")
	}
	if r.metadataSha256In != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "metadata_sha256__in", r.metadataSha256In, "csv")
	}
	if r.name != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name", r.name, "")
	}
	if r.nameContains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name__contains", r.nameContains, "")
	}
	if r.nameIcontains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name__icontains", r.nameIcontains, "")
	}
	if r.nameIn != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name__in", r.nameIn, "csv")
	}
	if r.nameStartswith != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name__startswith", r.nameStartswith, "")
	}
	if r.offset != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "offset", r.offset, "")
	}
	if r.ordering != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "ordering", r.ordering, "csv")
	}
	if r.pulpHrefIn != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pulp_href__in", r.pulpHrefIn, "csv")
	}
	if r.pulpIdIn != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pulp_id__in", r.pulpIdIn, "csv")
	}
	if r.fields != nil {
		t := *r.fields
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				parameterAddToHeaderOrQuery(localVarQueryParams, "fields", s.Index(i), "multi")
			}
		} else {
			parameterAddToHeaderOrQuery(localVarQueryParams, "fields", t, "multi")
		}
	}
	if r.excludeFields != nil {
		t := *r.excludeFields
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				parameterAddToHeaderOrQuery(localVarQueryParams, "exclude_fields", s.Index(i), "multi")
			}
		} else {
			parameterAddToHeaderOrQuery(localVarQueryParams, "exclude_fields", t, "multi")
		}
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest struct {
	ctx context.Context
	ApiService *PulpAnsibleDefaultApiV3NamespacesAPIService
	name string
	fields *[]string
	excludeFields *[]string
}

// A list of fields to include in the response.
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest) Fields(fields []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest {
	r.fields = &fields
	return r
}

// A list of fields to exclude from the response.
func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest) ExcludeFields(excludeFields []string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest {
	r.excludeFields = &excludeFields
	return r
}

func (r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest) Execute() (*AnsibleAnsibleNamespaceMetadataResponse, *http.Response, error) {
	return r.ApiService.PulpAnsibleGalaxyDefaultApiV3NamespacesReadExecute(r)
}

/*
PulpAnsibleGalaxyDefaultApiV3NamespacesRead Method for PulpAnsibleGalaxyDefaultApiV3NamespacesRead

Legacy v3 endpoint.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param name
 @return PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest

Deprecated
*/
func (a *PulpAnsibleDefaultApiV3NamespacesAPIService) PulpAnsibleGalaxyDefaultApiV3NamespacesRead(ctx context.Context, name string) PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest {
	return PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest{
		ApiService: a,
		ctx: ctx,
		name: name,
	}
}

// Execute executes the request
//  @return AnsibleAnsibleNamespaceMetadataResponse
// Deprecated
func (a *PulpAnsibleDefaultApiV3NamespacesAPIService) PulpAnsibleGalaxyDefaultApiV3NamespacesReadExecute(r PulpAnsibleDefaultApiV3NamespacesAPIPulpAnsibleGalaxyDefaultApiV3NamespacesReadRequest) (*AnsibleAnsibleNamespaceMetadataResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *AnsibleAnsibleNamespaceMetadataResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PulpAnsibleDefaultApiV3NamespacesAPIService.PulpAnsibleGalaxyDefaultApiV3NamespacesRead")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/pulp_ansible/galaxy/default/api/v3/namespaces/{name}/"
	localVarPath = strings.Replace(localVarPath, "{"+"name"+"}", parameterValueToString(r.name, "name"), -1)  // NOTE: paths aren't escaped because Pulp uses hrefs as path parameters

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.fields != nil {
		t := *r.fields
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				parameterAddToHeaderOrQuery(localVarQueryParams, "fields", s.Index(i), "multi")
			}
		} else {
			parameterAddToHeaderOrQuery(localVarQueryParams, "fields", t, "multi")
		}
	}
	if r.excludeFields != nil {
		t := *r.excludeFields
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				parameterAddToHeaderOrQuery(localVarQueryParams, "exclude_fields", s.Index(i), "multi")
			}
		} else {
			parameterAddToHeaderOrQuery(localVarQueryParams, "exclude_fields", t, "multi")
		}
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
