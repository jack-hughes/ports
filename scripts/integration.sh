#!/bin/bash
echo "--- [running integration tests] ---"

# Check for healthy response data
mkdir -p ../artifacts
curl --location --request GET 'localhost:8181/ports/INPAV' -o artifacts/tmp.json
if diff -u "./artifacts/tmp.json" "./test/testdata/get.json"; then
  echo "--- [test and fetched data have identical contents, passing] ---"
else
  : # this lists differences
  exit 1
fi

# Check the number of items we expect are returned
curl --location --request GET 'localhost:8181/ports' -o artifacts/tmp.json
entries=$(jq length artifacts/tmp.json)
if [ $entries -ne 1632 ]; then
  echo "--- [expected 1632 entries, got $entries] ---"
  rm artifacts/tmp.json
else
  echo "--- [expected entries met, passing] ---"
  rm artifacts/tmp.json
  exit 0
fi
