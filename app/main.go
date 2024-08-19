package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
)

// PluginActivateResponse represents the response to the /Plugin.Activate request
type PluginActivateResponse struct {
    Implements []string `json:"Implements"`
}

// VolumeDriverCapabilitiesResponse represents the response to the /VolumeDriver.Capabilities request
type VolumeDriverCapabilitiesResponse struct {
    Capabilities struct {
        Scope string `json:"Scope"`
    } `json:"Capabilities"`
}

func main() {
    btrfsManager := NewBtrfsManager("/var/lib/docker-snap-volumes") // Initialize your BtrfsManager
    pluginAPI := NewPluginAPI(btrfsManager) // Create an instance of PluginAPI

    r := mux.NewRouter()

    // Handle the /Plugin.Activate route
    r.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, r *http.Request) {
        response := PluginActivateResponse{
            Implements: []string{"VolumeDriver"},
        }
        json.NewEncoder(w).Encode(response)
    }).Methods("POST")

    // Handle the /VolumeDriver.Capabilities route
    r.HandleFunc("/VolumeDriver.Capabilities", func(w http.ResponseWriter, r *http.Request) {
        response := VolumeDriverCapabilitiesResponse{}
        response.Capabilities.Scope = "local"
        json.NewEncoder(w).Encode(response)
    }).Methods("POST")

    // Handle routes for the standard Docker volume plugin functionality
    r.HandleFunc("/VolumeDriver.Create", pluginAPI.CreateVolume).Methods("POST")
    r.HandleFunc("/VolumeDriver.Remove", pluginAPI.RemoveVolume).Methods("POST")
    r.HandleFunc("/VolumeDriver.Mount", pluginAPI.MountVolume).Methods("POST")
    r.HandleFunc("/VolumeDriver.Unmount", pluginAPI.UnmountVolume).Methods("POST")

    // Define the Unix socket path
    socketPath := "/run/docker/plugins/snapvol.sock"

    // Remove any existing socket file
    if _, err := os.Stat(socketPath); err == nil {
        os.Remove(socketPath)
    }

    // Create a Unix socket listener
    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        log.Fatalf("Error creating Unix socket listener: %v", err)
    }

    // Start the HTTP server on the Unix socket
    log.Println("Starting the SnapVol Docker Volume Plugin")
    err = http.Serve(listener, r)
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
        os.Exit(1)
    }
}
