package main

import (
    "encoding/json"
    "net/http"
)

type PluginAPI struct {
    btrfsManager *BtrfsManager
}

func NewPluginAPI(btrfsManager *BtrfsManager) *PluginAPI {
    return &PluginAPI{
        btrfsManager: btrfsManager,
    }
}

func (api *PluginAPI) CreateVolume(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    name := req["Name"].(string)

    if err := api.btrfsManager.CreateVolume(name); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"Mountpoint": api.btrfsManager.GetMountPoint(name)})
}

func (api *PluginAPI) RemoveVolume(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    name := req["Name"].(string)

    if err := api.btrfsManager.RemoveVolume(name); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (api *PluginAPI) MountVolume(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    name := req["Name"].(string)

    if err := api.btrfsManager.MountVolume(name); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"Mountpoint": api.btrfsManager.GetMountPoint(name)})
}

func (api *PluginAPI) UnmountVolume(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    name := req["Name"].(string)

    if err := api.btrfsManager.UnmountVolume(name); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
