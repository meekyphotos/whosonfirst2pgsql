#!/usr/bin/env bash

mkdir /tmp/ramdisk
chmod 777 /tmp/ramdisk
sudo mount -t tmpfs -o size=25G whosonfirst /tmp/ramdisk
bzip2 -d whosonfirst-data-admin-latest.tar.bz2
cp -r /data/whosonfirst-data-admin-latest/ /tmp/ramdisk

./whosonfirst2pgsql -i /tmp/ramdisk -db