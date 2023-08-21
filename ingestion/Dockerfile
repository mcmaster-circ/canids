# builder
FROM golang:1.20 as build

COPY ingestion/ /project/
WORKDIR /project
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s" -o /entrypoint main.go

# final
FROM zeek/zeek:6.0
RUN apt -y update && apt -y install cron

RUN mkdir /db
# COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
# COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /entrypoint /entrypoint
WORKDIR /db
ENTRYPOINT /entrypoint upload $flags --hostname $hostname /logs