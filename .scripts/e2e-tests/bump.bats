#!/usr/bin/env bats

# e2e tests for the `vrsn bump` command

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-vf'
commit_msg='testing commit'

setup_file() {
	echo "### suite setup ###"
	load ./setup-git.sh
	configure-git "$main_branch"

	load ./setup-git-repo.sh
	setup-git-repo-with-version-file "$test_dir"
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
}

@test "vrsn bump w. VERSION file: valid bump" {
	git checkout -b "$test_branch"
	run vrsn bump major
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 1.0.0'

	new=$(head -n1 VERSION)
	assert_equal "1.0.0" "$new"
}

@test "vrsn bump w. VERSION file: invalid bump" {
	git checkout -b "$test_branch"
	run vrsn bump fail
	assert_failure
	assert_line --index 0 'invalid argument "fail" for "vrsn bump"'

	new=$(head -n1 VERSION)
	assert_equal "0.0.1" "$new"
}

@test "vrsn bump w. VERSION file: --commit default commit message" {
	git checkout -b "$test_branch"
	run vrsn bump minor --commit
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.1.0'

	new=$(head -n1 VERSION)
	assert_equal "0.1.0" "$new"

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial "bump version"
	git reset "$(git rev-parse HEAD^1)"
}

@test "vrsn bump w. VERSION file: --commit --commit-msg" {
	git checkout -b "$test_branch"
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

@test "vrsn bump w. VERSION file: commit in config file" {
	git checkout -b "$test_branch"

	cfg_file="$BATS_TEST_DIRNAME/check.toml"
	run vrsn bump minor --config="$cfg_file"
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.1.0'

	new=$(head -n1 VERSION)
	assert_equal "0.1.0" "$new"

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial "$commit_msg"
	git reset "$(git rev-parse HEAD^1)"
}
