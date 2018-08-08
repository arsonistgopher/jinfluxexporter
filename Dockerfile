FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

#RUN go build -o main .

CMD ["app", "-identity=vmx01", "-influxdb=junos", "-influxhost=http://192.168.10.200:8086", "-influxperiod=5", "-username=jet", "-password=Passw0rd", "-target=81.138.165.210"]
