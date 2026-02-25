FROM ghcr.io/engigu/baihu:latest

ENV REMOTE_FOLDER="huggingface:/baihu"

RUN apt-get update && apt-get install -y rclone jq
RUN rclone config -h
COPY ./restore.sh /
RUN sed -i '/^exec \.\/baihu$/d' /app/docker-entrypoint.sh
RUN cat /restore.sh >> /app/docker-entrypoint.sh
#RUN sed -i '50r /restore.sh' /app/docker-entrypoint.sh

