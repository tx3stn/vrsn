#!/usr/bin/env bats

# e2e tests for the `vrsn get` command with git tags

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-get-tag'

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
	git tag -a "0.0.1" -m "Release 0.0.1" -f
}

teardown() {
	echo "### test teardown ###"
	load ./teardown-git.sh
	tidy-git-changes "$main_branch" "$test_branch"
	delete-tags
}

@test "vrsn get w. git tags: prints the current tag" {
	run vrsn get --git-tag
	assert_success
	assert_line --index 0 '0.0.1'
}

@test "vrsn get w. git tags: prints the latest tag when there are multiple" {
	git tag -a "0.1.0" -m "Release 0.1.0"

	run vrsn get --git-tag
	assert_success
	assert_line --index 0 '0.1.0'
}

@test "vrsn get w. git tags: version sorts tags rather than lexicographic" {
	git tag -a "0.0.2" -m "Release 0.0.2"
	git tag -a "0.0.10" -m "Release 0.0.10"

	run vrsn get --git-tag
	assert_success
	assert_line --index 0 '0.0.10'
}

@test "vrsn get w. git tags: errors when no tags exist" {
	load ./teardown-git.sh
	delete-tags

	run vrsn get --git-tag
	assert_failure
	assert_output --partial 'no git tags found'

	git tag -a "0.0.1" -m "Release 0.0.1"
}

@test "vrsn get w. git tags: in config file" {
	cfg_file="$BATS_TEST_DIRNAME/tag.toml"
	run vrsn get --config="$cfg_file"
	assert_success
	assert_line --index 0 '0.0.1'
}
