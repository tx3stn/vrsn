#!/usr/bin/sh

# script to provision integration test git repo using a version file.
cd /project-tag || exit

git init
echo "# example" >README.md
git add README.md
git commit -m "initial commit"
git tag -a "0.0.1" -m "Release 0.0.1"
