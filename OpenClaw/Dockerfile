FROM node:22-slim

ENV REMOTE_FOLDER="huggingface:/openclaw"
# 1. 基础依赖
RUN apt-get update && apt-get install -y --no-install-recommends \
    git openssh-client build-essential python3 python3-pip \ 
    g++ rclone nginx make ca-certificates && rm -rf /var/lib/apt/lists/*

RUN pip3 install --no-cache-dir huggingface_hub --break-system-packages

# Install code-server
RUN curl -L https://coder.com/install.sh | sh

# 2、创建rclone配置文件夹
RUN rclone config -h

# 3. 安装 OpenClaw
RUN npm install -g pm2 openclaw@latest --unsafe-perm 

# 4. 设置工作目录并拷贝脚本
WORKDIR /app
COPY sync.sh .
COPY start-openclaw.sh .
RUN chmod +x start-openclaw.sh
RUN chmod +x sync.sh
RUN chmod -R 777 /root

# 5. 环境变量
ENV PORT=7860 HOME=/root


EXPOSE 7860
CMD ["./start-openclaw.sh"]
