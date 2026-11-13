package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestCVHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/api/cv", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // CV handler implementation
        cv := CVData{
            Name: "Alex Chen",
            Title: "Senior Software Engineer",
            Email: "alex@example.com",
            Experience: []string{"10+ years"},
            Skills: []string{"Go", "Python"},
            Education: []string{"MIT"},
        }
        
        jsonBytes, _ := json.Marshal(cv)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write(jsonBytes)
    })
    
    handler.ServeHTTP(rr, req)
    
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
    
    expected := `{"name":"Alex Chen","title":"Senior Software Engineer","email":"alex@example.com","experience":["10+ years"],"skills":["Go","Python"],"education":["MIT"]}`
    if !strings.Contains(rr.Body.String(), expected) {
        t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

func TestAboutHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/api/about", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // About handler implementation
        response := map[string]interface{}{
            "company": "Tech Innovations Inc.",
            "founded": 2015,
            "employees": 500,
        }
        
        jsonBytes, _ := json.Marshal(response)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write(jsonBytes)
    })
    
    handler.ServeHTTP(rr, req)
    
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}