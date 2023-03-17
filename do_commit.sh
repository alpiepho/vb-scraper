git add "*.json" "*.text" "*.go" "*.sh"
git commit -m "update"
git stash
git checkout gh-pages
git pull
git stash apply stash@{0}
git add -A
git commit -m "update"
git push
git checkout main
git stash pop
git add -A
git commit -m "update"
git push
