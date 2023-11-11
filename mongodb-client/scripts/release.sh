#!/bin/bash
read -p "Enter the version: " VERSION

# Create a Git tag
git tag -a "v$VERSION" -m "Release $VERSION"
git push origin v$VERSION
