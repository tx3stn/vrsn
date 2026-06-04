#!/usr/bin/env bats

# e2e tests for the `files` config option with multiple version files

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-multi'

setup_file() {
	echo "### suite setup ###"
	load ./setup-git.sh
	configure-git "$main_branch"

	load ./setup-git-repo.sh
	setup-git-repo-with-version-file "$test_dir"

	printf '{"version":"0.0.1"}' >package.json
	git add package.json
	git commit -m "add package.json"
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

@test "vrsn bump w. files in config: bumps all files in a single commit" {
	git checkout -b "$test_branch"

	cfg_file="$BATS_TEST_DIRNAME/multi-file.toml"
	run vrsn bump minor --config="$cfg_file"
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.1.0'
	assert_line --index 1 'version file committed'

	assert_equal "0.1.0" "$(head -n1 VERSION)"
	assert_equal "0.1.0" "$(cut -d\" -f4 <package.json)"

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial 'bump to 0.1.0'

	run git --no-pager diff-tree --no-commit-id --name-only -r HEAD
	assert_success
	assert_line --partial 'VERSION'
	assert_line --partial 'package.json'

	git reset "$(git rev-parse HEAD^1)"
}

@test "vrsn bump w. files in project config: picks up vrsn.toml from current directory" {
	git checkout -b "$test_branch"

	cp "$BATS_TEST_DIRNAME/multi-file.toml" vrsn.toml
	run vrsn bump minor
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.1.0'
	assert_line --index 1 'version file committed'

	assert_equal "0.1.0" "$(head -n1 VERSION)"
	assert_equal "0.1.0" "$(cut -d\" -f4 <package.json)"

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial 'bump to 0.1.0'

	rm vrsn.toml
	git reset "$(git rev-parse HEAD^1)"
}

@test "vrsn bump w. files in config: errors when versions do not match" {
	git checkout -b "$test_branch"
	printf '{"version":"0.0.2"}' >package.json

	cfg_file="$BATS_TEST_DIRNAME/multi-file-verbose.toml"
	run vrsn bump patch --config="$cfg_file"
	assert_failure
	assert_line --partial 'file VERSION has version 0.0.1'
	assert_line --partial 'file package.json has version 0.0.2'
	assert_output --partial 'version files do not contain matching versions'

	assert_equal "0.0.1" "$(head -n1 VERSION)"
}

@test "vrsn check w. files in config: valid bump" {
	git checkout -b "$test_branch"
	echo "0.1.0" >VERSION
	printf '{"version":"0.1.0"}' >package.json

	cfg_file="$BATS_TEST_DIRNAME/multi-file.toml"
	run vrsn check --config="$cfg_file"
	assert_success
	assert_line --index 0 'was: 0.0.1'
	assert_line --index 1 'now: 0.1.0'
	assert_line --index 2 'valid version bump'
}

@test "vrsn check w. files in config: errors when versions do not match" {
	git checkout -b "$test_branch"
	echo "0.1.0" >VERSION
	printf '{"version":"0.2.0"}' >package.json

	cfg_file="$BATS_TEST_DIRNAME/multi-file.toml"
	run vrsn check --config="$cfg_file"
	assert_failure
	assert_output --partial 'version files do not contain matching versions'
}
