#!/bin/bash
#
# Copyright Istio Authors
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

set -o errexit

if [ "$#" -ne 1 ]; then
    echo "Incorrect parameters"
    echo "Usage: build-services.sh <version>"
    exit 1
fi

VERSION=$1
PREFIX="docker.io/polarismesh"
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

pushd "$SCRIPTDIR/reviews"
  #java build the app.
  docker run --rm -u root -v "$(pwd)":/home/maven/project -w /home/maven/project maven:3.8.1-openjdk-8-slim mvn clean package
  #plain build -- no ratings
  docker build --pull -t "${PREFIX}/examples-bookinfo-reviews-v1:${VERSION}" -t "${PREFIX}/examples-bookinfo-reviews-v1:latest" --build-arg service_version=v1 .
  docker push "${PREFIX}/examples-bookinfo-reviews-v1:${VERSION}"
  #with ratings black stars
  docker build --pull -t "${PREFIX}/examples-bookinfo-reviews-v2:${VERSION}" -t "${PREFIX}/examples-bookinfo-reviews-v2:latest" --build-arg service_version=v2 \
   --build-arg enable_ratings=true .
  docker push "${PREFIX}/examples-bookinfo-reviews-v2:${VERSION}"
  #with ratings red stars
  docker build --pull -t "${PREFIX}/examples-bookinfo-reviews-v3:${VERSION}" -t "${PREFIX}/examples-bookinfo-reviews-v3:latest" --build-arg service_version=v3 \
   --build-arg enable_ratings=true --build-arg star_color=red .
  docker push "${PREFIX}/examples-bookinfo-reviews-v3:${VERSION}"
popd
