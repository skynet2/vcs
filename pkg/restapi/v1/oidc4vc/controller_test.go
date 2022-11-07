/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc4vc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/trustbloc/vcs/pkg/restapi/v1/issuer"
	"github.com/trustbloc/vcs/pkg/restapi/v1/oidc4vc"
	"github.com/trustbloc/vcs/pkg/storage/mongodb/oidc4vcstatestore"
)

//nolint:lll
func TestController_OidcPushedAuthorizationRequest(t *testing.T) {
	var (
		mockOAuthProvider     = NewMockOAuth2Provider(gomock.NewController(t))
		mockInteractionClient = NewMockIssuerInteractionClient(gomock.NewController(t))
		q                     url.Values
	)

	tests := []struct {
		name  string
		setup func()
		check func(t *testing.T, rec *httptest.ResponseRecorder, err error)
	}{
		{
			name: "success",
			setup: func() {
				mockOAuthProvider.EXPECT().NewPushedAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{}, nil)
				mockOAuthProvider.EXPECT().NewPushedAuthorizeResponse(gomock.Any(), gomock.Any(), gomock.Any()).Return(&fosite.PushedAuthorizeResponse{}, nil)
				mockOAuthProvider.EXPECT().WritePushedAuthorizeResponse(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

				mockInteractionClient.EXPECT().PushAuthorizationDetails(gomock.Any(), gomock.Any()).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(nil)),
					}, nil)

				q = url.Values{}
				q.Add("op_state", "opState")
				q.Add("authorization_details", `{"type":"openid_credential","credential_type":"UniversityDegreeCredential","format":"ldp_vc"}`)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "invalid pushed authorize request",
			setup: func() {
				mockOAuthProvider.EXPECT().NewPushedAuthorizeRequest(gomock.Any(), gomock.Any()).Return(nil, errors.New("par error"))
				q = url.Values{}
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "par error")
			},
		},
		{
			name: "fail to unmarshal authorization details",
			setup: func() {
				mockOAuthProvider.EXPECT().NewPushedAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{}, nil)

				q = url.Values{}
				q.Add("op_state", "opState")
				q.Add("authorization_details", "invalid")
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "invalid-value[authorization_details]")
			},
		},
		{
			name: "fail to validate authorization details",
			setup: func() {
				mockOAuthProvider.EXPECT().NewPushedAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{}, nil)

				q = url.Values{}
				q.Add("op_state", "opState")
				q.Add("authorization_details", `{"type":"invalid","credential_type":"UniversityDegreeCredential","format":"ldp_vc"}`)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "type should be 'openid_credential'")
			},
		},
		{
			name: "fail to push authorization details",
			setup: func() {
				mockOAuthProvider.EXPECT().NewPushedAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{}, nil)
				mockInteractionClient.EXPECT().PushAuthorizationDetails(gomock.Any(), gomock.Any()).Return(nil, errors.New("push authorization details error"))

				q = url.Values{}
				q.Add("op_state", "opState")
				q.Add("authorization_details", `{"type":"openid_credential","credential_type":"UniversityDegreeCredential","format":"ldp_vc"}`)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "push authorization details error")
			},
		},
		{
			name: "invalid status code for push authorization details",
			setup: func() {
				mockOAuthProvider.EXPECT().NewPushedAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{}, nil)

				mockInteractionClient.EXPECT().PushAuthorizationDetails(gomock.Any(), gomock.Any()).Return(
					&http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBuffer(nil)),
					}, nil)

				q = url.Values{}
				q.Add("op_state", "opState")
				q.Add("authorization_details", `{"type":"openid_credential","credential_type":"UniversityDegreeCredential","format":"ldp_vc"}`)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "push authorization details: status code")
			},
		},
		{
			name: "fail to create new pushed authorize response",
			setup: func() {
				mockOAuthProvider.EXPECT().NewPushedAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{}, nil)
				mockOAuthProvider.EXPECT().NewPushedAuthorizeResponse(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("new pushed authorize response error"))

				mockInteractionClient.EXPECT().PushAuthorizationDetails(gomock.Any(), gomock.Any()).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(nil)),
					}, nil)

				q = url.Values{}
				q.Add("op_state", "opState")
				q.Add("authorization_details", `{"type":"openid_credential","credential_type":"UniversityDegreeCredential","format":"ldp_vc"}`)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "new pushed authorize response error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			controller := oidc4vc.NewController(&oidc4vc.Config{
				OAuth2Provider:          mockOAuthProvider,
				IssuerInteractionClient: mockInteractionClient,
				IssuerVCSPublicHost:     "https://issuer.example.com",
			})

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(q.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()

			err := controller.OidcPushedAuthorizationRequest(echo.New().NewContext(req, rec))
			tt.check(t, rec, err)
		})
	}
}

