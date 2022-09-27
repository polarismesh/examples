#  Tencent is pleased to support the open source community by making Polaris available.
#  Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
#  Licensed under the BSD 3-Clause License (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#  https://opensource.org/licenses/BSD-3-Clause
#  Unless required by applicable law or agreed to in writing, software distributed
#  under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
#  CONDITIONS OF ANY KIND, either express or implied. See the License for the
#  specific language governing permissions and limitations under the License.

#!/bin/bash

set -o errexit

if [ "$#" -ne 1 ]; then
  echo "Incorrect parameters"
  echo "Usage: build_docker.sh <version>"
  exit 1
fi

VERSION=$1
echo "version is ${VERSION}"
PREFIX="docker.io/polarismesh"
SCRIPTDIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

pushd "$SCRIPTDIR/"
#java build the app.
mvn clean package

pushd user-service
TAG_NAME="${PREFIX}/sct-gray-releasing-user-v1:${VERSION}"
docker build --network=host --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="1.0.0" .
docker push "${TAG_NAME}"

TAG_NAME="${PREFIX}/sct-gray-releasing-user-v2:${VERSION}"
docker build --network=host --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="2.0.0" .
docker push "${TAG_NAME}"
popd

pushd credit-service
TAG_NAME="${PREFIX}/sct-gray-releasing-credit-v1:${VERSION}"
docker build --network=host --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="1.0.0" .
docker push "${TAG_NAME}"
TAG_NAME="${PREFIX}/sct-gray-releasing-credit-v2:${VERSION}"
docker build --network=host --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="2.0.0" .
docker push "${TAG_NAME}"
popd

pushd promotion-service
TAG_NAME="${PREFIX}/sct-gray-releasing-promotion-v1:${VERSION}"
docker build --network=host --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="1.0.0" .
docker push "${TAG_NAME}"
TAG_NAME="${PREFIX}/sct-gray-releasing-promotion-v2:${VERSION}"
docker build --network=host --pull -t "${TAG_NAME}" --build-arg pkg_version="${VERSION}" --build-arg logic_version="2.0.0" .
docker push "${TAG_NAME}"
popd

popd