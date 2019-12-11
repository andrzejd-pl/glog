FROM andrzejd/go-env

RUN go get -u github.com/andrzejd-pl/glog
RUN go install github.com/andrzejd-pl/glog

EXPOSE 80

CMD ["glog"]
