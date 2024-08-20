#!/bin/sh

basepath="/opt/snapvol-plugin"
plugin_binary="snapvol"

docker build --no-cache -t ${plugin_binary} .

docker plugin disable snapvol-plugin:latest

mkdir -p ${basepath}

docker stop temp_snapvol
docker rm temp_snapvol

docker run --rm -it --name temp_snapvol --entrypoint /bin/sh \
-v ${basepath}:/mnt snapvol \
-c "cd / && tar czf /mnt/snapvol-plugin.tar.gz bin sbin usr/bin usr/sbin lib usr/lib etc config.json"

tar -xzf ${basepath}/snapvol-plugin.tar.gz -C ${basepath} config.json
tar -xzf ${basepath}/snapvol-plugin.tar.gz -C ${basepath}/rootfs/ .

# Make sure the socket directory exists and is writable by the Docker daemon.
mkdir -p /run/docker/plugins
chown root:docker /run/docker/plugins
chmod 755 /run/docker/plugins

# make sure the volume directory exists and is writable by the Docker daemon.
mkdir -p /var/lib/docker-snap-volumes
chown root:docker /var/lib/docker-snap-volumes
chmod 755 /var/lib/docker-snap-volumes

docker plugin rm snapvol-plugin

docker plugin create snapvol-plugin ${basepath}

docker plugin enable snapvol-plugin
