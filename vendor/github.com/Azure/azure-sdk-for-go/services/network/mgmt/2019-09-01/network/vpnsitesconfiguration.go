package network

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// VpnSitesConfigurationClient is the network Client
type VpnSitesConfigurationClient struct {
	BaseClient
}

// NewVpnSitesConfigurationClient creates an instance of the VpnSitesConfigurationClient client.
func NewVpnSitesConfigurationClient(subscriptionID string) VpnSitesConfigurationClient {
	return NewVpnSitesConfigurationClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewVpnSitesConfigurationClientWithBaseURI creates an instance of the VpnSitesConfigurationClient client using a
// custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds,
// Azure stack).
func NewVpnSitesConfigurationClientWithBaseURI(baseURI string, subscriptionID string) VpnSitesConfigurationClient {
	return VpnSitesConfigurationClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Download gives the sas-url to download the configurations for vpn-sites in a resource group.
// Parameters:
// resourceGroupName - the resource group name.
// virtualWANName - the name of the VirtualWAN for which configuration of all vpn-sites is needed.
// request - parameters supplied to download vpn-sites configuration.
func (client VpnSitesConfigurationClient) Download(ctx context.Context, resourceGroupName string, virtualWANName string, request GetVpnSitesConfigurationRequest) (result VpnSitesConfigurationDownloadFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VpnSitesConfigurationClient.Download")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: request,
			Constraints: []validation.Constraint{{Target: "request.OutputBlobSasURL", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("network.VpnSitesConfigurationClient", "Download", err.Error())
	}

	req, err := client.DownloadPreparer(ctx, resourceGroupName, virtualWANName, request)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VpnSitesConfigurationClient", "Download", nil, "Failure preparing request")
		return
	}

	result, err = client.DownloadSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VpnSitesConfigurationClient", "Download", result.Response(), "Failure sending request")
		return
	}

	return
}

// DownloadPreparer prepares the Download request.
func (client VpnSitesConfigurationClient) DownloadPreparer(ctx context.Context, resourceGroupName string, virtualWANName string, request GetVpnSitesConfigurationRequest) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"virtualWANName":    autorest.Encode("path", virtualWANName),
	}

	const APIVersion = "2019-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{virtualWANName}/vpnConfiguration", pathParameters),
		autorest.WithJSON(request),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DownloadSender sends the Download request. The method will close the
// http.Response Body if it receives an error.
func (client VpnSitesConfigurationClient) DownloadSender(req *http.Request) (future VpnSitesConfigurationDownloadFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// DownloadResponder handles the response to the Download request. The method always
// closes the http.Response Body.
func (client VpnSitesConfigurationClient) DownloadResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
