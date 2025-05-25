#!/usr/bin/sh

# script to provision e2e test git repo using git tags.
git config --global --add safe.directory "$1"
cd "$1" || exit

git init
echo "# example" >README.md
git add README.md
git commit -m "initial commit"
git tag -a "0.0.1" -m "Release 0.0.1"
