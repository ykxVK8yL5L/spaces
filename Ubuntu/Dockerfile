FROM linkease/desktop-ubuntu2-standard-amd64:latest
RUN sudo apt update
RUN sudo apt install -y snapd vim python3-pip python-is-python3
RUN sudo snap install lapce --edge
# Install Chrome
RUN curl https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb -o /chrome.deb
RUN dpkg -i /chrome.deb || apt-get install -yf
RUN rm /chrome.deb
RUN python -m pip install playwright --break-system-packages
RUN python -m playwright install-deps
RUN python -m playwright install chrome
