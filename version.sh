#!/bin/sh

VERSION=$(
  jq -r '.version' version.json
)

echo "version=$VERSION"