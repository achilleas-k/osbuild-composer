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


// UsersRolesAPIService UsersRolesAPI service
type UsersRolesAPIService service

type UsersRolesAPIUsersRolesCreateRequest struct {
	ctx context.Context
	ApiService *UsersRolesAPIService
	authUserHref string
	userRole *UserRole
}

func (r UsersRolesAPIUsersRolesCreateRequest) UserRole(userRole UserRole) UsersRolesAPIUsersRolesCreateRequest {
	r.userRole = &userRole
	return r
}

func (r UsersRolesAPIUsersRolesCreateRequest) Execute() (*UserRoleResponse, *http.Response, error) {
	return r.ApiService.UsersRolesCreateExecute(r)
}

/*
UsersRolesCreate Create an user role

ViewSet for UserRole.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param authUserHref
 @return UsersRolesAPIUsersRolesCreateRequest
*/
func (a *UsersRolesAPIService) UsersRolesCreate(ctx context.Context, authUserHref string) UsersRolesAPIUsersRolesCreateRequest {
	return UsersRolesAPIUsersRolesCreateRequest{
		ApiService: a,
		ctx: ctx,
		authUserHref: authUserHref,
	}
}

// Execute executes the request
//  @return UserRoleResponse
func (a *UsersRolesAPIService) UsersRolesCreateExecute(r UsersRolesAPIUsersRolesCreateRequest) (*UserRoleResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *UserRoleResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "UsersRolesAPIService.UsersRolesCreate")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "{auth_user_href}roles/"
	localVarPath = strings.Replace(localVarPath, "{"+"auth_user_href"+"}", parameterValueToString(r.authUserHref, "authUserHref"), -1)  // NOTE: paths aren't escaped because Pulp uses hrefs as path parameters

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.userRole == nil {
		return localVarReturnValue, nil, reportError("userRole is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json", "application/x-www-form-urlencoded", "multipart/form-data"}

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
	// body params
	localVarPostBody = r.userRole
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

type UsersRolesAPIUsersRolesDeleteRequest struct {
	ctx context.Context
	ApiService *UsersRolesAPIService
	authUsersUserRoleHref string
}

func (r UsersRolesAPIUsersRolesDeleteRequest) Execute() (*http.Response, error) {
	return r.ApiService.UsersRolesDeleteExecute(r)
}

/*
UsersRolesDelete Delete an user role

ViewSet for UserRole.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param authUsersUserRoleHref
 @return UsersRolesAPIUsersRolesDeleteRequest
*/
func (a *UsersRolesAPIService) UsersRolesDelete(ctx context.Context, authUsersUserRoleHref string) UsersRolesAPIUsersRolesDeleteRequest {
	return UsersRolesAPIUsersRolesDeleteRequest{
		ApiService: a,
		ctx: ctx,
		authUsersUserRoleHref: authUsersUserRoleHref,
	}
}

// Execute executes the request
func (a *UsersRolesAPIService) UsersRolesDeleteExecute(r UsersRolesAPIUsersRolesDeleteRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "UsersRolesAPIService.UsersRolesDelete")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "{auth_users_user_role_href}"
	localVarPath = strings.Replace(localVarPath, "{"+"auth_users_user_role_href"+"}", parameterValueToString(r.authUsersUserRoleHref, "authUsersUserRoleHref"), -1)  // NOTE: paths aren't escaped because Pulp uses hrefs as path parameters

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type UsersRolesAPIUsersRolesListRequest struct {
	ctx context.Context
	ApiService *UsersRolesAPIService
	authUserHref string
	contentObject *string
	domain *string
	limit *int32
	offset *int32
	ordering *[]string
	pulpHrefIn *[]string
	pulpIdIn *[]string
	role *string
	roleContains *string
	roleIcontains *string
	roleIn *[]string
	roleStartswith *string
	fields *[]string
	excludeFields *[]string
}

// content_object
func (r UsersRolesAPIUsersRolesListRequest) ContentObject(contentObject string) UsersRolesAPIUsersRolesListRequest {
	r.contentObject = &contentObject
	return r
}

// Foreign Key referenced by HREF
func (r UsersRolesAPIUsersRolesListRequest) Domain(domain string) UsersRolesAPIUsersRolesListRequest {
	r.domain = &domain
	return r
}

// Number of results to return per page.
func (r UsersRolesAPIUsersRolesListRequest) Limit(limit int32) UsersRolesAPIUsersRolesListRequest {
	r.limit = &limit
	return r
}

// The initial index from which to return the results.
func (r UsersRolesAPIUsersRolesListRequest) Offset(offset int32) UsersRolesAPIUsersRolesListRequest {
	r.offset = &offset
	return r
}

// Ordering  * &#x60;role&#x60; - Role * &#x60;-role&#x60; - Role (descending) * &#x60;description&#x60; - Description * &#x60;-description&#x60; - Description (descending) * &#x60;pulp_created&#x60; - Pulp created * &#x60;-pulp_created&#x60; - Pulp created (descending) * &#x60;pk&#x60; - Pk * &#x60;-pk&#x60; - Pk (descending)
func (r UsersRolesAPIUsersRolesListRequest) Ordering(ordering []string) UsersRolesAPIUsersRolesListRequest {
	r.ordering = &ordering
	return r
}

// Multiple values may be separated by commas.
func (r UsersRolesAPIUsersRolesListRequest) PulpHrefIn(pulpHrefIn []string) UsersRolesAPIUsersRolesListRequest {
	r.pulpHrefIn = &pulpHrefIn
	return r
}

// Multiple values may be separated by commas.
func (r UsersRolesAPIUsersRolesListRequest) PulpIdIn(pulpIdIn []string) UsersRolesAPIUsersRolesListRequest {
	r.pulpIdIn = &pulpIdIn
	return r
}

func (r UsersRolesAPIUsersRolesListRequest) Role(role string) UsersRolesAPIUsersRolesListRequest {
	r.role = &role
	return r
}

func (r UsersRolesAPIUsersRolesListRequest) RoleContains(roleContains string) UsersRolesAPIUsersRolesListRequest {
	r.roleContains = &roleContains
	return r
}

func (r UsersRolesAPIUsersRolesListRequest) RoleIcontains(roleIcontains string) UsersRolesAPIUsersRolesListRequest {
	r.roleIcontains = &roleIcontains
	return r
}

// Multiple values may be separated by commas.
func (r UsersRolesAPIUsersRolesListRequest) RoleIn(roleIn []string) UsersRolesAPIUsersRolesListRequest {
	r.roleIn = &roleIn
	return r
}

func (r UsersRolesAPIUsersRolesListRequest) RoleStartswith(roleStartswith string) UsersRolesAPIUsersRolesListRequest {
	r.roleStartswith = &roleStartswith
	return r
}

// A list of fields to include in the response.
func (r UsersRolesAPIUsersRolesListRequest) Fields(fields []string) UsersRolesAPIUsersRolesListRequest {
	r.fields = &fields
	return r
}

// A list of fields to exclude from the response.
func (r UsersRolesAPIUsersRolesListRequest) ExcludeFields(excludeFields []string) UsersRolesAPIUsersRolesListRequest {
	r.excludeFields = &excludeFields
	return r
}

func (r UsersRolesAPIUsersRolesListRequest) Execute() (*PaginatedUserRoleResponseList, *http.Response, error) {
	return r.ApiService.UsersRolesListExecute(r)
}

/*
UsersRolesList List user roles

ViewSet for UserRole.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param authUserHref
 @return UsersRolesAPIUsersRolesListRequest
*/
func (a *UsersRolesAPIService) UsersRolesList(ctx context.Context, authUserHref string) UsersRolesAPIUsersRolesListRequest {
	return UsersRolesAPIUsersRolesListRequest{
		ApiService: a,
		ctx: ctx,
		authUserHref: authUserHref,
	}
}

// Execute executes the request
//  @return PaginatedUserRoleResponseList
func (a *UsersRolesAPIService) UsersRolesListExecute(r UsersRolesAPIUsersRolesListRequest) (*PaginatedUserRoleResponseList, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *PaginatedUserRoleResponseList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "UsersRolesAPIService.UsersRolesList")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "{auth_user_href}roles/"
	localVarPath = strings.Replace(localVarPath, "{"+"auth_user_href"+"}", parameterValueToString(r.authUserHref, "authUserHref"), -1)  // NOTE: paths aren't escaped because Pulp uses hrefs as path parameters

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.contentObject != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "content_object", r.contentObject, "")
	}
	if r.domain != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "domain", r.domain, "")
	}
	if r.limit != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "limit", r.limit, "")
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
	if r.role != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "role", r.role, "")
	}
	if r.roleContains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "role__contains", r.roleContains, "")
	}
	if r.roleIcontains != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "role__icontains", r.roleIcontains, "")
	}
	if r.roleIn != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "role__in", r.roleIn, "csv")
	}
	if r.roleStartswith != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "role__startswith", r.roleStartswith, "")
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

