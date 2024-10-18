#!/bin/bash
set -e

if [ ! -s /QLToolsV2/config/config.yaml ]; then
  printf "检测到config配置目录下不存在config.yaml，从示例文件复制一份用于初始化...\n"
  cp -fv /QLToolsV2/config/example.config.yaml /QLToolsV2/config/config.yaml
fi

if [ -s /QLToolsV2/config/config.yaml ]; then
  printf "检测到config配置目录下存在config.yaml，即将启动...\n"

  ./QLToolsV2-linux-"${TARGET_ARCH}"

fi

exec "$@"
