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
	git checkout "$main_branch"
	git branch -D "$test_branch"
}

@test "vrsn check with VERSION file: no bump" {
	git checkout -b "$test_branch"
	run vrsn check
	assert_failure
	assert_line --index 0 'was: 0.0.1'
	assert_line --index 1 'now: 0.0.1'
	assert_line --index 2 --partial 'version has not been bumped'
}
