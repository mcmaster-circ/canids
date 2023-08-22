# builder
FROM golang:1.20 as build

COPY ingestion/ /project/
WORKDIR /project
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s" -o /entrypoint main.go

# final
FROM zeek/zeek:6.0
RUN apt -y update && apt -y install cron
COPY ./ingestion/config/zeekctl.cfg /usr/locla/bin/etc/zeekctl.cfg
COPY --chmod=644 ./ingestion/config/zeek-crontab /etc/cron.d/zeek
#TODO: How to point to backend on host network?
RUN mkdir /db
RUN mkdir /logs
# COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
# COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /entrypoint /entrypoint
WORKDIR /db
USER gopher
ENTRYPOINT /entrypoint upload $flags --hostname $hostname /logs