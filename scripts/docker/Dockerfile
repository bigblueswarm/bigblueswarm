FROM ubuntu:bionic AS builder
WORKDIR /tmp
# Go installation
ENV GOPATH $HOME/go
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin
RUN apt update && apt install -y curl git
RUN curl -OL https://golang.org/dl/go1.17.6.linux-amd64.tar.gz
RUN tar -C /usr/local -xvf go1.17.6.linux-amd64.tar.gz
RUN go version

# build bbb telegraf plugin
RUN mkdir -p $GOPATH
WORKDIR $GOPATH
RUN git clone https://github.com/SLedunois/bigbluebutton-telegraf-plugin
WORKDIR $GOPATH/bigbluebutton-telegraf-plugin
RUN go build -o pkg/bbb-telegraf ./cmd/main.go

FROM  imdt/bigbluebutton:2.4.x-develop

RUN unlink /etc/nginx/sites-enabled/default
COPY nginx.conf /etc/nginx/sites-available/bigbluebutton

# Telegraf installation
RUN curl -s https://repos.influxdata.com/influxdb.key | apt-key add -
RUN . /etc/lsb-release
RUN echo "deb https://repos.influxdata.com/ubuntu bionic stable" | tee /etc/apt/sources.list.d/influxdb.list

RUN apt update
RUN apt install -y telegraf
RUN mkdir -p /etc/telegraf/plugins/bigbluebutton
COPY --from=builder /go/bigbluebutton-telegraf-plugin/pkg/bbb-telegraf /etc/telegraf/plugins/bigbluebutton/bbb-telegraf
RUN chmod +x /etc/telegraf/plugins/bigbluebutton/bbb-telegraf
COPY ./telegraf.conf /etc/telegraf/telegraf.conf
COPY ./bbb-telegraf.conf /etc/telegraf/plugins/bigbluebutton/bbb-telegraf.conf

ENTRYPOINT ["/bin/systemd", "--system", "--unit=multi-user.target"]
