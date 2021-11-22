#!/bin/bash
cd "$(dirname "$(dirname "$0")}")" || return

BRANCH="$(git branch --show-current)"
if [[ ! $BRANCH == "main" ]]; then
  echo "$BRANCH is not main branch"
  exit 1
fi

VERSION="$(cat version)"
echo "old: $VERSION"
IFS="."
read -ra ADDR <<<"$VERSION"
V_MAJOR="${ADDR[0]}"
V_MINOR="${ADDR[1]}"
V_PATCH="${ADDR[2]}"
V_PATCH="$(expr "$V_PATCH" + 1)"
VERSION="$V_MAJOR.$V_MINOR.$V_PATCH"
rm -f version
echo "$VERSION" >>version
echo "new: $VERSION"
git add version
git commit -m "Increase version"
git push origin main

scripts/_build.sh || exit 1
scripts/_nsi.sh || exit 1

git tag v"$VERSION"
git push origin --tags

TOKEN_HEADER="Authorization: token $(cat githubtoken)"

RELEASE_ID=$(
  curl \
    -X POST \
    -H "$TOKEN_HEADER" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{\"tag_name\": \"v$VERSION\", \"body\": \"[Release Notes](https://github.com/wirekang/mouseable/blob/main/release-notes.md#v$VERSION)\"}" \
    "https://api.github.com/repos/wirekang/mouseable/releases" |
    awk '/"id":/{print substr($2, 1, length($2)-1); exit;}'
)

echo "RELEASE_ID: $RELEASE_ID"

uploadAsset(){
 FILE="$1"

 curl \
   -X POST \
   -H "$TOKEN_HEADER" \
   -H "Accept: application/vnd.github.v3+json" \
   -H "Content-Type: $(file -b --mime-type "$FILE")" \
   --data-binary @"$FILE" \
   "https://uploads.github.com/repos/wirekang/mouseable/releases/$RELEASE_ID/assets?name=$(basename "$FILE")"
}

uploadAsset "build/installer.exe"
uploadAsset "build/portable.exe"
