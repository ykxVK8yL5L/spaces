#!/bin/sh

if [ -n "$RCLONE_CONF" ]; then
    rclone delete $REMOTE_FOLDER  #删除之前的备份，如需要保留之前的备份可以注释掉
    rclone sync /app/conf $REMOTE_FOLDER 
else
    echo "没有检测到Rclone配置信息" 
fi