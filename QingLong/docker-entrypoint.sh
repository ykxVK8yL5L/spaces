#!/bin/bash

dir_shell=/ql/shell
. $dir_shell/share.sh
. $dir_shell/env.sh

echo -e "======================写入rclone配置========================\n"
echo "$RCLONE_CONF" > ~/.config/rclone/rclone.conf

echo -e "======================1. 检测配置文件========================\n"
import_config "$@"
make_dir /etc/nginx/conf.d
make_dir /run/nginx
init_nginx
fix_config

pm2 l &>/dev/null

echo -e "======================2. 安装依赖========================\n"
patch_version


echo -e "======================3. 启动nginx========================\n"
nginx -s reload 2>/dev/null || nginx -c /etc/nginx/nginx.conf
echo -e "nginx启动成功...\n"

echo -e "======================4. 启动pm2服务========================\n"
reload_update
reload_pm2

if [[ $AutoStartBot == true ]]; then
  echo -e "======================5. 启动bot========================\n"
  nohup ql bot >$dir_log/bot.log 2>&1 &
  echo -e "bot后台启动中...\n"
fi

if [[ $EnableExtraShell == true ]]; then
  echo -e "====================6. 执行自定义脚本========================\n"
  nohup ql extra >$dir_log/extra.log 2>&1 &
  echo -e "自定义脚本后台执行中...\n"
fi


echo -e "############################################################\n"
echo -e "容器启动成功..."
echo -e "############################################################\n"


echo -e "##########写入登陆信息############"
#echo "{ \"username\": \"$ADMIN_USERNAME\", \"password\": \"$ADMIN_PASSWORD\" }" > /ql/data/config/auth.json
dir_root=/ql && source /ql/shell/api.sh 
init_auth_info() {
  local body="$1"
  local tip="$2"
  local currentTimeStamp=$(date +%s)
  local api=$(
    curl -s --noproxy "*" "http://0.0.0.0:5600/api/user/init?t=$currentTimeStamp" \
      -X 'PUT' \
      -H "Accept: application/json" \
      -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36" \
      -H "Content-Type: application/json;charset=UTF-8" \
      -H "Origin: http://0.0.0.0:5700" \
      -H "Referer: http://0.0.0.0:5700/crontab" \
      -H "Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7" \
      --data-raw "{$body}" \
      --compressed
  )
  code=$(echo "$api" | jq -r .code)
  message=$(echo "$api" | jq -r .message)
  if [[ $code == 200 ]]; then
    echo -e "${tip}成功🎉"
  else
    echo -e "${tip}失败(${message})"
  fi
}

init_auth_info "\"username\": \"$ADMIN_USERNAME\", \"password\": \"$ADMIN_PASSWORD\"" "Change Password"

if [ -n "$RCLONE_CONF" ]; then
  echo -e "##########同步备份############"
  # 指定远程文件夹路径，格式为 remote:path
  REMOTE_FOLDER="huggingface:/qinglong"

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
      mkdir /ql/.tmp/data
      rclone sync $REMOTE_FOLDER /ql/.tmp/data && real_time=true ql reload data
    fi
  elif [[ "$OUTPUT" == *"directory not found"* ]]; then
    echo "错误：文件夹不存在"
  else
    echo "错误：$OUTPUT"
  fi
else
    echo "没有检测到Rclone配置信息"
fi

if [ -n "$NOTIFY_CONFIG" ]; then
    python /notify.py
    dir_root=/ql && source /ql/shell/api.sh && notify_api '青龙服务启动通知' '青龙面板成功启动'
else
    echo "没有检测到通知配置信息，不进行通知"
fi

tail -f /dev/null

exec "$@"
