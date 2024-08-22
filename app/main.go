package main

import (
    "encoding/json"
    "fmt"
    "log"
    "log/syslog"
    "net"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/gorilla/mux"
)

const (
	driverName = "snapvol-plugin"
	runPath    = "/run/docker/plugins"
	socketName = "snapvol.sock"
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

func getPluginSocketPath() (string, error) {
	// Execute the Docker command to inspect the plugin
	cmd := exec.Command("docker", "plugin", "inspect", driverName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running docker plugin inspect: %w", err)
	}

	// Parse the JSON output to get the plugin ID
	var plugins []struct {
		ID string `json:"Id"`
	}
	if err := json.Unmarshal(output, &plugins); err != nil {
		return "", fmt.Errorf("error parsing JSON output: %w", err)
	}

	if len(plugins) == 0 {
		return "", fmt.Errorf("no plugin found with the name %s", driverName)
	}

	// Build the socket path
	pluginID := plugins[0].ID
	socketPath := filepath.Join(runPath, "plugins", pluginID, socketName)

	return socketPath, nil
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

	socketPath, err := getPluginSocketPath()
	if err != nil {
		log.Fatalf("Failed to get plugin socket path: %v", err)
	}

	fmt.Println("Socket path:", socketPath)

    // Remove any existing socket file
    if _, err := os.Stat(socketPath); err == nil {
        os.Remove(socketPath)
    }

    // Create a Unix socket listener
    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        log.Fatalf("Error creating Unix socket listener: %v", err)
    }

    // Adjust the file permissions of the socket
    if err := os.Chmod(socketPath, 0600); err != nil {
        log.Fatalf("Error setting permissions on the Unix socket: %v", err)
    }

    // Connect to the syslog daemon
    syslogger, err := syslog.New(syslog.LOG_NOTICE, "snapvol-plugin")
    if err != nil {
        log.Fatalf("Failed to connect to syslog: %v", err)
    }
    log.SetOutput(syslogger)
    
    // Start the HTTP server on the Unix socket
    log.Println("Starting the SnapVol Docker Volume Plugin")
    err = http.Serve(listener, r)
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
        os.Exit(1)
    }
}
