#!/usr/bin/sh

# script to provision e2e test git repo using a version file.
git config --global --add safe.directory "$1"
cd "$1" || exit

git init
echo "0.0.1" >VERSION
git add VERSION
git commit -m "initial commit"
