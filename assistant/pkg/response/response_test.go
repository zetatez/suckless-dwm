package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestOk(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Ok(c, map[string]string{"key": "value"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Code != CodeSuccess {
		t.Errorf("expected code %d, got %d", CodeSuccess, resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("expected message 'success', got %q", resp.Message)
	}
}

func TestErr_InvalidParams(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Err(c, CodeInvalidParams, "invalid input")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Code != CodeInvalidParams {
		t.Errorf("expected code %d, got %d", CodeInvalidParams, resp.Code)
	}
}

func TestErr_Unauthorized(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Err(c, CodeUnauthorized, "not authenticated")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestErr_Forbidden(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Err(c, CodeForbidden, "access denied")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestErr_NotFound(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Err(c, CodeNotFound, "resource not found")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestErr_ServerError(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Err(c, CodeServerError, "internal error")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestErr_DatabaseError(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Err(c, CodeDatabaseError, "db connection failed")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestErr_ThirdPartyError(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Err(c, CodeThirdPartyErr, "external service error")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("expected status %d, got %d", http.StatusBadGateway, w.Code)
	}
}

func TestErr_UnknownCode(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Err(c, 99999, "unknown error")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d for unknown code, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestStatusFromCode(t *testing.T) {
	tests := []struct {
		code           int
		expectedStatus int
	}{
		{CodeSuccess, http.StatusOK},
		{CodeInvalidParams, http.StatusBadRequest},
		{CodeDatabaseError, http.StatusInternalServerError},
		{CodeThirdPartyErr, http.StatusBadGateway},
		{CodeNotFound, http.StatusNotFound},
		{CodeUnauthorized, http.StatusUnauthorized},
		{CodeForbidden, http.StatusForbidden},
		{CodeServerError, http.StatusInternalServerError},
		{99999, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		result := statusFromCode(tt.code)
		if result != tt.expectedStatus {
			t.Errorf("statusFromCode(%d) = %d, want %d", tt.code, result, tt.expectedStatus)
		}
	}
}

func TestOk_WithNilData(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		Ok(c, nil)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Code != CodeSuccess {
		t.Errorf("expected code %d, got %d", CodeSuccess, resp.Code)
	}
}
