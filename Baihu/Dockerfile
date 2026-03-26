FROM ghcr.io/engigu/baihu:latest

ENV REMOTE_FOLDER="huggingface:/baihu"

RUN apt-get update && apt-get install -y rclone jq
RUN rclone config -h
COPY ./restore.sh /
RUN sed -i '/^exec \.\/baihu server$/d' /app/docker-entrypoint.sh
RUN cat /restore.sh >> /app/docker-entrypoint.sh
#RUN sed -i '50r /restore.sh' /app/docker-entrypoint.sh
RUN mise exec python@3.13.12 -- pip install requests
RUN mise use -g python@3.13.12
RUN mise use -g node@23.11.1
