#!/bin/bash
set -e

# 检查 /QLToolsV2/config 目录是否存在，不存在则创建
if [ ! -d /QLToolsV2/config ]; then
  printf "检测到config配置目录不存在，正在创建...\n"
  mkdir -p /QLToolsV2/config
fi

if [ ! -s /QLToolsV2/config/config.yaml ]; then
  printf "检测到config配置目录下不存在config.yaml，从示例文件复制一份用于初始化...\n"
  cp -fv /QLToolsV2/example.config.yaml /QLToolsV2/config/config.yaml
fi

if [ -s /QLToolsV2/config/config.yaml ]; then
  printf "检测到config配置目录下存在config.yaml，即将启动...\n"

  ./QLToolsV2-linux-"${TARGET_ARCH}"

fi

exec "$@"
