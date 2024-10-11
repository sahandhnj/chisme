package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Resource represents a database, cache, SSL certificate, etc.
type Resource struct {
	ID          int      `json:"id"`
	Shortname   string   `json:"shortname"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	ServerID    int      `json:"server_id"`
	IP          string   `json:"ip"`
	Status      string   `json:"status"`
	LogStream   string   `json:"log_stream"`
	Description string   `json:"description"`
	Metrics     *Metrics `json:"metrics,omitempty"`
}

// Metrics represents performance metrics of a resource.
type Metrics struct {
	CPUUsage         string `json:"cpu_usage"`
	MemoryUsage      string `json:"memory_usage"`
	ResponseTimeMs   int    `json:"response_time_ms"`
	QueriesPerSecond int    `json:"queries_per_second,omitempty"`
	DiskUsage        string `json:"disk_usage,omitempty"`
	CacheHitsPerSec  int    `json:"cache_hits_per_second,omitempty"`
	Uptime           string `json:"uptime,omitempty"`
	ErrorRate        string `json:"error_rate,omitempty"`
}

// Checkpoint represents internal health checks or services within an application.
type Checkpoint struct {
	ID          int      `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	CheckStatus string   `json:"check_status"`
	LogStream   string   `json:"log_stream"`
	Metrics     *Metrics `json:"metrics,omitempty"`
}

// Application struct represents an application and its resources.
type Application struct {
	ID          int          `json:"id"`
	Shortname   string       `json:"shortname"`
	Name        string       `json:"name"`
	ServerID    int          `json:"server_id"`
	Status      string       `json:"status"`
	IP          string       `json:"ip"`
	LogStream   string       `json:"log_stream"`
	Description string       `json:"description"`
	Metrics     Metrics      `json:"metrics"`
	Checkpoints []Checkpoint `json:"checkpoints"`
	Resources   []int        `json:"resources"`
}

// ApplicationShort
type ApplicationShort struct {
	ID        int    `json:"id"`
	Shortname string `json:"shortname"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

// Server struct represents a server and the applications running on it.
type Server struct {
	ID          int    `json:"id"`
	Shortname   string `json:"shortname"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	IP          string `json:"ip"`
	Description string `json:"description"`
	LogStream   string `json:"log_stream"`
}

// Data struct to hold everything loaded from the JSON file.
type Data struct {
	Servers      []Server      `json:"servers"`
	Applications []Application `json:"applications"`
	Resources    []Resource    `json:"resources"`
}

// Load data from JSON file
func loadData() (Data, error) {
	var data Data
	bytes, err := os.ReadFile("example/data.json")
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil
}

// Helper function to extract path variables from URL
func pathValue(r *http.Request, key string) string {
	// We can use URL.Path directly and split based on `/`
	pathSegments := strings.Split(r.URL.Path, "/")
	for i, segment := range pathSegments {
		if segment == key && i+1 < len(pathSegments) {
			return pathSegments[i+1]
		}
	}
	return ""
}

// Get all servers
func getServers(w http.ResponseWriter, r *http.Request) {
	data, err := loadData()
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data.Servers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getServerByID(w http.ResponseWriter, r *http.Request) {
	idStr := pathValue(r, "server")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	data, err := loadData()
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}

	var serverFound *Server
	for _, srv := range data.Servers {
		if srv.ID == id {
			serverFound = &srv
			break
		}
	}

	if serverFound == nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	var serverApplications []ApplicationShort
	for _, app := range data.Applications {
		if app.ServerID == id {
			simplifiedApp := ApplicationShort{
				ID:        app.ID,
				Shortname: app.Shortname,
				Name:      app.Name,
				Status:    app.Status,
			}
			serverApplications = append(serverApplications, simplifiedApp)
		}
	}

	var serverResources []Resource
	for _, resource := range data.Resources {
		if resource.ServerID == id {
			serverResources = append(serverResources, resource)
		}
	}

	response := struct {
		Server          *Server            `json:"server"`
		Applications    []ApplicationShort `json:"applications"`
		ServerResources []Resource         `json:"resources"`
	}{
		Server:          serverFound,
		Applications:    serverApplications,
		ServerResources: serverResources,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get a specific application on a server
func getApplicationByID(w http.ResponseWriter, r *http.Request) {
	serverIDStr := pathValue(r, "server")
	appIDStr := pathValue(r, "application")

	serverID, err := strconv.Atoi(serverIDStr)
	appID, err2 := strconv.Atoi(appIDStr)
	if err != nil || err2 != nil {
		http.Error(w, "Invalid server or application ID", http.StatusBadRequest)
		return
	}

	data, err := loadData()
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}

	var appFound *Application
	for _, app := range data.Applications {
		if app.ID == appID && app.ServerID == serverID {
			appFound = &app
			break
		}
	}

	if appFound == nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	var appResources []Resource
	for _, depID := range appFound.Resources {
		for _, dep := range data.Resources {
			if dep.ID == depID {
				appResources = append(appResources, dep)
			}
		}
	}

	response := struct {
		Application Application `json:"application"`
		Resources   []Resource  `json:"resources"`
	}{
		Application: *appFound,
		Resources:   appResources,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
