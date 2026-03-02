FROM ghcr.io/engigu/message-nest:latest AS builder

FROM node:20
ENV TZ=Asia/Shanghai
ENV REMOTE_FOLDER="huggingface:/messagenest"

#安装依赖
RUN apt update \
        && apt install -y ca-certificates rclone jq cron \
        && update-ca-certificates \
        && rm -rf /var/lib/apt/lists/*

RUN npm install -g pm2
#Rclone配置
RUN rclone config -h
#复制文件并给权限
COPY --from=builder /app/ /app/
RUN chmod -R 777 /app
COPY ./backup.sh /
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /backup.sh
RUN chmod +x /entrypoint.sh
#创建定时备份任务
RUN echo "3 2 * * * root /backup.sh" > /etc/cron.d/backup
RUN chmod 0644 /etc/cron.d/backup
RUN crontab /etc/cron.d/backup

WORKDIR /app
EXPOSE 8000

ENTRYPOINT ["/entrypoint.sh"]



