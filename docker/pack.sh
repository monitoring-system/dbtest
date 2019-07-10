#!/bin/sh

cp -r ../randgen-server/autoyyzz autoyyzz
docker build -t framework-dbtest .
rm -rf autoyyzz
