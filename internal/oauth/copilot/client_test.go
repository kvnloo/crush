package copilot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInitiatorTransportSetsHeader(t *testing.T) {
	tests := []struct {
		name string
		// body builds the request body; nil means a bodyless request.
		body func() (int, string)
		want string
	}{
		{
			name: "nil body defaults to user",
			body: func() (int, string) { return 0, "" }, // req.Body == nil
			want: "user",
		},
		{
			name: "NoBody defaults to user",
			body: func() (int, string) { return -1, "" }, // req.Body == http.NoBody
			want: "user",
		},
		{
			name: "user-only history stays user",
			body: func() (int, string) { return 1, `{"messages":[{"role":"user"}]}` },
			want: "user",
		},
		{
			name: "assistant history becomes agent",
			body: func() (int, string) { return 1, `{"messages":[{"role":"assistant"}]}` },
			want: "agent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.Header.Get("X-Initiator")
			}))
			defer srv.Close()

			kind, payload := tt.body()
			var req *http.Request
			var err error
			switch kind {
			case 0:
				req, err = http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
			case -1:
				req, err = http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, http.NoBody)
			default:
				req, err = http.NewRequestWithContext(t.Context(), http.MethodPost, srv.URL, strings.NewReader(payload))
			}
			require.NoError(t, err)

			client := &http.Client{Transport: &initiatorTransport{}}
			resp, err := client.Do(req)
			require.NoError(t, err)
			resp.Body.Close()

			require.Equal(t, tt.want, got)
		})
	}
}

// TestPollForTokenTiming verifies that PollForToken respects the
// server-mandated interval plus padding to prevent premature polls that
// trigger slow_down errors.
func TestPollForTokenTiming(t *testing.T) {
	// NOTE: Not parallel because we modify package-level test URLs.

	tests := []struct {
		name            string
		serverInterval  int
		serverResponses []string
		wantMinGaps     []time.Duration
		wantSuccess     bool
	}{
		{
			name:           "initial 5s interval with authorization_pending",
			serverInterval: 5,
			serverResponses: []string{
				`{"error": "authorization_pending"}`,
				`{"error": "authorization_pending"}`,
				`{"access_token": "gho_test"}`,
			},
			wantMinGaps: []time.Duration{
				7 * time.Second,
				7 * time.Second,
			},
			wantSuccess: true,
		},
		{
			name:           "slow_down increases interval by 5s per RFC 8628",
			serverInterval: 5,
			serverResponses: []string{
				`{"error": "slow_down"}`,
				`{"error": "slow_down"}`,
				`{"access_token": "gho_test"}`,
			},
			wantMinGaps: []time.Duration{
				7 * time.Second,
				12 * time.Second,
			},
			wantSuccess: true,
		},
		{
			name:           "multiple slow_downs compound the backoff",
			serverInterval: 5,
			serverResponses: []string{
				`{"error": "slow_down"}`,
				`{"error": "slow_down"}`,
				`{"error": "slow_down"}`,
				`{"access_token": "gho_test"}`,
			},
			wantMinGaps: []time.Duration{
				7 * time.Second,
				12 * time.Second,
				17 * time.Second,
			},
			wantSuccess: true,
		},
		{
			name:           "server interval below 5s gets bumped to 5s",
			serverInterval: 3,
			serverResponses: []string{
				`{"error": "authorization_pending"}`,
				`{"access_token": "gho_test"}`,
			},
			wantMinGaps: []time.Duration{
				7 * time.Second,
			},
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// NOTE: Not parallel - test modifies package state.

			var mu sync.Mutex
			var pollTimes []time.Time
			responseIdx := 0

			// Mock token endpoint that records poll times and returns responses.
			tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				now := time.Now()
				mu.Lock()
				pollTimes = append(pollTimes, now)
				idx := responseIdx
				responseIdx++
				mu.Unlock()

				w.Header().Set("Content-Type", "application/json")
				if idx < len(tt.serverResponses) {
					fmt.Fprint(w, tt.serverResponses[idx])
				} else {
					fmt.Fprint(w, `{"access_token": "gho_test"}`)
				}
			}))
			defer tokenSrv.Close()

			// Mock Copilot token endpoint.
			copilotSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check for Bearer token in Authorization header.
				auth := r.Header.Get("Authorization")
				if !strings.HasPrefix(auth, "Bearer gho_") {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				resp := map[string]interface{}{
					"token":      "test_copilot_token",
					"expires_at": time.Now().Add(24 * time.Hour).Unix(),
				}
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer copilotSrv.Close()

			// Replace global URLs for test.
			testAccessTokenURL = tokenSrv.URL
			testCopilotTokenURL = copilotSrv.URL
			t.Cleanup(func() {
				testAccessTokenURL = ""
				testCopilotTokenURL = ""
			})

			dc := &DeviceCode{
				DeviceCode:      "test_device_code",
				UserCode:        "TEST-1234",
				VerificationURI: "https://github.com/login/device",
				ExpiresIn:       300,
				Interval:        tt.serverInterval,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			defer cancel()

			_, err := PollForToken(ctx, dc)
			if tt.wantSuccess {
				require.NoError(t, err)
			}

			mu.Lock()
			defer mu.Unlock()

			// Verify the gaps between polls meet minimum requirements.
			require.GreaterOrEqual(t, len(pollTimes), 2, "should have at least 2 polls")
			for i := 1; i < len(pollTimes) && i-1 < len(tt.wantMinGaps); i++ {
				gap := pollTimes[i].Sub(pollTimes[i-1])
				require.GreaterOrEqual(
					t,
					gap,
					tt.wantMinGaps[i-1],
					"poll %d happened too early: gap=%v, wanted >= %v",
					i,
					gap,
					tt.wantMinGaps[i-1],
				)
			}
		})
	}
}

// TestPollForTokenContextCancellation verifies that PollForToken returns
// promptly when the context is cancelled.
func TestPollForTokenContextCancellation(t *testing.T) {
	t.Parallel()

	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"error": "authorization_pending"}`)
	}))
	defer tokenSrv.Close()

	testAccessTokenURL = tokenSrv.URL
	t.Cleanup(func() { testAccessTokenURL = "" })

	dc := &DeviceCode{
		DeviceCode:      "test_device_code",
		UserCode:        "TEST-1234",
		VerificationURI: "https://github.com/login/device",
		ExpiresIn:       300,
		Interval:        5,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	_, err := PollForToken(ctx, dc)
	elapsed := time.Since(start)

	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
	require.Less(t, elapsed, 2*time.Second, "should return promptly on cancellation")
}

// TestPollForTokenTimeout verifies that PollForToken returns an error when
// the device code expires.
func TestPollForTokenTimeout(t *testing.T) {
	// NOTE: Not parallel because we modify package-level test URLs.

	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"error": "authorization_pending"}`)
	}))
	defer tokenSrv.Close()

	testAccessTokenURL = tokenSrv.URL
	t.Cleanup(func() { testAccessTokenURL = "" })

	dc := &DeviceCode{
		DeviceCode:      "test_device_code",
		UserCode:        "TEST-1234",
		VerificationURI: "https://github.com/login/device",
		ExpiresIn:       3,
		Interval:        5,
	}

	ctx := context.Background()
	_, err := PollForToken(ctx, dc)

	require.Error(t, err)
	require.Contains(t, err.Error(), "authorization timed out")
}
