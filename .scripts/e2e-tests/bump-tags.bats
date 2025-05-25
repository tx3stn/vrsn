#!/usr/bin/env bats

# e2e tests for the `vrsn bump` command

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-tag'

setup_file() {
	echo "### suite setup ###"
	load ./setup-git.sh
	configure-git "$main_branch"

	load ./setup-git-repo.sh
	setup-git-repo-with-tags "$test_dir"
}

teardown_file() {
	echo "### suite teardown ###"
	rm -rf "$test_dir"
}

setup() {
	echo "### test setup ###"
	bats_load_library bats-support
	bats_load_library bats-assert
	cd "$test_dir" || exit 1
}

teardown() {
	echo "### test teardown ###"
	load ./teardown-git.sh
	tidy-git-changes "$main_branch" "$test_branch"
	delete-tags
}

@test "vrsn bump w. git tags: valid bump" {
	git checkout -b "$test_branch"
	echo "update" >>README.md
	git add README.md
	git commit -m "update"

	run vrsn bump patch --git-tag
	assert_success
	assert_line --index 0 'git tag version bumped from 0.0.1 to 0.0.2'

	new=$(git --no-pager tag --list --points-at HEAD)
	assert_equal "0.0.2" "$new"
}
