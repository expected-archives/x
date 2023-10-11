# This script walk through all subdirectories and run go mod tidy if go.mod exists.
# It prints out the directory name if go mod tidy is run.

set -e

find . -name "go.mod" | xargs -n1 dirname | while read dir; do
  echo "go mod tidy in $dir"
  (cd "$dir" && go mod tidy)
done
