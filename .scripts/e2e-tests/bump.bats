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

@test "vrsn bump w. VERSION file: interactive bump" {
	git checkout -b "$test_branch"
	expect_script="bump-interactive.exp"
	new_path="$test_dir/$expect_script"
	cp "$BATS_TEST_DIRNAME/$expect_script" "$new_path"

	eval "run $new_path"
	assert_success

	new=$(head -n1 VERSION)
	assert_equal "0.1.0" "$new"
	rm "$new_path"
}

@test "vrsn bump w. VERSION file: invalid bump" {
	git checkout -b "$test_branch"
	run vrsn bump fail
	assert_failure
	assert_line --index 0 'invalid argument "fail" for "vrsn bump"'

	new=$(head -n1 VERSION)
	assert_equal "0.0.1" "$new"
}

@test "vrsn bump w. VERSION file: valid bump --file" {
	git checkout -b "$test_branch"
	file='package.json'
	printf '{"version":"v0.6.32"}' >"$file"
	run vrsn bump patch --file="$file"
	assert_success
	assert_line --index 0 'version bumped from v0.6.32 to v0.6.33'

	new=$(cut -d\" -f4 <"$file")
	assert_equal "v0.6.33" "$new"
	rm "$file"
}

@test "vrsn bump w. AndroidManifest.xml: valid bump --file" {
	git checkout -b "$test_branch"
	file='AndroidManifest.xml'
	printf '<manifest\n    android:versionCode="10203"\n    android:versionName="1.2.3">\n</manifest>\n' >"$file"
	run vrsn bump minor --file="$file"
	assert_success
	assert_line --index 0 'version bumped from 1.2.3 to 1.3.0'

	new=$(grep -o 'android:versionName="[^"]*"' "$file" | cut -d\" -f2)
	assert_equal "1.3.0" "$new"
	# versionCode is left untouched.
	code=$(grep -o 'android:versionCode="[^"]*"' "$file" | cut -d\" -f2)
	assert_equal "10203" "$code"
	rm "$file"
}

@test "vrsn bump w. AndroidManifest variant: matches glob --file" {
	git checkout -b "$test_branch"
	file='AndroidManifest.debug.xml'
	printf '<manifest\n    android:versionName="0.6.32">\n</manifest>\n' >"$file"
	run vrsn bump patch --file="$file"
	assert_success
	assert_line --index 0 'version bumped from 0.6.32 to 0.6.33'

	new=$(grep -o 'android:versionName="[^"]*"' "$file" | cut -d\" -f2)
	assert_equal "0.6.33" "$new"
	rm "$file"
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

@test "vrsn bump w. VERSION file: --commit-msg with template variable" {
	git checkout -b "$test_branch"
	run vrsn bump minor --commit --commit-msg 'release {{.Version}}'
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.1.0'

	new=$(head -n1 VERSION)
	assert_equal "0.1.0" "$new"

	run git --no-pager log --oneline -n 1
	assert_success
	assert_line --index 0 --partial "release 0.1.0"
	git reset "$(git rev-parse HEAD^1)"
}

@test "vrsn bump w. VERSION file: invalid --commit-msg template" {
	git checkout -b "$test_branch"
	run vrsn bump minor --commit --commit-msg 'release {{.Version'
	assert_failure
	assert_output --partial 'error parsing message template'

	new=$(head -n1 VERSION)
	assert_equal "0.0.1" "$new"
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
