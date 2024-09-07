FROM ubuntu:22.04
ADD https://github.com/ykxVK8yL5L/alist/releases/download/latest/alist-linux-amd64.tar.gz /usr/bin/

RUN cd /usr/bin && \
  tar -xzf alist-linux-amd64.tar.gz && \
  rm alist-linux-amd64.tar.gz

RUN chmod a+x /usr/bin/alist

# Set up a new user named "user" with user ID 1000
RUN useradd -m -u 1000 user

# Switch to the "user" user
USER user

# Set home to the user's home directory
ENV HOME=/home/user \
	PATH=/home/user/.local/bin:$PATH

# Set the working directory to the user's home directory
WORKDIR $HOME/alist

CMD cd $HOME/alist && alist server

EXPOSE 5244