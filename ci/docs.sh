#!/usr/bin/env bash

set -e
TAGNAME=$1
GH_REPO="${GITHUB_REPOSITORY:-qa-guru/cm}"
GH_REF="github.com/${GH_REPO}.git"

if [[ -z "${GITHUB_TOKEN:-}" ]]; then
	echo "Skipping docs publish: GITHUB_TOKEN not set"
	exit 0
fi

git config user.name "${GITHUB_ACTOR:-github-actions}"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git remote add upstream "https://x-access-token:${GITHUB_TOKEN}@${GH_REF}" 2>/dev/null || \
	git remote set-url upstream "https://x-access-token:${GITHUB_TOKEN}@${GH_REF}"

if ! git ls-remote --exit-code upstream gh-pages >/dev/null 2>&1; then
	echo "Skipping docs publish: gh-pages branch not found on ${GH_REPO}"
	exit 0
fi

git fetch upstream

echo "Deleting old output"
rm -rf ${GITHUB_WORKSPACE}/docs/output
mkdir ${GITHUB_WORKSPACE}/docs/output
git worktree prune
rm -rf ${GITHUB_WORKSPACE}/.git/worktrees/docs/output/

echo "Checking out gh-pages branch into docs/output"
git worktree add -B gh-pages ${GITHUB_WORKSPACE}/docs/output upstream/gh-pages

echo "Removing existing files"
mkdir -p ${GITHUB_WORKSPACE}/docs/output/${TAGNAME}
rm -rf ${GITHUB_WORKSPACE}/docs/output/${TAGNAME}/*

echo "Copying images"
cp -R ${GITHUB_WORKSPACE}/docs/img ${GITHUB_WORKSPACE}/docs/output/${TAGNAME}/img
echo "Copying files to root"
cp -Rv ${GITHUB_WORKSPACE}/docs/files/* ${GITHUB_WORKSPACE}/docs/output

echo "Generating docs"
docker run -v ${GITHUB_WORKSPACE}/docs/:/documents/ --name asciidoc-to-html asciidoctor/docker-asciidoctor asciidoctor -a revnumber=${TAGNAME} -D /documents/output/${TAGNAME} index.adoc

echo "Updating gh-pages branch"
cd ${GITHUB_WORKSPACE}/docs/output && git add --all && git status && git commit -m "Publishing to gh-pages"

git push upstream HEAD:gh-pages
