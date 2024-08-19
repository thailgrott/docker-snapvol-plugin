#!/bin/sh

basepath="/opt/snapvol-plugin"
plugin_binary="snapvol"

docker build --no-cache -t snapvol .

mkdir -p ${basepath}

docker run -d --name temp_snapvol ${plugin_binary} sleep infinity

docker cp temp_snapvol:/usr/local/bin/${plugin_binary} ${basepath}
docker cp temp_snapvol:/usr/local/bin/config.json ${basepath}

docker stop temp_snapvol
docker rm temp_snapvol

mkdir -p ${basepath}/rootfs/usr/local/bin

mv ${basepath}/${plugin_binary} ${basepath}/rootfs/usr/local/bin/
chmod 755 ${basepath}/rootfs/usr/local/bin/${plugin_binary}

tar -czf ${basepath}/snapvol-plugin.tar.gz -C ${basepath} rootfs config.json

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
