FROM node:16-alpine as yarnbuilder
WORKDIR /app
COPY frontend/package.json ./
COPY frontend/yarn.lock ./
RUN yarn install --frozen-lockfile

FROM yarnbuilder as builder
COPY frontend/ ./
RUN yarn build

FROM caddy:2.6.4-alpine
COPY --from=builder /app /app
ENTRYPOINT ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]
