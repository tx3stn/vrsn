#!/usr/bin/sh

# script to provision integration test git repo using a version file.
cd /project-vf || exit

git init
echo "0.0.1" >VERSION
git add VERSION
git commit -m "initial commit"
