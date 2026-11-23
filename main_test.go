package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func testRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	return newRouter()
}

func perform(router *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestHealthHandler(t *testing.T) {
	rr := perform(testRouter(t), http.MethodGet, "/health")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "healthy" || body["service"] != "matteosalvatoreg" {
		t.Fatalf("unexpected health body: %v", body)
	}
}

func TestCVHandlerUsesRealData(t *testing.T) {
	rr := perform(testRouter(t), http.MethodGet, "/api/cv")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var cv CVData
	if err := json.Unmarshal(rr.Body.Bytes(), &cv); err != nil {
		t.Fatal(err)
	}
	if cv.Name != "Matthew Giancola" {
		t.Errorf("name = %q, want Matthew Giancola", cv.Name)
	}
	if cv.Email != "matthew.giancola42@gmail.com" {
		t.Errorf("email = %q", cv.Email)
	}
	if len(cv.Experience) < 5 || len(cv.Education) == 0 || len(cv.Skills) == 0 {
		t.Fatalf("incomplete CV: %d jobs, %d schools, %d skill groups", len(cv.Experience), len(cv.Education), len(cv.Skills))
	}
	if strings.Contains(rr.Body.String(), "Alex Chen") || strings.Contains(rr.Body.String(), "example.com") {
		t.Error("placeholder identity still present in /api/cv")
	}
}

func TestAboutHandlerUsesRealSiteInfo(t *testing.T) {
	rr := perform(testRouter(t), http.MethodGet, "/api/about")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if strings.Contains(rr.Body.String(), "Tech Innovations Inc.") {
		t.Error("placeholder company still present in /api/about")
	}
	if !strings.Contains(rr.Body.String(), "matteosalvatore.dev") {
		t.Error("real domain missing from /api/about")
	}
}

func TestStaticSiteIsServedByTheAPI(t *testing.T) {
	router := testRouter(t)

	page := perform(router, http.MethodGet, "/")
	if page.Code != http.StatusOK {
		t.Fatalf("page status = %d, want 200", page.Code)
	}
	if !strings.Contains(page.Body.String(), "Matthew Giancola") {
		t.Error("CV page did not contain the owner's name")
	}
	if !strings.HasPrefix(page.Header().Get("Content-Type"), "text/html") {
		t.Errorf("page content type = %q", page.Header().Get("Content-Type"))
	}

	image := perform(router, http.MethodGet, "/img.png")
	if image.Code != http.StatusOK {
		t.Fatalf("image status = %d, want 200", image.Code)
	}
	if image.Header().Get("Content-Type") != "image/png" {
		t.Errorf("image content type = %q", image.Header().Get("Content-Type"))
	}
	if !strings.HasPrefix(image.Body.String(), "\x89PNG") {
		t.Error("image response is not a PNG")
	}
}

func TestSecurityHeadersAndUnknownRoute(t *testing.T) {
	rr := perform(testRouter(t), http.MethodGet, "/does-not-exist")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("missing X-Content-Type-Options header")
	}
	if rr.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("missing X-Frame-Options header")
	}
}
