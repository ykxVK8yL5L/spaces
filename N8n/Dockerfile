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
ARG DB_TYPE=$DB_TYPE
ARG DB_POSTGRESDB_SCHEMA=$DB_POSTGRESDB_SCHEMA
ARG DB_POSTGRESDB_HOST=$DB_POSTGRESDB_HOST
ARG DB_POSTGRESDB_DATABASE=$DB_POSTGRESDB_DATABASE
ARG DB_POSTGRESDB_PORT=$DB_POSTGRESDB_PORT
ARG DB_POSTGRESDB_USER=$DB_POSTGRESDB_USER
ARG DB_POSTGRESDB_PASSWORD=$DB_POSTGRESDB_PASSWORD

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
    postgresql-client

# Set environment variables
ENV PUPPETEER_SKIP_DOWNLOAD=true
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser


# 设置自定义的 npm global 安装路径
ENV NPM_CONFIG_PREFIX=/npm-global

# 创建一个共享的目录来安装全局包
RUN mkdir -p /npm-global && \
    chmod -R 777 /npm-global

# 将该目录加入到 PATH 环境变量中
ENV PATH=/npm-global/bin:$PATH

RUN npm config set prefix '/npm-global'
# Install n8n globally
RUN npm install -g n8n@1.86.1

# Create necessary directories
RUN mkdir -p $DATABASE_PATH $CONFIG_PATH $WORKFLOWS_PATH $LOGS_PATH \
    && chmod -R 777 $BASE_PATH

# Set working directory
WORKDIR /data


# Start n8n

CMD ["n8n", "start"]
