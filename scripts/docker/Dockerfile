FROM  imdt/bigbluebutton:2.4.x-develop

RUN unlink /etc/nginx/sites-enabled/default
COPY nginx.conf /etc/nginx/sites-available/bigbluebutton

# Telegraf installation
RUN curl -s https://repos.influxdata.com/influxdb.key | sudo apt-key add -
RUN . /etc/lsb-release
RUN echo "deb https://repos.influxdata.com/ubuntu bionic stable" | sudo tee /etc/apt/sources.list.d/influxdb.list

RUN sudo apt update
RUN sudo apt install -y telegraf
COPY ./telegraf.conf /etc/telegraf/telegraf.conf

ENTRYPOINT ["/bin/systemd", "--system", "--unit=multi-user.target"]