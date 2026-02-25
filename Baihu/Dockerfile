FROM ghcr.io/engigu/baihu:latest

ENV REMOTE_FOLDER="huggingface:/baihu"

RUN apt-get update && apt-get install rclone
RUN rclone config -h
COPY ./restore.sh /
RUN sed -i '50r /restore.sh' /app/docker-entrypoint.sh

