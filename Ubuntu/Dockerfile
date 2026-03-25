FROM linkease/desktop-ubuntu2-standard-amd64:latest
RUN sudo apt update
# Uninstall Firefox
RUN sudo apt purge firefox -y
RUN sudo apt autoremove -y
RUN sudo apt install -y vim wget python3-pip python-is-python3
# Install Lapce
RUN wget https://github.com/lapce/lapce/releases/download/v0.4.6/lapce.ubuntu.noble.amd64.deb
RUN sudo apt install -yf ./lapce.ubuntu.noble.amd64.deb
RUN rm lapce.ubuntu.noble.amd64.deb
# Install Chrome
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN sudo apt install -yf ./google-chrome-stable_current_amd64.deb 
RUN rm google-chrome-stable_current_amd64.deb
RUN sudo sed -i 's|Exec=.*google-chrome.*|Exec=/usr/bin/google-chrome --no-sandbox --disable-dev-shm-usage %U|' /usr/share/applications/google-chrome*.desktop
# Install Playwright
RUN python -m pip install playwright --break-system-packages
#RUN python -m playwright install chrome --with-deps
RUN python -m playwright install-deps