//nolint:lll
func TestController_OidcAuthorize(t *testing.T) {
	var (
		mockOAuthProvider     = NewMockOAuth2Provider(gomock.NewController(t))
		mockStateStore        = NewMockStateStore(gomock.NewController(t))
		mockInteractionClient = NewMockIssuerInteractionClient(gomock.NewController(t))
		params                oidc4vc.OidcAuthorizeParams
	)

	tests := []struct {
		name  string
		setup func()
		check func(t *testing.T, rec *httptest.ResponseRecorder, err error)
	}{
		{
			name: "success",
			setup: func() {
				params = oidc4vc.OidcAuthorizeParams{
					ResponseType: "code",
					OpState:      "opState",
				}

				scope := []string{"openid", "profile"}

				mockOAuthProvider.EXPECT().NewAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{
					Request: fosite.Request{RequestedScope: scope},
				}, nil)

				mockOAuthProvider.EXPECT().NewAuthorizeResponse(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						ar fosite.AuthorizeRequester,
						session fosite.Session,
					) (fosite.AuthorizeResponder, error) {
						assert.Equal(t, params.OpState, ar.(*fosite.AuthorizeRequest).State)

						return &fosite.AuthorizeResponse{}, nil
					})

				b, err := json.Marshal(&issuer.PrepareClaimDataAuthorizationResponse{
					AuthorizationRequest: issuer.OAuthParameters{},
				})
				require.NoError(t, err)

				mockInteractionClient.EXPECT().PrepareAuthorizationRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						req issuer.PrepareAuthorizationRequestJSONRequestBody,
						reqEditors ...issuer.RequestEditorFn,
					) (*http.Response, error) {
						assert.Equal(t, params.ResponseType, req.ResponseType)
						assert.Equal(t, params.OpState, req.OpState)
						assert.Equal(t, lo.ToPtr(scope), req.Scope)

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBuffer(b)),
						}, nil
					})

				mockStateStore.EXPECT().SaveAuthorizeState(gomock.Any(), params.OpState, gomock.Any()).Return(nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.NoError(t, err)
				require.Equal(t, http.StatusSeeOther, rec.Code)
				require.NotEmpty(t, rec.Header().Get("Location"))
			},
		},
		{
			name: "invalid authorize request",
			setup: func() {
				params = oidc4vc.OidcAuthorizeParams{
					ResponseType: "code",
					OpState:      "opState",
				}

				mockOAuthProvider.EXPECT().NewAuthorizeRequest(gomock.Any(), gomock.Any()).Return(nil, errors.New("authorize error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "authorize error")
			},
		},
		{
			name: "prepare claim data authorization",
			setup: func() {
				params = oidc4vc.OidcAuthorizeParams{
					ResponseType: "code",
					OpState:      "opState",
				}

				scope := []string{"openid", "profile"}

				mockOAuthProvider.EXPECT().NewAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{
					Request: fosite.Request{RequestedScope: scope},
				}, nil)

				mockInteractionClient.EXPECT().PrepareAuthorizationRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						req issuer.PrepareAuthorizationRequestJSONRequestBody,
						reqEditors ...issuer.RequestEditorFn,
					) (*http.Response, error) {
						assert.Equal(t, params.ResponseType, req.ResponseType)
						assert.Equal(t, params.OpState, req.OpState)
						assert.Equal(t, lo.ToPtr(scope), req.Scope)

						return nil, errors.New("prepare claim data authorization error")
					})
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "prepare claim data authorization")
			},
		},
		{
			name: "invalid status code for prepare claim data authorization",
			setup: func() {
				params = oidc4vc.OidcAuthorizeParams{
					ResponseType: "code",
					OpState:      "opState",
				}

				scope := []string{"openid", "profile"}

				mockOAuthProvider.EXPECT().NewAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{
					Request: fosite.Request{RequestedScope: scope},
				}, nil)

				mockInteractionClient.EXPECT().PrepareAuthorizationRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						req issuer.PrepareAuthorizationRequestJSONRequestBody,
						reqEditors ...issuer.RequestEditorFn,
					) (*http.Response, error) {
						assert.Equal(t, params.ResponseType, req.ResponseType)
						assert.Equal(t, params.OpState, req.OpState)
						assert.Equal(t, lo.ToPtr(scope), req.Scope)

						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBuffer(nil)),
						}, nil
					})
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "prepare claim data authorization: status code")
			},
		},
		{
			name: "fail to create authorize response",
			setup: func() {
				params = oidc4vc.OidcAuthorizeParams{
					ResponseType: "code",
					OpState:      "opState",
				}

				scope := []string{"openid", "profile"}

				mockOAuthProvider.EXPECT().NewAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{
					Request: fosite.Request{RequestedScope: scope},
				}, nil)

				mockOAuthProvider.EXPECT().NewAuthorizeResponse(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						ar fosite.AuthorizeRequester,
						session fosite.Session,
					) (fosite.AuthorizeResponder, error) {
						assert.Equal(t, params.OpState, ar.(*fosite.AuthorizeRequest).State)

						return nil, errors.New("create authorize response error")
					})

				b, err := json.Marshal(&issuer.PrepareClaimDataAuthorizationResponse{
					AuthorizationRequest: issuer.OAuthParameters{},
				})
				require.NoError(t, err)

				mockInteractionClient.EXPECT().PrepareAuthorizationRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						req issuer.PrepareAuthorizationRequestJSONRequestBody,
						reqEditors ...issuer.RequestEditorFn,
					) (*http.Response, error) {
						assert.Equal(t, params.ResponseType, req.ResponseType)
						assert.Equal(t, params.OpState, req.OpState)
						assert.Equal(t, lo.ToPtr(scope), req.Scope)

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBuffer(b)),
						}, nil
					})
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "create authorize response error")
			},
		},
		{
			name: "fail to save authorize state",
			setup: func() {
				params = oidc4vc.OidcAuthorizeParams{
					ResponseType: "code",
					OpState:      "opState",
				}

				scope := []string{"openid", "profile"}

				mockOAuthProvider.EXPECT().NewAuthorizeRequest(gomock.Any(), gomock.Any()).Return(&fosite.AuthorizeRequest{
					Request: fosite.Request{RequestedScope: scope},
				}, nil)

				mockOAuthProvider.EXPECT().NewAuthorizeResponse(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						ar fosite.AuthorizeRequester,
						session fosite.Session,
					) (fosite.AuthorizeResponder, error) {
						assert.Equal(t, params.OpState, ar.(*fosite.AuthorizeRequest).State)

						return &fosite.AuthorizeResponse{}, nil
					})

				b, err := json.Marshal(&issuer.PrepareClaimDataAuthorizationResponse{
					AuthorizationRequest: issuer.OAuthParameters{},
				})
				require.NoError(t, err)

				mockInteractionClient.EXPECT().PrepareAuthorizationRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						req issuer.PrepareAuthorizationRequestJSONRequestBody,
						reqEditors ...issuer.RequestEditorFn,
					) (*http.Response, error) {
						assert.Equal(t, params.ResponseType, req.ResponseType)
						assert.Equal(t, params.OpState, req.OpState)
						assert.Equal(t, lo.ToPtr(scope), req.Scope)

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBuffer(b)),
						}, nil
					})

				mockStateStore.EXPECT().SaveAuthorizeState(gomock.Any(), params.OpState, gomock.Any()).Return(
					errors.New("save state error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "save authorize state")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			controller := oidc4vc.NewController(&oidc4vc.Config{
				OAuth2Provider:          mockOAuthProvider,
				StateStore:              mockStateStore,
				IssuerInteractionClient: mockInteractionClient,
				IssuerVCSPublicHost:     "https://issuer.example.com",
			})

			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()

			err := controller.OidcAuthorize(echo.New().NewContext(req, rec), params)
			tt.check(t, rec, err)
		})
	}
}

