#!/bin/bash

set -eo pipefail
set -x

TEMP_CONFIG_DIR=/data/temp
CUSTOM_CONFIG_DIR=/data/custom
CONFIG_DIR=/data/hazelcast

whoami

echo "WHAT THE FUCKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK"

if [ -f "$TEMP_CONFIG_DIR"/hazelcast.yaml ]; then
  echo "file is copied in temp"
else
  echo "file is not copieddddddddd in temp"
fi

echo "$TEMP_CONFIG_DIR"/*

echo $CONFIG_DIR

echo "doneeeeeeeeeeeeeeeeeeeeeeeeee"

ls -l /data/hazelcast
echo "doneeeeeeeeeeeeeeeeeeeeeeeeeekkkkkkkkkkkkkkkkkkkkkkkk"


cp -f -R -L "$TEMP_CONFIG_DIR"/* "$CONFIG_DIR"

ls -l /data

echo "done and dusted"

ls -l /data/temp

echo "doneeeeeeeeeeeeeeeeeeeeeeeeee"

ls -l /data/hazelcast

#sleep 10


echo "$CONFIG_DIR"/*


if [ -e "$CONFIG_DIR"/hazelcast.yaml ]; then
  echo "file is copied"
else
  echo "file is not copieddddddddd"
fi

for FILE_DIR in "$CONFIG_DIR"/*; do
  # extract file name
  FILE_NAME=$(basename -- "$FILE_DIR")

  # extract file extension
  EXTENSION="${FILE_NAME##*.}"

  if [ -f "$CONFIG_DIR"/"$FILE_NAME" ]; then
    echo "file is copied in loop"
  else
    echo "file is not copieddddddddd"
  fi

  if [[ "$EXTENSION" == "yaml" ]]; then
    # merge user provided custom config with the updated one
    if [ -f $CUSTOM_CONFIG_DIR/"$FILE_NAME" ]; then
        NETWORK="hazelcast.network"
        ADVANCED_NETWORK="advanced-network"
#        if yq e "has(\"$NETWORK\")" "$FILE_DIR" | grep -q "true" && yq e "has(\"$ADVANCED_NETWORK\")" "$CUSTOM_CONFIG_DIR"/$FILE_NAME | grep -q "true"; then
         yq delete -i "$FILE_DIR" "$NETWORK"
         echo "Key '$NETWORK' deleted from '$FILE_DIR' successfully.^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^"
         cat "$FILE_DIR"
         echo "done with config^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^"
#        fi
        yq merge -i --overwrite "$FILE_DIR" "$CUSTOM_CONFIG_DIR"/"$FILE_NAME"
    fi
  fi
done