package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestExperienceYearCalculations(t *testing.T) {
	if got := startYearFromYears(15, 2026); got != 2011 {
		t.Errorf("startYearFromYears(15, 2026) = %d, want 2011", got)
	}
	if experienceStartYear != 2011 {
		t.Errorf("experienceStartYear = %d, want 2011", experienceStartYear)
	}

	tests := []struct {
		startYear   int
		currentYear int
		want        int
	}{
		{2011, 2026, 15},
		{2011, 2030, 19},
		{2011, 2011, 0},
		{2011, 2009, 0},
	}
	for _, tt := range tests {
		if got := yearsSince(tt.startYear, tt.currentYear); got != tt.want {
			t.Errorf("yearsSince(%d, %d) = %d, want %d", tt.startYear, tt.currentYear, got, tt.want)
		}
	}
	if got := experienceYears(2026); got != 15 {
		t.Errorf("experienceYears(2026) = %d, want 15", got)
	}
}

func TestCVSummaryUsesCalculatedYears(t *testing.T) {
	rr := perform(testRouter(t), http.MethodGet, "/api/cv")
	var cv CVData
	if err := json.Unmarshal(rr.Body.Bytes(), &cv); err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("%d+ years (since %d)", experienceYears(time.Now().Year()), experienceStartYear)
	if !strings.Contains(cv.Summary, want) {
		t.Errorf("summary = %q, want it to contain %q", cv.Summary, want)
	}
	if strings.Contains(cv.Summary, "15+ years across") {
		t.Error("summary still contains the old hardcoded years claim")
	}
}

func TestPageCalculatesYearsInBrowser(t *testing.T) {
	page := perform(testRouter(t), http.MethodGet, "/")
	body := page.Body.String()
	for _, snippet := range []string{
		"function startYearFromYears",
		"function yearsSince",
		"ORIGINAL_YEARS_CLAIM = 15",
		"ORIGINAL_CLAIM_YEAR = 2026",
		"{EXPERIENCE_YEARS}",
		"{EXPERIENCE_START_YEAR}",
	} {
		if !strings.Contains(body, snippet) {
			t.Errorf("page is missing %q", snippet)
		}
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
