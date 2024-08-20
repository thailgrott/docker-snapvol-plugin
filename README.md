# Snapvol Docker Volume Plugin

## Overview
Snapvol is a Docker volume plugin designed to manage persistent named volumes on a BTRFS subvolumes as the local storage backend. It supports volume creation, removal, listing, mounting, and unmounting.

## Features

- **Efficiency:** Designed to be simple, fast, and have low memory requirements.
- **Snapshot Management:** Create, list, delete, and restore snapshots. (in devolopment)
- **Subvolume Management:** Handles a nested subvolume for each Docker persistent named volume.

## 2. Plugin Functionality

- **Volume Functions:** Implements Docker’s standard volume functions with BTRFS support.
- **Logging & Exception Handling:** Built-in mechanisms for logging and error handling.
- **Socket & API:** Uses a Unix socket and HTTP(s) API for additional functionality.

## Project Structure
```
snapvol/
├── app/
│   ├── btrfs_manager.go
│   ├── main.go
│   └── plugin_api.go
├── config.json
├── Dockerfile
├── go.mod
├── go.sum
├── plugin-installer.sh
└── README.md
```

## Installation

### Building and Installing the Plugin

Clone the project files to a directory to build the plugin. Ensure Docker is running on your system. A selected BTRFS subvolume should be mounted at, '/var/lib/docker-snap-volumes'. The directory /run/docker/plugins/ should exist for the plugin's snapvol.sock unix socket file. The unix socket directory should be readable and writable by the docker group. 

Then run the following installer script: 

```sh
./plugin-installer.sh
```

The installer script should unregister any previous plugin from a previous installation. Next the script should build the docker image from the current project files. A temporary docker container is created to extract the binary and registration file for the plugin. The default plugin base path /opt/snapvol-plugin is created and files moved here. The expected plugin directory structure is created, and the binary is moved to rootfs/usr/local/bin under the plugin base path. The binary is set to executable. A backup archive of the plugin is created. The installer script should try and create the socket directory with the expected permissions. Any previously installed plugin version is removed. The plugin is registered with the config.json in the plugin base path and enabled.

## Using the Plugin

1. Create a volume:

```sh
docker volume create -d snapvol-plugin myvolume
```

2. List volumes:

```sh
docker volume ls
```

3. Remove a volume:

```sh
docker volume rm myvolume
```


