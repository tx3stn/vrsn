#!/usr/bin/sh

# script to provision integration test git repo using a version file.
cd /project-tag || exit

git config --global init.defaultBranch "main"
git config --global user.email "int-tests@vrsn.com"
git config --global user.name "integration tests"
git init
echo "# example" >README.md
git add README.md
git commit -m "initial commit"
git tag -a "0.0.1" -m "Release 0.0.1"
