FROM golang

ENV IDENTITY="vmx01"
ENV TOPIC="important"
ENV PERIOD=1
ENV KAFKAPORT=9092
ENV KAFKAHOST="10.42.0.1"
ENV TARGET="10.42.0.132"

RUN apt-get install -y git && \
    go get github.com/arsonistgopher/jkafkaexporter

CMD jkafkaexporter -identity=$IDENTITY -kafkatopic=$TOPIC -kafkaperiod=$PERIOD -kafkaport=$KAFKAPORT -kafkahost=$KAFKAHOST -password Passw0rd -username jet -sshport 22 -target=$TARGET