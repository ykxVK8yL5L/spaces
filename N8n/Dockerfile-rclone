FROM node:18-alpine

# Set user to root for installation
USER root

# Arguments that can be passed at build time
ARG N8N_PATH=/usr/local/lib/node_modules/n8n
ARG BASE_PATH=/root/.n8n
ARG DATABASE_PATH=$BASE_PATH/database
ARG CONFIG_PATH=$BASE_PATH/config
ARG WORKFLOWS_PATH=$BASE_PATH/workflows
ARG LOGS_PATH=$BASE_PATH/logs
ARG N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS=$N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS
ARG N8N_HOST=$N8N_HOST
ARG N8N_PORT=$N8N_PORT
ARG N8N_PROTOCOL=https
ARG N8N_EDITOR_BASE_URL=$N8N_EDITOR_BASE_URL
ARG WEBHOOK_URL=$WEBHOOK_URL
ARG GENERIC_TIMEZONE=$GENERIC_TIMEZONE
ARG TZ=$TZ
ARG N8N_ENCRYPTION_KEY=$N8N_ENCRYPTION_KEY
ARG N8N_ENCRYPTION_KEY=$N8N_ENCRYPTION_KEY
ARG RCLONE_CONF=$RCLONE_CONF

# Install system dependencies
RUN apk add --no-cache \
    git \
    python3 \
    py3-pip \
    make \
    g++ \
    build-base \
    cairo-dev \
    pango-dev \
    chromium \
    rclone 

# Set environment variables
ENV PUPPETEER_SKIP_DOWNLOAD=true
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser



# Install n8n globally
RUN npm install -g n8n@1.86.1

# Create necessary directories
RUN mkdir -p $DATABASE_PATH $CONFIG_PATH $WORKFLOWS_PATH $LOGS_PATH \
    && chmod -R 777 $BASE_PATH

RUN mkdir -p /home/node/.config/rclone && chmod -R 777 /home/node/.config/rclone

COPY entrypoint.sh /entrypoint.sh

RUN chmod 777 /entrypoint.sh

# 创建rclone配置文件
RUN rclone config -h

# Set working directory
WORKDIR /data


CMD ["/entrypoint.sh"]

# Start n8n

# CMD ["n8n", "start"]
