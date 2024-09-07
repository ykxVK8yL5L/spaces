# Use Python 3.9 as the base image
FROM python:3.9

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages
RUN apt-get update && \
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
    rclone 


# Install Node.js (LTS version)
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash - && \
    apt-get install -y nodejs

RUN npm install -g pnpm@8.3.1 pm2 ts-node 


# Install code-server
RUN curl -fsSL https://code-server.dev/install.sh | sh -s -- --version=4.23.0-rc.2

# Install ollama
RUN curl -fsSL https://ollama.com/install.sh | sh

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
RUN sudo -u coder code-server --install-extension ms-python.python

# Switch to the coder user for running code-server
USER coder

ENV HOME=/home/coder \
	PATH=/home/coder/.local/bin:$PATH

COPY --chown=coder start_server.sh $HOME
RUN chmod +x $HOME/start_server.sh

WORKDIR /home/coder

# 创建rclone配置文件
RUN rclone config -h

# Start code-server with authentication
# CMD ["sh", "-c", "code-server --bind-addr 0.0.0.0:7860"]

CMD ["sh", "-c", "/home/coder/start_server.sh"]
# ENTRYPOINT ["/home/coder/start_server.sh"]

# Expose the default code-server port
EXPOSE 8080

# End of Dockerfile