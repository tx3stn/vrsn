#!/usr/bin/env bats

# e2e tests for the `vrsn check` command

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

@test "vrsn check w. VERSION file: no bump" {
	git checkout -b "$test_branch"
	run vrsn check
	assert_failure
	assert_line --index 0 'was: 0.0.1'
	assert_line --index 1 'now: 0.0.1'
	assert_line --index 2 --partial 'version has not been bumped'
}

@test "vrsn check w. VERSION file: valid bump" {
	git checkout -b "$test_branch"
	echo "0.1.0" >VERSION
	run vrsn check
	assert_success
	assert_line --index 0 'was: 0.0.1'
	assert_line --index 1 'now: 0.1.0'
	assert_line --index 2 'valid version bump'
}

@test "vrsn check w. VERSION file: invalid bump" {
	git checkout -b "$test_branch"
	echo "0.2.0" >VERSION
	run vrsn check
	assert_failure
	assert_line --index 0 'was: 0.0.1'
	assert_line --index 1 'now: 0.2.0'
	assert_line --index 2 --partial 'invalid version bump'
}

@test "vrsn check w. VERSION file: --base-branch" {
	git checkout -b "$test_branch"
	echo "0.1.0" >VERSION
	git add VERSION
	git commit -m "bump version"

	git checkout "$main_branch"
	run vrsn check --base-branch "$test_branch"
	assert_failure
	assert_line --index 0 'was: 0.1.0'
	assert_line --index 1 'now: 0.0.1'
	assert_line --index 2 --partial 'invalid version bump'
}

@test "vrsn check w. VERSION file: base-branch in config file" {
	git checkout -b "$test_branch"
	echo "0.1.0" >VERSION
	git add VERSION
	git commit -m "bump version"
	git checkout "$main_branch"
	sleep 0.2

	cfg_file="$BATS_TEST_DIRNAME/check.toml"

	cat "$cfg_file"
	run vrsn check --config="$cfg_file"
	assert_failure
	assert_line --index 0 'was: 0.1.0'
	assert_line --index 1 'now: 0.0.1'
	assert_line --index 2 --partial 'invalid version bump'
}

@test "vrsn check w. --was & --now flags" {
	run vrsn check --was 9.0.0 --now 10.0.0
	assert_success
	assert_line --index 0 'was: 9.0.0'
	assert_line --index 1 'now: 10.0.0'
	assert_line --index 2 'valid version bump'
}
