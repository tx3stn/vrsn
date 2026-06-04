#!/usr/bin/env bats

# e2e tests for the `vrsn bump` command combining a version file and git tags

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-dual'
commit_msg='testing commit'
tag_msg='custom tag message'

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
	git tag -a "0.0.1" -m "Release 0.0.1" -f
}

teardown() {
	echo "### test teardown ###"
	load ./teardown-git.sh
	tidy-git-changes "$main_branch" "$test_branch"
	delete-tags
}

@test "vrsn bump w. file and git tag: bumps file, commits and tags the commit" {
	git checkout -b "$test_branch"
	run vrsn bump patch --git-tag --file=VERSION --commit
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.0.2'
	assert_line --index 1 'version file committed'
	assert_line --index 2 'git tag 0.0.2 added'

	new=$(head -n1 VERSION)
	assert_equal "0.0.2" "$new"

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial "bump version"

	new_tag=$(git --no-pager tag --list --points-at HEAD)
	assert_equal "0.0.2" "$new_tag"
}

@test "vrsn bump w. file and git tag: --commit-msg and --tag-msg" {
	git checkout -b "$test_branch"
	run vrsn bump patch --git-tag --file=VERSION --commit --commit-msg="$commit_msg" --tag-msg="$tag_msg"
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.0.2'

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial "$commit_msg"

	run git --no-pager tag --list --points-at HEAD -n1
	assert_success
	assert_line --index 0 --partial "$tag_msg"
}

@test "vrsn bump w. file and git tag: errors without --commit" {
	git checkout -b "$test_branch"
	run vrsn bump patch --git-tag --file=VERSION
	assert_failure
	assert_line --index 0 --partial 'cannot combine --git-tag with version files unless commit is enabled'

	new=$(head -n1 VERSION)
	assert_equal "0.0.1" "$new"

	run git --no-pager tag --list '0.0.2'
	assert_output ''
}

@test "vrsn bump w. file and git tag: in config file" {
	git checkout -b "$test_branch"

	cfg_file="$BATS_TEST_DIRNAME/dual.toml"
	run vrsn bump minor --config="$cfg_file" --file=VERSION
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.1.0'
	assert_line --index 1 'version file committed'
	assert_line --index 2 'git tag 0.1.0 added'

	new=$(head -n1 VERSION)
	assert_equal "0.1.0" "$new"

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial "$commit_msg"

	run git --no-pager tag --list --points-at HEAD -n1
	assert_success
	assert_line --index 0 --partial "$tag_msg"
}

@test "vrsn bump w. file and git tag: --git-tag alone stays tag only" {
	git checkout -b "$test_branch"
	echo "update" >note.txt
	git add note.txt
	git commit -m "update"

	run vrsn bump patch --git-tag
	assert_success
	assert_line --index 0 'git tag version bumped from 0.0.1 to 0.0.2'

	new=$(head -n1 VERSION)
	assert_equal "0.0.1" "$new"

	new_tag=$(git --no-pager tag --list --points-at HEAD)
	assert_equal "0.0.2" "$new_tag"

	run git --no-pager log --oneline -n 1
	assert_line --index 0 --partial "update"
}
