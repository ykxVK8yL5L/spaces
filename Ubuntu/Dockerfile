FROM linkease/desktop-ubuntu2-standard-amd64:latest
RUN sudo apt update
RUN sudo apt install -y vim python3-pip python-is-python3
RUN python -m pip install playwright --break-system-packages
RUN python -m playwright install-deps
RUN python -m playwright install chrome