#!/bin/sh

VERSION=$(
  jq -r '.version' package.json
)

echo "version=$VERSION"