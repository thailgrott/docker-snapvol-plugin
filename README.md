# Snapvol Docker Volume Plugin

## Overview
Snapvol is a Docker volume plugin designed to manage BTRFS subvolumes. It supports volume creation, removal, mounting, and unmounting.

## Features
- Create, remove, mount, and unmount Docker volumes
- Future support for BTRFS snapshot management

## Project Structure
```
snapvol/
├── app/
│   └── btrfs_manager.go
├── main.go
├── plugin_api.go
├── config.json
├── Dockerfile
└── README.md
```

## Installation

### Building the Plugin
```sh
docker build -t snapvol .
```

### Installing the Plugin

1. Copy the plugin binary to the Docker plugin directory:

```sh
cp snapvol /usr/local/bin/snapvol
```

2. Register the plugin with Docker:

```sh
docker plugin create --config /path/to/config.json
```

## Using the Plugin

1. Create a volume:

```sh
docker volume create -d snap-volume-plugin myvolume
```

2. List volumes:

```sh
docker volume ls
```

3. Remove a volume:

```sh
docker volume rm myvolume
```


