package sbi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// GHSA-xw5p-5pgh-4xq5: the AMF communication handlers returned the raw
// deserialization error string in the ProblemDetails Detail field, exposing
// fully-qualified internal Go struct names to the client. The handlers must
// keep the detailed error in the server log and return a generic
// client-facing message instead.
func TestDeserializeErrorsDoNotLeakInternals(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := &Server{}

	// Truncated JSON guarantees the deserialization failure path.
	malformed := []byte(`{"notificationUri":`)

	handlers := []struct {
		name    string
		handler func(*gin.Context)
	}{
		{"HTTPAMFStatusChangeSubscribeModify", s.HTTPAMFStatusChangeSubscribeModify},
		{"HTTPCreateUEContext", s.HTTPCreateUEContext},
		{"HTTPEBIAssignment", s.HTTPEBIAssignment},
		{"HTTPRegistrationStatusUpdate", s.HTTPRegistrationStatusUpdate},
		{"HTTPReleaseUEContext", s.HTTPReleaseUEContext},
		{"HTTPUEContextTransfer", s.HTTPUEContextTransfer},
		{"HTTPN1N2MessageTransfer", s.HTTPN1N2MessageTransfer},
		{"HTTPN1N2MessageSubscribe", s.HTTPN1N2MessageSubscribe},
		{"HTTPAMFStatusChangeSubscribe", s.HTTPAMFStatusChangeSubscribe},
	}

	for _, tc := range handlers {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/",
				bytes.NewReader(malformed))
			c.Request.Header.Set("Content-Type", "application/json")

			tc.handler(c)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400 for malformed body, got %d", w.Code)
			}
			body := w.Body.String()
			// The client-facing detail must be the generic message, not the
			// raw deserialization error (which embeds internal details).
			if !bytes.Contains([]byte(body),
				[]byte("Failed to deserialize request body")) {
				t.Fatalf("expected generic detail message, got: %s", body)
			}
		})
	}
}
