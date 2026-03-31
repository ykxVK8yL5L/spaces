#!/bin/bash

set -e

# 1. 补全目录
mkdir -p /root/.openclaw/agents/main/sessions
mkdir -p /root/.openclaw/credentials
mkdir -p /root/.openclaw/sessions


# 2. 执行恢复
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
        echo "远程文件夹不为空开始还原"
        ./sync.sh restore
        echo "恢复完成."   
    fi
  elif [[ "$OUTPUT" == *"directory not found"* ]]; then
    echo "错误：文件夹不存在"
  else
    echo "错误：$OUTPUT"
  fi
else
    echo "没有检测到Rclone配置信息"
fi


# ── 3. Fix DNS ────────────────────────────────────────────────
echo "nameserver 8.8.8.8" >> /etc/resolv.conf
echo "nameserver 8.8.4.4" >> /etc/resolv.conf
echo ">>> DNS fixed."

# ── 4. Chromium ───────────────────────────────────────────────
export PLAYWRIGHT_BROWSERS_PATH=/root/.openclaw/browsers
CHROMIUM_PATH=$(find /root/.openclaw/browsers -name "chrome" -type f 2>/dev/null | head -1)

if [ -z "$CHROMIUM_PATH" ]; then
    echo ">>> Installing Chromium..."
    OPENCLAW_NM=$(npm root -g 2>/dev/null)/openclaw/node_modules/playwright-core/cli.js
    if timeout 180 node "$OPENCLAW_NM" install chromium; then
        echo ">>> Chromium OK"
    else
        echo ">>> WARN: Chromium install failed"
    fi
    CHROMIUM_PATH=$(find /root/.openclaw/browsers -name "chrome" -type f 2>/dev/null | head -1)
else
    echo ">>> Chromium found: $CHROMIUM_PATH"
fi

# 3. 处理 API 地址
CLEAN_BASE=$(echo "$OPENAI_API_BASE" | sed "s|/chat/completions||g" | sed "s|/v1/|/v1|g" | sed "s|/v1$|/v1|g")

# 4. 生成配置文件
cat > /root/.openclaw/openclaw.json <<EOF
{
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "$CLEAN_BASE",
        "apiKey": "$OPENAI_API_KEY",
        "api": "openai-completions",
        "models": [
          { "id": "$MODEL", "name": "$MODEL", "contextWindow": 128000 }
        ]
      }
    }
  },
  "agents": { "defaults": { "model": { "primary": "nvidia/$MODEL" } } },
  "commands": {
    "restart": true
  },
  "gateway": {
    "mode": "local",
    "bind": "lan",
    "port": 7861,
    "trustedProxies": ["0.0.0.0/0"],
    "auth": { "mode": "token", "token": "$OPENCLAW_GATEWAY_PASSWORD" },
    "controlUi": {
      "enabled": true,
      "allowInsecureAuth": true,
      "allowedOrigins": ["*"],
      "dangerouslyDisableDeviceAuth": true,
      "dangerouslyAllowHostHeaderOriginFallback": true
    },
  }
}
EOF



cat > /etc/nginx/nginx.conf <<'EOF'
worker_processes 1;
events {
    worker_connections 1024;
}

http {
    map $http_upgrade $connection_upgrade {
        default upgrade;
        ''      close;
    }

    server {
        listen 7860;
        server_name _;
        
        location / {
            proxy_pass http://127.0.0.1:7861/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_set_header X-Forwarded-Prefix /openclaw/;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
            proxy_set_header X-Forwarded-Host $host;
        }
        
        location /telegram/webhook {
            proxy_pass http://127.0.0.1:8787;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

    }
}
EOF

nginx -t && nginx || { echo "nginx 配置失败"; cat /var/log/nginx/error.log; exit 1; }


# 5. 启动定时备份 (每 5 分钟)
# (while true; do
#   sleep 3600
#   echo ">>> Running scheduled backup..."
#   ./sync.sh backup
# done) &


# 6. 运行
openclaw doctor --fix

exec openclaw gateway run --port 7861
