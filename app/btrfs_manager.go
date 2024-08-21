package main

import (
    "fmt"
    "os"
    "os/exec"
    "log"
)

type BtrfsManager struct {
    mountPoint string
}

func NewBtrfsManager(mountPoint string) *BtrfsManager {
    return &BtrfsManager{mountPoint: mountPoint}
}

func (m *BtrfsManager) CreateVolume(name string) error {
    volumePath := fmt.Sprintf("%s/%s", m.mountPoint, name)
    if err := os.MkdirAll(volumePath, 0755); err != nil {
        return fmt.Errorf("failed to create directory %s: %w", volumePath, err)
    }
    cmd := exec.Command("btrfs", "subvolume", "create", volumePath)
    return cmd.Run()
}

func (m *BtrfsManager) RemoveVolume(name string) error {
    volumePath := fmt.Sprintf("%s/%s", m.mountPoint, name)
    cmd := exec.Command("btrfs", "subvolume", "delete", volumePath)
    return cmd.Run()
}

func (m *BtrfsManager) MountVolume(name string) error {
    volumePath := fmt.Sprintf("%s/%s", m.mountPoint, name)
    log.Printf("Mounting BTRFS volume: %s", volumePath)
    return nil
}

func (m *BtrfsManager) UnmountVolume(name string) error {
    volumePath := fmt.Sprintf("%s/%s", m.mountPoint, name)
    log.Printf("Unmounting BTRFS volume: %s", volumePath)
    return nil
}

func (m *BtrfsManager) GetMountPoint(name string) string {
    return fmt.Sprintf("%s/%s", m.mountPoint, name)
}
