#!/bin/bash

set -eo pipefail
set -x

TEMP_CONFIG_DIR=/data/temp
CUSTOM_CONFIG_DIR=/data/custom
CONFIG_DIR=/data/hazelcast

cp -f -R -L "$TEMP_CONFIG_DIR"/* "$CONFIG_DIR"

for FILE_DIR in "$CONFIG_DIR"/*; do
  # extract file name
  FILE_NAME=$(basename -- "$FILE_DIR")

  # extract file extension
  EXTENSION="${FILE_NAME##*.}"

  if [[ "$EXTENSION" == "yaml" ]]; then
    # merge user provided custom config with the updated one
    if [ -f $CUSTOM_CONFIG_DIR/"$FILE_NAME" ]; then
        NETWORK="hazelcast.network"
        ADVANCED_NETWORK="advanced-network"
        if yq r "$FILE_DIR" hazelcast.network | grep -q "true" && yq r "$CUSTOM_CONFIG_DIR"/$FILE_NAME  hazelcast.advanced-network | grep -q "true"; then
           yq delete -i "$FILE_DIR" "$NETWORK"
           echo "Key '$NETWORK' deleted from '$FILE_DIR' successfully."
        fi
        yq merge -i --overwrite "$FILE_DIR" "$CUSTOM_CONFIG_DIR"/"$FILE_NAME"
    fi
  fi
done