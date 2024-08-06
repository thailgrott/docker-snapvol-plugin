package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "snapvol/app"
)

func main() {
    log.SetOutput(os.Stdout)
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    btrfsManager := app.NewBtrfsManager("/var/lib/docker-snap-volumes")
    pluginAPI := app.NewPluginAPI(btrfsManager)

    router := mux.NewRouter()
    router.HandleFunc("/VolumeDriver.Create", pluginAPI.CreateVolume).Methods("POST")
    router.HandleFunc("/VolumeDriver.Remove", pluginAPI.RemoveVolume).Methods("POST")
    router.HandleFunc("/VolumeDriver.Mount", pluginAPI.MountVolume).Methods("POST")
    router.HandleFunc("/VolumeDriver.Unmount", pluginAPI.UnmountVolume).Methods("POST")

    log.Println("Starting HTTP server on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Failed to start plugin API: %v", err)
    }
}
