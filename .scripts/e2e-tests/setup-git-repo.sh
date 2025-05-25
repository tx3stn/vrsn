#!/usr/bin/env bash

# script to provision e2e test git repo using a version file.
setup-git-repo-with-version-file() {
	mkdir "$1"
	git config --global --add safe.directory "$1"
	cd "$1" || exit

	git init
	echo "0.0.1" >VERSION
	git add VERSION
	git commit -m "initial commit"
}

# script to provision e2e test git repo using git tags.
setup-git-repo-with-tags() {
	mkdir "$1"
	git config --global --add safe.directory "$1"
	cd "$1" || exit

	git init
	echo "# example" >README.md
	git add README.md
	git commit -m "initial commit"
	git tag -a "0.0.1" -m "Release 0.0.1"
}
