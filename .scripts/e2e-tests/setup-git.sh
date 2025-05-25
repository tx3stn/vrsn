configure-git() {
	git config --global init.defaultBranch "$1"
	git config --global user.email "int-tests@vrsn.com"
	git config --global user.name "integration tests"
}