//nolint:lll
func TestController_OidcRedirect(t *testing.T) {
	var (
		mockOAuthProvider     = NewMockOAuth2Provider(gomock.NewController(t))
		mockStateStore        = NewMockStateStore(gomock.NewController(t))
		mockInteractionClient = NewMockIssuerInteractionClient(gomock.NewController(t))
		params                oidc4vc.OidcRedirectParams
	)

	tests := []struct {
		name  string
		setup func()
		check func(t *testing.T, rec *httptest.ResponseRecorder, err error)
	}{
		{
			name: "success",
			setup: func() {
				params = oidc4vc.OidcRedirectParams{
					Code:  "code",
					State: "state",
				}

				redirectURI := &url.URL{Scheme: "https", Host: "example.com", Path: "redirect"}

				mockStateStore.EXPECT().GetAuthorizeState(gomock.Any(), params.State).Return(&oidc4vcstatestore.AuthorizeState{
					RedirectURI: redirectURI,
				}, nil)
				mockInteractionClient.EXPECT().StoreAuthorizationCodeRequest(
					gomock.Any(),
					issuer.StoreAuthorizationCodeRequest{
						Code:    params.Code,
						OpState: params.State,
					}).Return(&http.Response{Body: io.NopCloser(bytes.NewBuffer(nil))}, nil)

				mockOAuthProvider.EXPECT().WriteAuthorizeResponse(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(
						ctx context.Context,
						rw http.ResponseWriter,
						ar fosite.AuthorizeRequester,
						responder fosite.AuthorizeResponder,
					) {
						assert.Equal(t, redirectURI, ar.GetRedirectURI())
						assert.Equal(t, params.State, ar.GetState())
					})
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "fail to store code",
			setup: func() {
				params = oidc4vc.OidcRedirectParams{
					Code:  "code",
					State: "state",
				}

				redirectURI := &url.URL{Scheme: "https", Host: "example.com", Path: "redirect"}

				mockStateStore.EXPECT().GetAuthorizeState(gomock.Any(), params.State).Return(&oidc4vcstatestore.AuthorizeState{
					RedirectURI: redirectURI,
				}, nil)
				mockInteractionClient.EXPECT().StoreAuthorizationCodeRequest(
					gomock.Any(),
					issuer.StoreAuthorizationCodeRequest{
						Code:    params.Code,
						OpState: params.State,
					}).Return(nil, errors.New("random error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "random error")
			},
		},
		{
			name: "fail to get authorize state",
			setup: func() {
				params = oidc4vc.OidcRedirectParams{
					Code:  "code",
					State: "state",
				}

				mockStateStore.EXPECT().GetAuthorizeState(gomock.Any(), params.State).Return(nil, errors.New("get state error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "get state error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			controller := oidc4vc.NewController(&oidc4vc.Config{
				OAuth2Provider:          mockOAuthProvider,
				StateStore:              mockStateStore,
				IssuerInteractionClient: mockInteractionClient,
				IssuerVCSPublicHost:     "https://issuer.example.com",
			})

			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()

			err := controller.OidcRedirect(echo.New().NewContext(req, rec), params)
			tt.check(t, rec, err)
		})
	}
}

func TestController_OidcToken(t *testing.T) {
	var (
		mockOAuthProvider     = NewMockOAuth2Provider(gomock.NewController(t))
		mockInteractionClient = NewMockIssuerInteractionClient(gomock.NewController(t))
	)

	tests := []struct {
		name  string
		setup func()
		check func(t *testing.T, rec *httptest.ResponseRecorder, err error)
	}{
		{
			name: "success",
			setup: func() {
				opState := uuid.NewString()
				mockOAuthProvider.EXPECT().NewAccessRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&fosite.AccessRequest{
						Request: fosite.Request{
							Session: &fosite.DefaultSession{
								Extra: map[string]interface{}{
									"opState": opState,
								},
							},
						},
					}, nil)

				mockInteractionClient.EXPECT().ExchangeAuthorizationCodeRequest(gomock.Any(),
					issuer.ExchangeAuthorizationCodeRequestJSONRequestBody{
						OpState: opState,
					}).Return(&http.Response{Body: io.NopCloser(bytes.NewBuffer(nil))}, nil)

				mockOAuthProvider.EXPECT().NewAccessResponse(gomock.Any(), gomock.Any()).Return(
					&fosite.AccessResponse{}, nil)

				mockOAuthProvider.EXPECT().WriteAccessResponse(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "fail to create new access request",
			setup: func() {
				mockOAuthProvider.EXPECT().NewAccessRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, errors.New("new access request error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "new access request error")
			},
		},
		{
			name: "fail to create new access response",
			setup: func() {
				mockOAuthProvider.EXPECT().NewAccessRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&fosite.AccessRequest{
						Request: fosite.Request{
							Session: &fosite.DefaultSession{
								Extra: map[string]interface{}{
									"opState": "1234",
								},
							},
						},
					}, nil)

				mockInteractionClient.EXPECT().ExchangeAuthorizationCodeRequest(gomock.Any(), gomock.Any()).
					Return(&http.Response{Body: io.NopCloser(bytes.NewBuffer(nil))}, nil)

				mockOAuthProvider.EXPECT().NewAccessResponse(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("new access response error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "new access response error")
			},
		},
		{
			name: "fail to exchange token",
			setup: func() {
				mockOAuthProvider.EXPECT().NewAccessRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&fosite.AccessRequest{
						Request: fosite.Request{
							Session: &fosite.DefaultSession{
								Extra: map[string]interface{}{
									"opState": "1234",
								},
							},
						},
					}, nil)

				mockInteractionClient.EXPECT().ExchangeAuthorizationCodeRequest(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("can not exchange token"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				require.ErrorContains(t, err, "can not exchange token")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			controller := oidc4vc.NewController(&oidc4vc.Config{
				OAuth2Provider:          mockOAuthProvider,
				IssuerInteractionClient: mockInteractionClient,
			})

			req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()

			err := controller.OidcToken(echo.New().NewContext(req, rec))
			tt.check(t, rec, err)
		})
	}
}

//nolint:lll
func TestController_OidcPreAuthorize(t *testing.T) {
	var (
		mockOAuthProvider     = NewMockOAuth2Provider(gomock.NewController(t))
		mockInteractionClient = NewMockIssuerInteractionClient(gomock.NewController(t))
		oauthClient           = NewMockOAuth2Client(gomock.NewController(t))
		preAuthorizeClient    = NewMockHttpClient(gomock.NewController(t))
	)

	tests := []struct {
		name  string
		body  io.Reader
		setup func()
		check func(t *testing.T, rec *httptest.ResponseRecorder, err error)
	}{
		{
			name: "success",
			body: strings.NewReader(url.Values{
				"grant_type":          {"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
				"pre-authorized_code": {"123456"},
				"user_pin":            {"5678"},
			}.Encode()),
			setup: func() {
				mockInteractionClient.EXPECT().ValidatePreAuthorizedCodeRequest(gomock.Any(), gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"scopes" : ["a","b"], "op_state" : "opp123"}`)),
					}, nil)
				oauthClient.EXPECT().GeneratePKCE().Return("verifier", "challenge", "S256", nil)
				finalURL := "https://127.0.0.1/authorize"

				cfg := oauth2.Config{
					ClientID:     "pre-auth-client",
					ClientSecret: "foobar",
					RedirectURL:  "/oidc/token",
					Scopes:       []string{"a", "b"},
					Endpoint: oauth2.Endpoint{
						AuthURL:   "/oidc/authorize",
						TokenURL:  "/oidc/token",
						AuthStyle: oauth2.AuthStyleInParams,
					},
				}

				oauthClient.EXPECT().AuthCodeURL(gomock.Any(), cfg, "opp123", gomock.Any()).Return(finalURL)

				preAuthorizeClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "/authorize", req.URL.Path)
					return &http.Response{
						StatusCode: http.StatusSeeOther,
						Header: map[string][]string{
							"Location": {"https://localhost/redirect?code=my-secure-code"},
						},
					}, nil
				})

				oauthClient.EXPECT().ExchangeWithCustomClient(gomock.Any(), cfg, "my-secure-code",
					gomock.Any(), gomock.Any()).Return(&oauth2.Token{
					AccessToken: "123456",
					Expiry:      time.Now().UTC().Add(10 * time.Minute),
				}, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				var resp oidc4vc.AccessTokenResponse
				assert.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
				assert.Equal(t, "123456", resp.AccessToken)
				assert.NotEmpty(t, *resp.ExpiresIn)
			},
		},
		{
			name:  "invalid grant type",
			setup: func() {},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.ErrorContains(t, err, "unexpected grant type. expected urn:ietf:params:oauth:grant-type:pre-authorized_code")
			},
		},
		{
			name: "name invalid pre-auth code",
			body: strings.NewReader(url.Values{
				"grant_type":          {"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
				"pre-authorized_code": {"123456"},
				"user_pin":            {"5678"},
			}.Encode()),
			setup: func() {
				mockInteractionClient.EXPECT().ValidatePreAuthorizedCodeRequest(gomock.Any(),
					issuer.ValidatePreAuthorizedCodeRequestJSONRequestBody{
						PreAuthorizedCode: "123456",
						UserPin:           lo.ToPtr("5678"),
					}).Return(nil, errors.New("invalid pin"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.ErrorContains(t, err, "invalid pin")
			},
		},
		{
			name: "invalid response from validator",
			body: strings.NewReader(url.Values{
				"grant_type":          {"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
				"pre-authorized_code": {"321"},
			}.Encode()),
			setup: func() {
				mockInteractionClient.EXPECT().ValidatePreAuthorizedCodeRequest(gomock.Any(),
					issuer.ValidatePreAuthorizedCodeRequestJSONRequestBody{
						PreAuthorizedCode: "321",
						UserPin:           lo.ToPtr(""),
					}).Return(&http.Response{
					Body:       io.NopCloser(strings.NewReader("{")),
					StatusCode: http.StatusOK,
				}, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.ErrorContains(t, err, "unexpected EOF")
			},
		},
		{
			name: "fail pkce",
			body: strings.NewReader(url.Values{
				"grant_type":          {"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
				"pre-authorized_code": {"321"},
			}.Encode()),
			setup: func() {
				mockInteractionClient.EXPECT().ValidatePreAuthorizedCodeRequest(gomock.Any(), gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"scopes" : ["a","b"]}`)),
					}, nil)
				oauthClient.EXPECT().GeneratePKCE().Return("", "", "", errors.New("unexpected pkce error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.ErrorContains(t, err, "unexpected pkce error")
			},
		},
		{
			name: "fail http pre-auth",
			body: strings.NewReader(url.Values{
				"grant_type":          {"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
				"pre-authorized_code": {"321"},
			}.Encode()),
			setup: func() {
				mockInteractionClient.EXPECT().ValidatePreAuthorizedCodeRequest(gomock.Any(), gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"scopes" : ["a","b"]}`)),
					}, nil)
				oauthClient.EXPECT().GeneratePKCE().Return("a", "b", "c", nil)
				oauthClient.EXPECT().AuthCodeURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://localhost")

				preAuthorizeClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("http error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.ErrorContains(t, err, "http error")
			},
		},
		{
			name: "fail http status code pre-auth",
			body: strings.NewReader(url.Values{
				"grant_type":          {"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
				"pre-authorized_code": {"321"},
			}.Encode()),
			setup: func() {
				mockInteractionClient.EXPECT().ValidatePreAuthorizedCodeRequest(gomock.Any(), gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"scopes" : ["a","b"]}`)),
					}, nil)
				oauthClient.EXPECT().GeneratePKCE().Return("a", "b", "c", nil)
				oauthClient.EXPECT().AuthCodeURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://localhost")

				preAuthorizeClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
				}, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.ErrorContains(t, err, "unexpected status code 200")
			},
		},
		{
			name: "fail location code pre-auth",
			body: strings.NewReader(url.Values{
				"grant_type":          {"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
				"pre-authorized_code": {"321"},
			}.Encode()),
			setup: func() {
				mockInteractionClient.EXPECT().ValidatePreAuthorizedCodeRequest(gomock.Any(), gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"scopes" : ["a","b"]}`)),
					}, nil)
				oauthClient.EXPECT().GeneratePKCE().Return("a", "b", "c", nil)
				oauthClient.EXPECT().AuthCodeURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://localhost")

				preAuthorizeClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusSeeOther,
					Header: map[string][]string{
						"Location": {"postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require"},
					},
				}, nil)
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.ErrorContains(t, err, "net/url: invalid userinfo")
			},
		},
		{
			name: "fail exchange code pre-auth",
			body: strings.NewReader(url.Values{
				"grant_type":          {"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
				"pre-authorized_code": {"321"},
			}.Encode()),
			setup: func() {
				mockInteractionClient.EXPECT().ValidatePreAuthorizedCodeRequest(gomock.Any(), gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"scopes" : ["a","b"]}`)),
					}, nil)
				oauthClient.EXPECT().GeneratePKCE().Return("a", "b", "c", nil)
				oauthClient.EXPECT().AuthCodeURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://localhost")

				preAuthorizeClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusSeeOther,
					Header: map[string][]string{
						"Location": {"https://localhost"},
					},
				}, nil)

				oauthClient.EXPECT().ExchangeWithCustomClient(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("exchange error"))
			},
			check: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.ErrorContains(t, err, "exchange error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			controller := oidc4vc.NewController(&oidc4vc.Config{
				OAuth2Provider:          mockOAuthProvider,
				IssuerInteractionClient: mockInteractionClient,
				OAuth2Client:            oauthClient,
				PreAuthorizeClient:      preAuthorizeClient,
			})

			req := httptest.NewRequest(http.MethodPost, "/", tt.body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()

			err := controller.OidcPreAuthorizedCode(echo.New().NewContext(req, rec))
			tt.check(t, rec, err)
		})
	}
}
