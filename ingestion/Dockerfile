# golang alpine base
FROM golang:alpine as goalpine
RUN adduser -D -g '' gopher
RUN apk update
RUN apk --no-cache add ca-certificates tzdata git

# builder w/ git injection
FROM goalpine as build
COPY injestion/ /project/
COPY protocol/ /protocol
WORKDIR /project
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s" -o /entrypoint main.go

# final
FROM alpine
RUN adduser -D -g '' gopher
RUN mkdir /db
RUN chown gopher:gopher /db
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /entrypoint /entrypoint
WORKDIR /db
USER gopher
ENTRYPOINT /entrypoint upload $flags --asset $asset --hostname $hostname /logs