type UsersRolesAPIUsersRolesReadRequest struct {
	ctx context.Context
	ApiService *UsersRolesAPIService
	authUsersUserRoleHref string
	fields *[]string
	excludeFields *[]string
}

// A list of fields to include in the response.
func (r UsersRolesAPIUsersRolesReadRequest) Fields(fields []string) UsersRolesAPIUsersRolesReadRequest {
	r.fields = &fields
	return r
}

// A list of fields to exclude from the response.
func (r UsersRolesAPIUsersRolesReadRequest) ExcludeFields(excludeFields []string) UsersRolesAPIUsersRolesReadRequest {
	r.excludeFields = &excludeFields
	return r
}

func (r UsersRolesAPIUsersRolesReadRequest) Execute() (*UserRoleResponse, *http.Response, error) {
	return r.ApiService.UsersRolesReadExecute(r)
}

/*
UsersRolesRead Inspect an user role

ViewSet for UserRole.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param authUsersUserRoleHref
 @return UsersRolesAPIUsersRolesReadRequest
*/
func (a *UsersRolesAPIService) UsersRolesRead(ctx context.Context, authUsersUserRoleHref string) UsersRolesAPIUsersRolesReadRequest {
	return UsersRolesAPIUsersRolesReadRequest{
		ApiService: a,
		ctx: ctx,
		authUsersUserRoleHref: authUsersUserRoleHref,
	}
}

// Execute executes the request
//  @return UserRoleResponse
func (a *UsersRolesAPIService) UsersRolesReadExecute(r UsersRolesAPIUsersRolesReadRequest) (*UserRoleResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *UserRoleResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "UsersRolesAPIService.UsersRolesRead")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "{auth_users_user_role_href}"
	localVarPath = strings.Replace(localVarPath, "{"+"auth_users_user_role_href"+"}", parameterValueToString(r.authUsersUserRoleHref, "authUsersUserRoleHref"), -1)  // NOTE: paths aren't escaped because Pulp uses hrefs as path parameters

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
