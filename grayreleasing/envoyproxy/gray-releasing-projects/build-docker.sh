#!/bin/bash

set -o errexit

if [ "$#" -ne 1 ]; then
  echo "Incorrect parameters"
  echo "Usage: build-docker.sh <version>"
  exit 1
fi

VERSION=$1
echo "version is ${VERSION}"
PREFIX="docker.io/polarismesh"
SCRIPTDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

pushd "$SCRIPTDIR/"
#java build the app.
docker run --rm -u root -v "$(pwd)":/home/maven/project -w /home/maven/project maven:3.8.1-openjdk-8-slim mvn clean package

pushd decorate-service
TAG_NAME="${PREFIX}/examples-gray-releasing-decorate-v1:${VERSION}"
docker build --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="1.0.0" .
docker push "${TAG_NAME}"

TAG_NAME="${PREFIX}/examples-gray-releasing-decorate-v2:${VERSION}"
docker build --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="2.0.0" .
docker push "${TAG_NAME}"
popd

pushd user-service
TAG_NAME="${PREFIX}/examples-gray-releasing-user-v1:${VERSION}"
docker build --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="1.0.0" .
docker push "${TAG_NAME}"

TAG_NAME="${PREFIX}/examples-gray-releasing-user-v2:${VERSION}"
docker build --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="2.0.0" .
docker push "${TAG_NAME}"
popd

pushd credit-service
TAG_NAME="${PREFIX}/examples-gray-releasing-credit-v1:${VERSION}"
docker build --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="1.0.0" .
docker push "${TAG_NAME}"
TAG_NAME="${PREFIX}/examples-gray-releasing-credit-v2:${VERSION}"
docker build --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="2.0.0" .
docker push "${TAG_NAME}"
popd

pushd promotion-service
TAG_NAME="${PREFIX}/examples-gray-releasing-promotion-v1:${VERSION}"
docker build --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="1.0.0" .
docker push "${TAG_NAME}"
TAG_NAME="${PREFIX}/examples-gray-releasing-promotion-v2:${VERSION}"
docker build --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="2.0.0" .
docker push "${TAG_NAME}"
popd

popd
