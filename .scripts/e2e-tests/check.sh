#!/usr/bin/env bats

# e2e tests for the `vrsn check` command

main_branch='main'
test_branch='bats-tests'

setup() {
	echo "test setup..."
	bats_load_library bats-support
	bats_load_library bats-assert
	cd /project-vf || exit 1
}

teardown() {
	echo "test teardown..."
	git add .
	git stash
	git stash drop
	git checkout "$main_branch"
	git branch -D "$test_branch"
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
