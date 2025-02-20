#!/usr/bin/env bash
set -e

# Some env variables
BRANCH="master"
REPO_URL="github.com/khulnasoft/docs.git"
AUTHOR_EMAIL="github-actions[bot]@users.noreply.github.com"
AUTHOR_USERNAME="github-actions[bot]"
VERSION_FILE="template_versions.json"
REPO_DIR="template"
COMMIT_URL="https://github.com/khulnasoft/velocity/template"
DOCUSAURUS_COMMAND="npm run docusaurus -- docs:version:template"

# Set commit author
git config --global user.email "${AUTHOR_EMAIL}"
git config --global user.name "${AUTHOR_USERNAME}"

git clone https://${TOKEN}@${REPO_URL} velocity-docs

# Handle push event
if [ "$EVENT" == "push" ]; then
latest_commit=$(git rev-parse --short HEAD)

for f in $(find . -type f -name "*.md" -not -path "./velocity-docs/*"); do
  log_output=$(git log --oneline "${BRANCH}" HEAD~1..HEAD --name-status -- "${f}")

    if [[ $log_output != "" || ! -f "velocity-docs/docs/${REPO_DIR}/$f" ]]; then
      mkdir -p velocity-docs/docs/${REPO_DIR}/$(dirname $f)
      cp "${f}" velocity-docs/docs/${REPO_DIR}/$f
  fi
done

# Handle release event
elif [ "$EVENT" == "release" ]; then
  # Extract package name from tag
  package_name="${TAG_NAME%/*}"
  major_version="${TAG_NAME#*/}"
  major_version="${major_version%%.*}"

  # Form new version name
  new_version="${package_name}_${major_version}.x.x"

  cd velocity-docs/ || true
  npm ci

  # Check if contrib_versions.json exists and modify it if required
  if [[ -f $VERSION_FILE ]]; then
    jq --arg new_version "$new_version" 'del(.[] | select(. == $new_version))' $VERSION_FILE > temp.json && mv temp.json $VERSION_FILE
  fi

  # Run docusaurus versioning command
  $DOCUSAURUS_COMMAND "${new_version}"

  if [[ -f $VERSION_FILE ]]; then
    jq 'sort | reverse' ${VERSION_FILE} > temp.json && mv temp.json ${VERSION_FILE}
  fi
fi

# Push changes
cd velocity-docs/ || true
git add .
if [[ $EVENT == "push" ]]; then
    git commit -m "Add docs from ${COMMIT_URL}/commit/${latest_commit}"
elif [[ $EVENT == "release" ]]; then
    git commit -m "Sync docs for release ${COMMIT_URL}/releases/tag/${TAG_NAME}"
fi

MAX_RETRIES=5
DELAY=5
retry=0

while ((retry < MAX_RETRIES))
do
    git push https://${TOKEN}@${REPO_URL} && break
    retry=$((retry + 1))
    git pull --rebase
    sleep $DELAY
done

if ((retry == MAX_RETRIES))
then
    echo "Failed to push after $MAX_RETRIES attempts. Exiting with 1."
    exit 1
fi

