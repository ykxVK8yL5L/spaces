#!/bin/sh

#开启定时任务
cron

npm install pm2 -g
#开启白虎服务
pm2 start "/app/Message-Nest" --name message

echo "10秒后开始恢复任务..."
sleep 10

#从日志中获取密码
echo  "======================从日志中获取密码========================\n"
DEFAULT_PASSWORD=$(tail -n 100 ~/.pm2/logs/message-out.log \
    | grep -oP '(?<=密码：)\S+' \
    | tail -n 1)

echo  "默认用户名: admin"
echo  "默认密码: $DEFAULT_PASSWORD"

echo "============重置密码==============="
# 获取登陆响应Token
MNToken=$(
curl -c cookies.txt \
  'http://localhost:8000/auth' \
  -H 'content-type: application/json' \
  --data-raw "{\"username\":\"admin\",\"passwd\":\"$DEFAULT_PASSWORD\"}" \
  | jq -r '.data.token'
)
sleep 1

RESET_RESPONSE=$(
  curl -b cookies.txt "http://localhost:8000/api/v1/settings/setpasswd" \
  -H "content-type: application/json" \
  -H "m-token: $MNToken" \
  --data-raw "{\"old_passwd\":\"$DEFAULT_PASSWORD\",\"new_passwd\":\"$ADMIN_PASSWORD\"}"
)


echo  "======================写入rclone配置========================\n"
echo "$RCLONE_CONF" > ~/.config/rclone/rclone.conf

if [ -n "$RCLONE_CONF" ]; then
  echo "##########同步备份############"

  # 使用 rclone ls 命令列出文件夹内容，将输出和错误分别捕获
  OUTPUT=$(rclone ls "$REMOTE_FOLDER" 2>&1)
  # 获取 rclone 命令的退出状态码
  EXIT_CODE=$?
  #echo "rclone退出代码:$EXIT_CODE"
  # 判断退出状态码
  if [ $EXIT_CODE -eq 0 ]; then
    # rclone 命令成功执行，检查文件夹是否为空
    if [ -z "$OUTPUT" ]; then
      #为空不处理
      echo "初次安装"
    else
      #echo "文件夹不为空"
      rclone sync $REMOTE_FOLDER /app/conf 
      # 找最新的文件名
      # latest_file=$(rclone lsjson $REMOTE_FOLDER | jq -r 'sort_by(.ModTime) | last | .Path')
      # 复制到目标目录
      # rclone copy $REMOTE_FOLDER/$latest_file /app/conf
      pm2 restart message
    fi
  elif [[ "$OUTPUT" == *"directory not found"* ]]; then
    echo "错误：文件夹不存在"
  else
    echo "错误：$OUTPUT"
  fi
else
    echo "没有检测到Rclone配置信息"
fi

tail -f /dev/null

