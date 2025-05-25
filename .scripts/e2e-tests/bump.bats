#!/usr/bin/env bats

# e2e tests for the `vrsn bump` command

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-vf'

setup_file() {
	echo "### suite setup ###"
	load ./setup-git.sh
	configure-git "$main_branch"

	load ./setup-git-repo.sh
	setup-git-repo-with-version-file "$test_dir"
}

teardown_file() {
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
}

@test "vrsn bump w. VERSION file: valid bump" {
	git checkout -b "$test_branch"
	run vrsn bump major
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 1.0.0'

	new=$(head -n1 VERSION)
	assert_equal "1.0.0" "$new"
}

@test "vrsn bump w. VERSION file: --commit" {
	git checkout -b "$test_branch"
	commit_msg='testing commit'
	run vrsn bump minor --commit --commit-msg "$commit_msg"
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.1.0'

	new=$(head -n1 VERSION)
	assert_equal "0.1.0" "$new"

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial "$commit_msg"
	git reset "$(git rev-parse HEAD^1)"
}
