source "./setup-git.sh"
source "./setup-git-repo.sh"

configure-git "main"

dir="/debug-project"
setup-git-repo-with-version-file "$dir"
