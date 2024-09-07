FROM node:20-slim as nodebuilder

FROM python:3.11-slim-bullseye as builder
ARG QL_MAINTAINER="whyour"
LABEL maintainer="${QL_MAINTAINER}"
ARG QL_URL=https://github.com/${QL_MAINTAINER}/qinglong.git
ARG QL_BRANCH=debian

ENV QL_DIR=/ql \
  QL_BRANCH=${QL_BRANCH}

COPY --from=nodebuilder /usr/local/bin/node /usr/local/bin/
COPY --from=nodebuilder /usr/local/lib/node_modules/. /usr/local/lib/node_modules/
RUN set -x && \
  ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && \
  apt-get update && \
  apt-get install --no-install-recommends -y libatomic1 git && \
  git config --global user.email "qinglong@@users.noreply.github.com" && \
  git config --global user.name "qinglong" && \
  git config --global http.postBuffer 524288000 && \
  git clone --depth=1 -b ${QL_BRANCH} ${QL_URL} ${QL_DIR} 

RUN mkdir /tmp/build
RUN cp ${QL_DIR}/package.json ${QL_DIR}/.npmrc ${QL_DIR}/pnpm-lock.yaml /tmp/build/

RUN npm i -g pnpm@8.3.1 && \
  cd /tmp/build && \
  pnpm install --prod

FROM python:3.11-slim-bullseye

ARG QL_MAINTAINER="whyour"
LABEL maintainer="${QL_MAINTAINER}"
ARG QL_URL=https://github.com/${QL_MAINTAINER}/qinglong.git
ARG QL_BRANCH=debian

ENV PNPM_HOME=/root/.local/share/pnpm \
  PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/.local/share/pnpm:/root/.local/share/pnpm/global/5/node_modules:$PNPM_HOME \
  NODE_PATH=/usr/local/bin:/usr/local/pnpm-global/5/node_modules:/usr/local/lib/node_modules:/root/.local/share/pnpm/global/5/node_modules \
  LANG=C.UTF-8 \
  SHELL=/bin/bash \
  PS1="\u@\h:\w \$ " \
  QL_DIR=/ql \
  QL_BRANCH=${QL_BRANCH}

COPY --from=nodebuilder /usr/local/bin/node /usr/local/bin/
COPY --from=nodebuilder /usr/local/lib/node_modules/. /usr/local/lib/node_modules/

RUN set -x && \
  ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && \
  ln -s /usr/local/lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx && \
  apt-get update && \
  apt-get upgrade -y && \
  apt-get install --no-install-recommends -y git \
  curl \
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
  rclone \
  unzip \
  libatomic1 && \
  apt-get clean && \
  ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
  echo "Asia/Shanghai" >/etc/timezone && \
  git config --global user.email "qinglong@@users.noreply.github.com" && \
  git config --global user.name "qinglong" && \
  git config --global http.postBuffer 524288000 && \
  npm install -g pnpm@8.3.1 pm2 ts-node && \
  rm -rf /root/.pnpm-store && \
  rm -rf /root/.local/share/pnpm/store && \
  rm -rf /root/.cache && \
  rm -rf /root/.npm && \
  chmod u+s /usr/sbin/cron && \
  ulimit -c 0

ARG SOURCE_COMMIT
RUN git clone --depth=1 -b ${QL_BRANCH} ${QL_URL} ${QL_DIR} && \
  cd ${QL_DIR} && \
  cp -f .env.example .env && \
  chmod 777 ${QL_DIR}/shell/*.sh && \
  chmod 777 ${QL_DIR}/docker/*.sh && \
  git clone --depth=1 -b ${QL_BRANCH} https://github.com/${QL_MAINTAINER}/qinglong-static.git /static && \
  mkdir -p ${QL_DIR}/static && \
  cp -rf /static/* ${QL_DIR}/static && \
  rm -rf /static && \
  rm -f ${QL_DIR}/docker/docker-entrypoint.sh

COPY docker-entrypoint.sh ${QL_DIR}/docker

RUN mkdir /ql/data && \
  mkdir /ql/data/config && \
  mkdir /ql/data/log && \
  mkdir /ql/data/db && \
  mkdir /ql/data/scripts && \
  mkdir /ql/data/repo && \
  mkdir /ql/data/raw && \
  mkdir /ql/data/deps && \
  chmod -R 777 /ql && \
  chmod -R 777 /var && \
  chmod -R 777 /usr/local && \
  chmod -R 777 /etc/nginx && \
  chmod -R 777 /run && \
  chmod -R 777 /usr && \
  chmod -R 777 /root 

COPY --from=builder /tmp/build/node_modules/. /ql/node_modules/

WORKDIR ${QL_DIR}


# Set up a new user named "user" with user ID 1000
RUN useradd -m -u 1000 user

# Switch to the "user" user
USER user

# 创建rclone配置文件
RUN rclone config -h

HEALTHCHECK --interval=5s --timeout=2s --retries=20 \
  CMD curl -sf --noproxy '*' http://127.0.0.1:5400/api/health || exit 1

ENTRYPOINT ["./docker/docker-entrypoint.sh"]

VOLUME /ql/data
  
EXPOSE 5700