#!/usr/bin/env bats

# e2e tests for config file resolution (project level and global configs)

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-config'
xdg_dir='/tmp/vrsn-test-xdg'
fake_home='/tmp/vrsn-test-home'

# writes a config file to the path in $1 with the commit message in $2 so
# tests can assert which config file was used.
write-config() {
	cat >"$1" <<-EOF
		verbose = false

		[bump]
		commit = true
		commit-msg = '$2'
		git-tag = false
		tag-msg = ''

		[check]
		base-branch = 'main'
	EOF
}

setup_file() {
	echo "### suite setup ###"
	load ./setup-git.sh
	configure-git "$main_branch"

	load ./setup-git-repo.sh
	setup-git-repo-with-version-file "$test_dir"

	mkdir -p "$xdg_dir"
	write-config "$xdg_dir/vrsn.toml" 'global config commit'

	# A fake home directory needs a copy of the git config so the git
	# commands vrsn runs still work when HOME is overridden.
	mkdir -p "$fake_home/.config"
	write-config "$fake_home/.config/vrsn.toml" 'home config commit'
	cp "$HOME/.gitconfig" "$fake_home/.gitconfig"
}

teardown_file() {
	echo "### suite teardown ###"
	rm -rf "$test_dir" "$xdg_dir" "$fake_home"
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

@test "vrsn config resolution: project vrsn.toml in current directory" {
	git checkout -b "$test_branch"
	write-config vrsn.toml 'project config commit'

	run vrsn bump patch
	assert_success
	assert_line --index 0 'version bumped from 0.0.1 to 0.0.2'
	assert_line --index 1 'version file committed'

	run git --no-pager log --oneline -n 1
	assert_line --index 0 --partial 'project config commit'

	rm vrsn.toml
	git reset "$(git rev-parse HEAD^1)"
}

@test "vrsn config resolution: project config takes precedence over global config" {
	git checkout -b "$test_branch"
	write-config vrsn.toml 'project config commit'

	run env XDG_CONFIG_DIR="$xdg_dir" vrsn bump patch
	assert_success

	run git --no-pager log --oneline -n 1
	assert_line --index 0 --partial 'project config commit'

	rm vrsn.toml
	git reset "$(git rev-parse HEAD^1)"
}

@test "vrsn config resolution: XDG_CONFIG_DIR config used when no project config" {
	git checkout -b "$test_branch"

	run env XDG_CONFIG_DIR="$xdg_dir" vrsn bump patch
	assert_success

	run git --no-pager log --oneline -n 1
	assert_line --index 0 --partial 'global config commit'

	git reset "$(git rev-parse HEAD^1)"
}

@test "vrsn config resolution: HOME config used when no project or XDG config" {
	git checkout -b "$test_branch"

	run env -u XDG_CONFIG_DIR HOME="$fake_home" vrsn bump patch
	assert_success

	run git --no-pager log --oneline -n 1
	assert_line --index 0 --partial 'home config commit'

	git reset "$(git rev-parse HEAD^1)"
}

@test "vrsn config resolution: --config flag takes precedence over project config" {
	git checkout -b "$test_branch"
	write-config vrsn.toml 'project config commit'

	run vrsn bump patch --config="$BATS_TEST_DIRNAME/check.toml"
	assert_success

	run git --no-pager log --oneline -n 1
	assert_line --index 0 --partial 'testing commit'

	rm vrsn.toml
	git reset "$(git rev-parse HEAD^1)"
}
