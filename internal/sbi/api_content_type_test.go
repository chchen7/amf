package sbi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// CVE-2026-41136: HTTPUEContextTransfer's Content-Type switch had no
// default case. An unsupported Content-Type left err nil, skipped the
// deserialization-error check, and invoked the processor with a
// completely uninitialized UeContextTransferRequest. Mirror the default
// case used by the analogous handlers (HTTPCreateUEContext).
func TestUEContextTransferRejectsUnsupportedContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := &Server{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/",
		bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "text/plain")

	s.HTTPUEContextTransfer(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unsupported content type, got %d", w.Code)
	}
}
