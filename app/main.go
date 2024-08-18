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
    r.HandleFunc("/VolumeDriver.Create", CreateVolume).Methods("POST")
    r.HandleFunc("/VolumeDriver.Remove", RemoveVolume).Methods("POST")
    r.HandleFunc("/VolumeDriver.Mount", MountVolume).Methods("POST")
    r.HandleFunc("/VolumeDriver.Unmount", UnmountVolume).Methods("POST")

    // Start the HTTP server
    log.Println("Starting the SnapVol Docker Volume Plugin")
    err := http.ListenAndServe("/run/docker/plugins/snapvol.sock", r)
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
        os.Exit(1)
    }
}
