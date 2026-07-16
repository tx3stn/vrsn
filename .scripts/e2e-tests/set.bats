#!/usr/bin/env bats

# e2e tests for the `vrsn set` command

main_branch='main'
test_branch='bats-tests'
test_dir='/tmp/project-set'

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

@test "vrsn set w. VERSION file: sets an arbitrary version" {
	git checkout -b "$test_branch"
	run vrsn set 9.9.9
	assert_success
	assert_line --index 0 'version set from 0.0.1 to 9.9.9'

	new=$(head -n1 VERSION)
	assert_equal "9.9.9" "$new"
}

@test "vrsn set w. VERSION file: allows non-increment downgrade" {
	git checkout -b "$test_branch"
	# bump rejects going backwards; set does not check the increment.
	run vrsn set 0.0.0
	assert_success
	assert_line --index 0 'version set from 0.0.1 to 0.0.0'

	new=$(head -n1 VERSION)
	assert_equal "0.0.0" "$new"
}

@test "vrsn set w. VERSION file: invalid version errors" {
	git checkout -b "$test_branch"
	run vrsn set banana
	assert_failure
	assert_output --partial 'error parsing version'

	# the version file is left unchanged when the version is invalid.
	new=$(head -n1 VERSION)
	assert_equal "0.0.1" "$new"
}

@test "vrsn set w. VERSION file: valid set --file" {
	git checkout -b "$test_branch"
	file='package.json'
	printf '{"version":"v0.6.32"}' >"$file"
	run vrsn set v1.2.3 --file="$file"
	assert_success
	assert_line --index 0 'version set from v0.6.32 to v1.2.3'

	new=$(cut -d\" -f4 <"$file")
	assert_equal "v1.2.3" "$new"

	# the VERSION file in the directory is left untouched.
	untouched=$(head -n1 VERSION)
	assert_equal "0.0.1" "$untouched"
	rm "$file"
}

@test "vrsn set w. AndroidManifest.xml: --android-version-code sets both attrs" {
	git checkout -b "$test_branch"
	file='AndroidManifest.xml'
	printf '<manifest\n    android:versionCode="10203"\n    android:versionName="1.2.3">\n</manifest>\n' >"$file"
	run vrsn set 3.4.5 --file="$file" --android-version-code
	assert_success
	assert_line --index 0 'version set from 1.2.3 to 3.4.5'

	new=$(grep -o 'android:versionName="[^"]*"' "$file" | cut -d\" -f2)
	assert_equal "3.4.5" "$new"
	# 3*10000 + 4*100 + 5 = 30405
	code=$(grep -o 'android:versionCode="[^"]*"' "$file" | cut -d\" -f2)
	assert_equal "30405" "$code"
	rm "$file"
}

@test "vrsn set w. AndroidManifest.xml: android-version-code from config file" {
	git checkout -b "$test_branch"
	file='AndroidManifest.xml'
	printf '<manifest\n    android:versionCode="10203"\n    android:versionName="1.2.3">\n</manifest>\n' >"$file"

	cfg_file="$BATS_TEST_DIRNAME/set.toml"
	run vrsn set 3.4.5 --file="$file" --config="$cfg_file"
	assert_success
	assert_line --index 0 'version set from 1.2.3 to 3.4.5'

	new=$(grep -o 'android:versionName="[^"]*"' "$file" | cut -d\" -f2)
	assert_equal "3.4.5" "$new"
	code=$(grep -o 'android:versionCode="[^"]*"' "$file" | cut -d\" -f2)
	assert_equal "30405" "$code"
	rm "$file"
}
