#!/usr/bin/env bash

mkdir /tmp/ramdisk
chmod 777 /tmp/ramdisk
mount -t tmpfs -o size=25G whosonfirst /tmp/ramdisk
bzip2 -d whosonfirst-data-admin-latest.tar.bz2
cp /data/whosonfirst-data-admin-latest.tar /tmp/ramdisk