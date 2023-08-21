FROM node:16-alpine as npmbuilder
WORKDIR /project/frontend
COPY frontend/package.json ./
RUN npm install

FROM npmbuilder as frontendbuilder
COPY . /project
RUN npm run build
# golang alpine base
FROM golang:alpine as goalpine
COPY --from=frontendbuilder /project /project
WORKDIR /project/backend
RUN adduser -D -g '' gopher
RUN apk update
RUN apk --no-cache add ca-certificates tzdata git

# builder w/ git injection
FROM goalpine as build
COPY .git/ /project/backend/.git/
RUN go mod vendor
RUN GIT_HASH=$(git rev-parse --short HEAD) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOWORK=off go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s -X main.gitHash=$GIT_HASH" -o /project/entrypoint main.go
RUN rm -rf .git/

# final
FROM alpine
RUN adduser -D -g '' gopher
RUN mkdir /storage
RUN chown gopher:gopher /storage
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /project/backend /project
COPY --from=build /project/entrypoint /project/entrypoint
WORKDIR /project
USER gopher
ENTRYPOINT ["/project/entrypoint"]
