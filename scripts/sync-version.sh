#!/bin/bash

# sync-version.sh - Automatically sync version across Go code, tests, and README.md

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.2.0"
    exit 1
fi

NEW_VERSION="$1"

# Validate version format
if [[ ! $NEW_VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Version must be in format v1.2.3"
    exit 1
fi

echo "Syncing version to: $NEW_VERSION"

# Cross-platform sed function
sed_inplace() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "$1" "$2"
    else
        sed -i "$1" "$2"
    fi
}

# Update version.go
if [ -f "pkg/version/version.go" ]; then
    sed_inplace "s/const Version = \".*\"/const Version = \"$NEW_VERSION\"/" "pkg/version/version.go"
    echo "Updated pkg/version/version.go"
fi

# Update version_test.go
if [ -f "pkg/version/version_test.go" ]; then
    sed_inplace "s/Version != \"v[0-9]\.[0-9]\.[0-9]\"/Version != \"$NEW_VERSION\"/" "pkg/version/version_test.go"
    sed_inplace "s/Expected version v[0-9]\.[0-9]\.[0-9]/Expected version $NEW_VERSION/" "pkg/version/version_test.go"
    echo "Updated pkg/version/version_test.go"
fi

# Update README.md
if [ -f "README.md" ]; then
    sed_inplace "s|releases/download/v[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*/|releases/download/$NEW_VERSION/|g" "README.md"
    echo "Updated README.md"
fi

echo "Version sync complete: $NEW_VERSION"