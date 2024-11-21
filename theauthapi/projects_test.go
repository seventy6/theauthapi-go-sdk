package theauthapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListProjects(t *testing.T) {
    // Create a mock server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/projects" && r.Method == http.MethodGet {
            projects := []Project{
                {ID: "1", Name: "Project 1"},
                {ID: "2", Name: "Project 2"},
            }
            json.NewEncoder(w).Encode(projects)
        } else {
            http.Error(w, "not found", http.StatusNotFound)
        }
    }))
    defer mockServer.Close()

    client := &Client{BaseURL: mockServer.URL, HTTPClient: mockServer.Client()}
    service := &ProjectsService{client: client}

    projects, err := service.List(context.Background())
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(projects) != 2 {
        t.Errorf("expected 2 projects, got %d", len(projects))
    }

    expectedNames := []string{"Project 1", "Project 2"}
    for i, project := range projects {
        if project.Name != expectedNames[i] {
            t.Errorf("expected project name %s, got %s", expectedNames[i], project.Name)
        }
    }
}