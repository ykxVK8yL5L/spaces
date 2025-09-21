#!/bin/bash

echo -e "======================写入rclone配置========================\n"
echo "$RCLONE_CONF" > ~/.config/rclone/rclone.conf

if [ -n "$RCLONE_CONF" ]; then
  echo -e "##########恢复n8n备份############"
  # 指定远程文件夹路径，格式为 remote:path
  REMOTE_FOLDER="huggingface:/n8n"

  # 使用 rclone ls 命令列出文件夹内容，将输出和错误分别捕获
  OUTPUT=$(rclone ls "$REMOTE_FOLDER" 2>&1)

  # 获取 rclone 命令的退出状态码
  EXIT_CODE=$?

  # 判断退出状态码
  if [ $EXIT_CODE -eq 0 ]; then
    # rclone 命令成功执行，检查文件夹是否为空
    if [ -z "$OUTPUT" ]; then
      #为空不处理
      #rclone sync --interactive /ql $REMOTE_FOLDER
      echo "初次安装"
    else
      #echo "文件夹不为空"
      rclone sync $REMOTE_FOLDER ~/.n8n 
    fi
  elif [[ "$OUTPUT" == *"directory not found"* ]]; then
    echo "错误：文件夹不存在"
  else
    echo "错误：$OUTPUT"
  fi
else
    echo "没有检测到Rclone配置信息"
fi


ln -s ~/.n8n ~/n8n-data
mkdir /etc/nginx/conf.d
mkdir /run/nginx
cp -fv ~/apps.conf /etc/nginx/conf.d
cp -fv ~/nginx.conf /etc/nginx/nginx.conf

echo -e "======================启动nginx========================\n"
nginx -s reload 2>/dev/null || nginx -c /etc/nginx/nginx.conf
echo -e "nginx启动成功...\n"

echo -e "======================启动pm2服务========================\n"
pm2 start n8n
export PASSWORD=$ADMIN_PASSWORD
pm2 start "code-server --bind-addr 0.0.0.0:7860 --port 7860" --name "code-server"
pm2 startup
pm2 save

# export PASSWORD=$ADMIN_PASSWORD
# code-server --bind-addr 0.0.0.0:7860 --port 7860

tail -f /dev/null

exec "$@"

