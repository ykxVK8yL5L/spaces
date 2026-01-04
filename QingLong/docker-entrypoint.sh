#!/bin/bash

export PATH="$HOME/bin:$PATH"

dir_shell=/ql/shell
. $dir_shell/share.sh


echo -e "======================写入rclone配置========================\n"
echo "$RCLONE_CONF" > ~/.config/rclone/rclone.conf

export_ql_envs() {
  export BACK_PORT="${ql_port}"
  export GRPC_PORT="${ql_grpc_port}"
}

log_with_style() {
  local level="$1"
  local message="$2"
  local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
  printf "\n[%s] [%7s]  %s\n" "${timestamp}" "${level}" "${message}"
}



# Fix DNS resolution issues in Alpine Linux
# Alpine uses musl libc which has known DNS resolver issues with certain domains
# Adding ndots:0 prevents unnecessary search domain appending
if [ -f /etc/alpine-release ]; then
  if ! grep -q "^options ndots:0" /etc/resolv.conf 2>/dev/null; then
    echo "options ndots:0" >> /etc/resolv.conf
    log_with_style "INFO" "🔧  0. 已配置 DNS 解析优化 (ndots:0)"
  fi
fi

log_with_style "INFO" "🚀  1. 检测配置文件..."
load_ql_envs
export_ql_envs
. $dir_shell/env.sh
import_config "$@"
fix_config

# Try to initialize PM2, but don't fail if it doesn't work
pm2 l &>/dev/null || log_with_style "WARN" "PM2 初始化可能失败，将在启动时尝试使用备用方案"

log_with_style "INFO" "⚙️  2. 启动 pm2 服务..."
reload_pm2

if [[ $AutoStartBot == true ]]; then
  log_with_style "INFO" "🤖  3. 启动 bot..."
  nohup ql bot >$dir_log/bot.log 2>&1 &
fi

if [[ $EnableExtraShell == true ]]; then
  log_with_style "INFO" "🛠️  4. 执行自定义脚本..."
  nohup ql extra >$dir_log/extra.log 2>&1 &
fi

log_with_style "SUCCESS" "🎉  容器启动成功!"


echo -e "======================3. 启动nginx========================\n"
nginx -s reload 2>/dev/null || nginx -c /etc/nginx/nginx.conf
echo -e "nginx启动成功...\n"

echo -e "##########写入登陆信息############"
#echo "{ \"username\": \"$ADMIN_USERNAME\", \"password\": \"$ADMIN_PASSWORD\" }" > /ql/data/config/auth.json
dir_root=/ql && source /ql/shell/api.sh 
init_auth_info() {
  local body="$1"
  local tip="$2"
  local currentTimeStamp=$(date +%s)
  local api=$(
    curl -s --noproxy "*" "http://0.0.0.0:5700/api/user/init?t=$currentTimeStamp" \
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
  # REMOTE_FOLDER="huggingface:/qinglong"

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
    dir_root=/ql && sleep 60 && source /ql/shell/api.sh && notify_api '青龙服务启动通知' '青龙面板成功启动'
else
    echo "没有检测到通知配置信息，不进行通知"
fi

#pm2 start code-server --name "code-server" -- --bind-addr 0.0.0.0:7860 --port 7860
# export PASSWORD=$ADMIN_PASSWORD
# code-server --bind-addr 0.0.0.0:7860 --port 7860

export PASSWORD=$ADMIN_PASSWORD
pm2 start "code-server --bind-addr 0.0.0.0:7860 --port 7860" --name "code-server"
pm2 startup
pm2 save


tail -f /dev/null

exec "$@"
