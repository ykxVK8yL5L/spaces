FROM python:3.9.5-slim-buster
LABEL author="Lan"
LABEL email="vast@tom.com"


RUN chmod -R 777 /var
RUN chmod 777 /etc/timezone

RUN set -x && \
  apt-get update && \
  apt-get install --no-install-recommends -y git && \
  git config --global user.email "filecodebox@@users.noreply.github.com" && \
  git config --global user.name "filecodebox" && \
  git config --global http.postBuffer 524288000 


RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone


RUN useradd -m -u 1000 user
USER user
ENV PATH="/home/user/.local/bin:$PATH"


RUN git clone https://github.com/vastsa/FileCodeBox /home/user/app

WORKDIR /home/user/app
RUN pip install -r requirements.txt
EXPOSE 12345
CMD ["python","main.py"]