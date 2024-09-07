#!/bin/sh

echo -e "======================写入rclone配置========================\n"
echo "$RCLONE_CONF" > ~/.config/rclone/rclone.conf


code-server --bind-addr 0.0.0.0:7860


