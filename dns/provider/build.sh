#!/bin/bash

set -euo pipefail

# 设置Docker仓库，默认为polarismesh
DOCKER_IMAGE_RPO="${DOCKER_RPO:-polarismesh}"
APP=polaris-dns-provider
push_image=true  # 默认推送镜像

# 解析参数
while [[ $# -gt 0 ]]; do
    case "$1" in
        --no-push)
            push_image=false
            shift
            ;;
        *)
            docker_tag="$1"
            shift
            ;;
    esac
done

# 参数检查
if [ -z "${docker_tag:-}" ]; then
    echo "Usage: bash $0 [--no-push] <docker_tag>"
    echo "e.g.: bash $0 v1.0"
    echo "e.g.: bash $0 --no-push v1.0-rc1"
    exit 1
fi

echo "Building Docker image: ${DOCKER_IMAGE_RPO}/${APP}:${docker_tag}"
$push_image || echo "Skipping image push (--no-push specified)"

# 确保在脚本所在目录执行
script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)
cd "$script_dir"

# 构建二进制
go build -buildvcs=false -o provider . || {
    echo "Error: Failed to build ${APP} binary" >&2
    exit 1
}

# 构建Docker镜像
docker build --network=host -t "${DOCKER_IMAGE_RPO}/${APP}:${docker_tag}" . || {
    echo "Error: Docker build failed" >&2
    exit 1
}

# 推送镜像（如果未指定--no-push）
if $push_image; then
    docker push "${DOCKER_IMAGE_RPO}/${APP}:${docker_tag}" || {
        echo "Error: Docker push failed" >&2
        exit 1
    }

    # 如果是正式版本，则同时标记为latest并推送
    if [[ ! "$docker_tag" =~ alpha|beta|rc ]]; then
        echo "Tagging as latest version"
        docker tag "${DOCKER_IMAGE_RPO}/${APP}:${docker_tag}" "${DOCKER_IMAGE_RPO}/${APP}:latest"
        docker push "${DOCKER_IMAGE_RPO}/${APP}:latest" || {
            echo "Error: Failed to push latest tag" >&2
            exit 1
        }
    fi
fi

echo "Docker build completed successfully"
if $push_image; then
    echo "Image pushed to registry"
fi
exit 0  # 确保返回成功状态
