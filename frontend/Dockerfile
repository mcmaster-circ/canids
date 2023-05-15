FROM node:16-alpine as yarnbuilder
WORKDIR /app
COPY frontend/package.json ./
COPY frontend/yarn.lock ./
RUN yarn install --frozen-lockfile

FROM yarnbuilder as builder
COPY frontend/ ./
RUN yarn build

FROM abiosoft/caddy:1.0.3-no-stats
COPY --from=builder /app /app
ENTRYPOINT ["/bin/parent", "caddy"]
CMD ["--conf", "/etc/Caddyfile", "--log", "stdout", "--agree=true"]
