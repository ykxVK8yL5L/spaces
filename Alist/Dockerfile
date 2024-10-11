FROM python:3.11-slim-bullseye
ADD https://github.com/ykxVK8yL5L/alist/releases/download/latest/alist-linux-amd64.tar.gz /usr/bin/

RUN cd /usr/bin && \
  tar -xzf alist-linux-amd64.tar.gz && \
  rm alist-linux-amd64.tar.gz

RUN chmod a+x /usr/bin/alist

RUN set -x && \
  apt-get update && \
  apt-get upgrade -y && \
  apt-get install --no-install-recommends -y git \
  rclone 

RUN chmod -R 777 /var/log

# Set up a new user named "user" with user ID 1000
RUN useradd -m -u 1000 user

# Switch to the "user" user
USER user

# Set home to the user's home directory
ENV HOME=/home/user \
  PATH=/home/user/.local/bin:$PATH

# Set the working directory to the user's home directory
WORKDIR $HOME/alist

RUN pip install yacron
COPY --chown=user yacrontab.yaml /etc/yacron.d/yacrontab.yaml

RUN rclone config -h

RUN --mount=type=secret,id=RCLONE_CONF,mode=0444,required=true \
  echo "$(cat /run/secrets/RCLONE_CONF)" > ~/.config/rclone/rclone.conf

COPY --chown=user start.sh /home/user/start.sh

RUN chmod a+x /home/user/start.sh

CMD [ "/home/user/start.sh" ]

EXPOSE 5244
