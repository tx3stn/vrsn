tidy-git-changes() {
	git add .
	git stash
	git stash drop
	git checkout "$1"

	branches=$(git --no-pager branch --list "$2")
	if [ "$branches" != "" ]; then
		git branch -D "$2"
	fi
}

delete-tags() {
	git --no-pager tag -d $(git --no-pager tag -l)
}
