#!/bin/bash
# A script to fetch the latest stable versions and then create a new Git release branch and tag for specific Terraform providers.

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Function to get the latest stable version from a Git repository ---
get_latest_version() {
  local repo_url="$1"
  # Fetch all tags, sort them by version, and get the latest stable version (not pre-release).
  # We use grep to filter for tags that match the vX.Y.Z pattern, excluding any with hyphens (e.g., v1.2.3-beta).
  local latest_version=$(git ls-remote --tags --sort='-v:refname' "$repo_url" | grep -o 'v[0-9]*\.[0-9]*\.[0-9]*$' | head -n 1)
  
  if [ -z "$latest_version" ]; then
    echo "Version not found"
  else
    # Remove the 'v' prefix for cleaner output
    echo "${latest_version:1}"
  fi
}

# Small helper to confirm an action
confirm() {
  local prompt="$1"
  echo ""
  read -p "$prompt (y/n) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Operation cancelled."
    exit 0
  fi
}

# --- Fetch and Display Latest Stable Versions ---
echo "--- Fetching Latest Stable Provider Versions ---"

# Define the GitHub repositories for each provider.
REPOSITORIES=(
  "jfrog/terraform-provider-shared"
)

# Loop through each repository, fetch its latest version, and display it.
for repo in "${REPOSITORIES[@]}"; do
  provider_name=$(basename "$repo")
  repo_url="https://github.com/${repo}"
  latest=$(get_latest_version "$repo_url")
  echo "Latest version for ${provider_name}: v$latest"
done

echo "-------------------------------------"
echo ""

# --- Inputs ---
PROVIDER_NAME="terraform-provider-shared"
echo "Using provider: ${PROVIDER_NAME}"
read -p "Please enter the new version number (e.g., 1.2.3): " NEW_VERSION

# Add 'v' prefix to the new version if it doesn't exist.
if [[ ! "$NEW_VERSION" =~ ^v ]]; then
  NEW_VERSION="v$NEW_VERSION"
fi

# --- Determine the correct branch to use ---
BRANCH_TO_CHECKOUT=""
case "$PROVIDER_NAME" in
  "terraform-provider-shared")
    BRANCH_TO_CHECKOUT="main"
    ;;
  *)
    echo "Error: Unknown provider name '$PROVIDER_NAME'."
    echo "Known providers are: terraform-provider-shared."
    exit 1
    ;;
esac

echo "--- Starting release process for provider '${PROVIDER_NAME}' and version ${NEW_VERSION} ---"

# --- Git Workflow ---
# 1. Checkout the correct base branch.
echo "About to checkout branch '${BRANCH_TO_CHECKOUT}'..."
confirm "Proceed to checkout '${BRANCH_TO_CHECKOUT}'?"
git checkout "${BRANCH_TO_CHECKOUT}"

# 2. Pull the latest code.
echo "About to pull latest code from '${BRANCH_TO_CHECKOUT}'..."
confirm "Proceed to pull from '${BRANCH_TO_CHECKOUT}'?"
git pull

# 3. Checkout a new branch for the release.
echo "About to create and checkout new release branch: ${NEW_VERSION}..."
confirm "Proceed to create branch '${NEW_VERSION}'?"
git checkout -b "${NEW_VERSION}"

# 4. Push the new branch to the remote repository.
echo "About to push new branch to origin: ${NEW_VERSION}..."
confirm "Proceed to push branch '${NEW_VERSION}' to origin?"
git push origin "${NEW_VERSION}"

# 5. Create a new tag from the new branch.
echo "About to create new tag: ${NEW_VERSION}..."
confirm "Proceed to create tag '${NEW_VERSION}'?"
git tag "${NEW_VERSION}"

# 6. Push the new tag to the remote repository.
echo "About to push new tag to origin: ${NEW_VERSION}..."
confirm "Proceed to push tag '${NEW_VERSION}' to origin?"
git push origin tag "${NEW_VERSION}"

echo ""
echo "--- Release process completed successfully for ${PROVIDER_NAME}! ---"

