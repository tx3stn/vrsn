#!/usr/bin/env bats

# e2e tests for the `vrsn get` command

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-get'

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

@test "vrsn get w. VERSION file: prints the current version" {
	run vrsn get
	assert_success
	assert_line --index 0 '0.0.1'
}

@test "vrsn get w. --file flag: prints the version from the specified file" {
	printf '{"version":"1.2.3"}' >package.json

	run vrsn get --file package.json
	assert_success
	assert_line --index 0 '1.2.3'
}

@test "vrsn get: errors when multiple version files found and no --file flag" {
	printf '{"version":"1.2.3"}' >package.json

	run vrsn get
	assert_failure
	assert_output --partial 'multiple version files found in directory'
}

@test "vrsn get w. files in config: prints the version in every file" {
	printf '{"version":"0.0.1"}' >package.json

	cfg_file="$BATS_TEST_DIRNAME/multi-file.toml"
	run vrsn get --config="$cfg_file"
	assert_success
	assert_line --index 0 'VERSION: 0.0.1'
	assert_line --index 1 'package.json: 0.0.1'
}

@test "vrsn get w. files in config: prints mismatched versions without erroring" {
	printf '{"version":"0.0.2"}' >package.json

	cfg_file="$BATS_TEST_DIRNAME/multi-file.toml"
	run vrsn get --config="$cfg_file"
	assert_success
	assert_line --index 0 'VERSION: 0.0.1'
	assert_line --index 1 'package.json: 0.0.2'
}

@test "vrsn get: errors when no version files found" {
	mkdir empty-dir
	cd empty-dir || exit 1

	run vrsn get
	assert_failure
	assert_output --partial 'no version files found in directory'

	cd "$test_dir" || exit 1
	rmdir empty-dir
}
