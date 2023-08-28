# builder
FROM golang:1.20 as build

COPY ingestion/ /project/
WORKDIR /project
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s" -o /entrypoint main.go

#final
FROM alpine
RUN adduser -D -g '' gopher
RUN mkdir /db
RUN chown gopher:gopher /db
# COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /entrypoint /entrypoint
WORKDIR /db
USER gopher
ENTRYPOINT /entrypoint upload $flags --hostname $hostname /logs