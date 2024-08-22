package main

import (
    "fmt"
    "os"
    "os/exec"
    "log"
    "path/filepath"
    "errors"
)

// BtrfsManager manages BTRFS volumes
type BtrfsManager struct {
    mountPoint      string
    volumeStorePath string
}

// NewBtrfsManager initializes a new BtrfsManager
func NewBtrfsManager(mountPoint string) *BtrfsManager {
    volumeStorePath := filepath.Join(mountPoint, "volume_store")
    return &BtrfsManager{
        mountPoint:      mountPoint,
        volumeStorePath: volumeStorePath,
    }
}

// isBtrfsVolume checks if the provided path is a BTRFS volume
func (m *BtrfsManager) isBtrfsVolume(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return fmt.Errorf("path does not exist: %s", path)
    }
    cmd := exec.Command("btrfs", "subvolume", "show", path)
    if err := cmd.Run(); err != nil {
        return errors.New("not a BTRFS subvolume")
    }
    return nil
}

// CreateVolume creates a new BTRFS subvolume
func (m *BtrfsManager) CreateVolume(name string) error {
    volumePath := filepath.Join(m.volumeStorePath, name)
    
    if err := m.isBtrfsVolume(m.mountPoint); err != nil {
        log.Printf("Mount point check failed: %v", err)
        return err
    }

    // Check if the volumeStorePath directory exists, create if not
    if _, err := os.Stat(m.volumeStorePath); os.IsNotExist(err) {
        if err := os.MkdirAll(m.volumeStorePath, 0755); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", m.volumeStorePath, err)
        }
    }

    cmd := exec.Command("btrfs", "subvolume", "create", volumePath)
    if err := cmd.Run(); err != nil {
        log.Printf("Error running btrfs command: %v", err)
        return err
    }
    return nil
}

// RemoveVolume deletes a BTRFS subvolume
func (m *BtrfsManager) RemoveVolume(name string) error {
    volumePath := filepath.Join(m.volumeStorePath, name)
    
    if err := m.isBtrfsVolume(volumePath); err != nil {
        log.Printf("Volume check failed: %v", err)
        return err
    }

    cmd := exec.Command("btrfs", "subvolume", "delete", volumePath)
    return cmd.Run()
}

// MountVolume logs the mounting of a BTRFS volume
func (m *BtrfsManager) MountVolume(name string) error {
    volumePath := filepath.Join(m.volumeStorePath, name)
    
    if err := m.isBtrfsVolume(volumePath); err != nil {
        log.Printf("Volume check failed: %v", err)
        return err
    }

    log.Printf("Mounting BTRFS volume: %s", volumePath)
    return nil
}

// UnmountVolume logs the unmounting of a BTRFS volume
func (m *BtrfsManager) UnmountVolume(name string) error {
    volumePath := filepath.Join(m.volumeStorePath, name)
    
    if err := m.isBtrfsVolume(volumePath); err != nil {
        log.Printf("Volume check failed: %v", err)
        return err
    }

    log.Printf("Unmounting BTRFS volume: %s", volumePath)
    return nil
}

// GetMountPoint returns the full path of the volume
func (m *BtrfsManager) GetMountPoint(name string) string {
    return filepath.Join(m.volumeStorePath, name)
}
