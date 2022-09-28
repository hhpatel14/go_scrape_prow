FROM quay.io/centos/centos:stream8

RUN dnf clean all && dnf update -y  && dnf install -y \
                                        golang \
                                        tmux


COPY influxdb.repo /etc/yum.repos.d/influxdb.repo

RUN dnf -y install telegraf wget unzip procps-ng 
RUN dnf -y install https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm
RUN wget -N https://chromedriver.storage.googleapis.com/97.0.4692.36/chromedriver_linux64.zip -P ~/
RUN unzip ~/chromedriver_linux64.zip -d ~/
RUN mv -f ~/chromedriver /usr/local/bin/chromedriver
RUN chown root:root /usr/local/bin/chromedriver
RUN chmod 0755 /usr/local/bin/chromedriver
ENV DISPLAY=:99

# update ulimits, in /etc/security/limits.conf
COPY limits.conf /etc/security/limits.conf

# install golang go_scrape_prow
COPY src/go.mod ./
COPY src/go.sum ./
RUN go mod download
COPY src/*.go ./
RUN go build -o /prow
RUN cp /prow /usr/bin/prow

# setup telegraf
COPY telegraf.conf /etc/telegraf/telegraf.conf
COPY telegraf.d/* /etc/telegraf/telegraf.d/
# RUN systemctl enable --now telegraf



#ENTRYPOINT ["/entrypoint.sh"]
CMD ["telegraf", "--debug", "--config-directory","/etc/telegraf/telegraf.d/","--input-filter", "exec"]
