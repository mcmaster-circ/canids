# golang alpine base
FROM golang:alpine as goalpine
RUN adduser -D -g '' gopher
RUN apk update
RUN apk --no-cache add ca-certificates tzdata git

# builder w/ git injection
FROM goalpine as build
COPY backend/ /project
COPY protocol/ /protocol
COPY .git/ /project/.git/
WORKDIR /project
RUN go mod vendor
RUN GIT_HASH=$(git rev-parse --short HEAD) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s -X main.gitHash=$GIT_HASH" -o /project/entrypoint main.go
RUN rm -rf .git/

# final
FROM alpine
RUN adduser -D -g '' gopher
RUN mkdir /storage
RUN chown gopher:gopher /storage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /project /project
WORKDIR /project
USER gopher
ENTRYPOINT ["/project/entrypoint"]
