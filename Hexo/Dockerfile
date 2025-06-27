FROM node:22-slim as nodebuilder
FROM rclone/rclone:latest AS rclone

FROM python:3.13-slim-bullseye

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ARG RCLONE_CONF=$RCLONE_CONF


COPY --from=nodebuilder /usr/local/bin/node /usr/local/bin/
COPY --from=nodebuilder /usr/local/lib/node_modules/. /usr/local/lib/node_modules/
COPY --from=rclone /usr/local/bin/rclone /usr/bin/rclone

# Install necessary packages
RUN ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && \
    apt-get update && \
    apt-get install -y \
    curl \
    sudo \
    build-essential \
    default-jdk \
    default-jre \
    g++ \
    gcc \
    libzbar0 \
    fish \
    ffmpeg \
    nmap \
    ca-certificates \
    zsh \
    cron \
    wget \
    tzdata \
    perl \
    openssl \
    openssh-client \
    nginx \
    jq \
    procps \
    netcat \
    sshpass \
    unzip \
    git \
    make \
    chromium \
    libatomic1 

RUN npm install -g pnpm pm2 hexo-cli

RUN chmod -R 777 /var && \
    chmod -R 777 /usr/local && \
    chmod -R 777 /etc/nginx && \
    chmod -R 777 /run && \
    chmod -R 777 /usr && \
    chmod -R 777 /root 

# Install code-server
RUN curl -fsSL https://code-server.dev/install.sh | sh -s -- --version=4.23.0-rc.2


# Create a user to run code-server
RUN useradd -m -s /bin/zsh coder && \
    echo 'coder ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Create code-server configuration directory
RUN mkdir -p /home/coder/.local/share/code-server/User
RUN chmod -R 777 /home/coder




# Add settings.json to enable dark mode
RUN echo '{ \
   "workbench.colorTheme": "Default Dark Modern", \
    "telemetry.enableTelemetry": true, \
    "telemetry.enableCrashReporter": true \
}' > /home/coder/.local/share/code-server/User/settings.json

# Change ownership of the configuration directory
RUN chown -R coder:coder /home/coder/.local/share/code-server

# Install Python extension for code-server
# RUN sudo -u coder code-server --install-extension ms-python.python

# Switch to the coder user for running code-server
USER coder

ENV HOME=/home/coder \
    PATH=/home/coder/.local/bin:$PATH

COPY --chown=coder start_server.sh $HOME
COPY --chown=coder apps.conf $HOME
COPY --chown=coder nginx.conf $HOME

RUN chmod +x $HOME/start_server.sh

WORKDIR /home/coder
RUN hexo init blog && cd blog && pnpm i

# 创建rclone配置文件
RUN rclone config -h

# Start code-server with authentication
# CMD ["sh", "-c", "code-server --bind-addr 0.0.0.0:7860"]

CMD ["sh", "-c", "/home/coder/start_server.sh"]
# ENTRYPOINT ["/home/coder/start_server.sh"]

# Expose the default code-server port
EXPOSE 5700

# End of Dockerfile
