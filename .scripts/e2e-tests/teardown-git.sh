tidy-git-changes() {
	git add .
	git stash
	git stash drop
	git checkout "$1"
	git branch -D "$2"
}
