FROM andrzejd/go-env

RUN go get -u github.com/andrzejd-pl/glog
RUN go install github.com/andrzejd-pl/glog

EXPOSE 80
ENV config_file /etc/glog/config.json

CMD glog -config=${config_file}
