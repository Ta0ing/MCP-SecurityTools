#!/bin/bash

# 创建输出目录
mkdir -p build

# 版本号
VERSION="1.0.0"

# 编译函数
build() {
    local os=$1
    local arch=$2
    local extension=""
    
    # Windows 需要 .exe 扩展名
    if [ "$os" == "windows" ]; then
        extension=".exe"
    fi
    
    echo "Building for $os/$arch..."
    GOOS=$os GOARCH=$arch go build -o "build/chaitin-mcp-${os}-${arch}${extension}" main.go
}

# Linux
build "linux" "amd64"
build "linux" "arm64"

# macOS
build "darwin" "amd64"
build "darwin" "arm64"

# Windows
build "windows" "amd64"
build "windows" "arm64"

echo "Build complete! Check the build directory for binaries."
