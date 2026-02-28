
npm install pm2 -g
#开启白虎服务
pm2 start "./baihu server" --name baihu

echo "10秒后开始恢复任务..."
sleep 10

#从日志中获取密码
echo  "======================从日志中获取密码========================\n"
DEFAULT_PASSWORD=$(tail -n 100 ~/.pm2/logs/baihu-out.log \
    | grep -oP '密\s*码:\s*\K[^,[:space:]]+' \
    | tail -n 1)

echo  "默认用户名: admin"
#echo  "默认密码: $DEFAULT_PASSWORD"

echo "============重置密码==============="
# 获取登陆响应Token
BHToken=$(
curl -c cookies.txt -s -D - -o /dev/null \
  'http://localhost:8052/api/v1/auth/login' \
  -H 'content-type: application/json' \
  --data-raw "{\"username\":\"admin\",\"password\":\"$DEFAULT_PASSWORD\"}" \
| awk -F'[=;]' '/Set-Cookie: BHToken=/{print $2}'
)

sleep 1

RESET_RESPONSE=$(
  curl -b cookies.txt 'http://localhost:8052/api/v1/settings/password' \
  -H 'content-type: application/json' \
  --data-raw "{\"old_password\":\"$DEFAULT_PASSWORD\",\"new_password\":\"$ADMIN_PASSWORD\"}"
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
      # rclone sync $REMOTE_FOLDER /app --exclude="/baihu" --exclude "/docker-entrypoint.sh"
      mkdir /app/backup_tmp
      # 找最新的文件名
      latest_file=$(rclone lsjson $REMOTE_FOLDER | jq -r 'sort_by(.ModTime) | last | .Path')
      # 复制到目标目录
      rclone copy $REMOTE_FOLDER/$latest_file /app/backup_tmp
      RESTORE_RESPON=$(curl -b cookies.txt "http://localhost:8052/api/v1/settings/restore" \
        -F "file=@/app/backup_tmp/$latest_file;type=application/zip" \
        -H "Accept: */*" \
        --compressed
      )
      rm -rf /app/backup_tmp
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

