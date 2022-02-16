# CANIDS

## Running in Development
An initial build will take a few minutes to setup the docker environment and run. Ensure that the keys are generated prior to building for development. Consecutive builds will utilize caching and will be much faster.

```sh
# generate certs
cd certfiles
./gen-certs.sh CA Ontario Hamilton FYELABS Ingestion Backend host.docker.internal
cd ..

# start server
docker-compose -f deploy-dev.yml up -d --build

# start ingestion client
cd ../ingestion
docker-compose up -d --build
```

Opens ports 80 (HTTP), 443 (HTTPS), 50000 (Ingestion client upload), 9200 (elasticsearch), 5601 (kibana), 8080 (swagger documentation) & 6060 (backend)

## Building for Production
```sh
# Update Caddyfile with your server's URL if it has one

# generate certs
cd certfiles
# run ./gen-certs.sh -h to see what each parameter represents
./gen-certs.sh CA Ontario Hamilton FYELABS Ingestion Backend host.docker.internal
cd ..

# build & package server
docker-compose -f deploy-prod.yml build

mkdir canids-release-v2.0.0
docker save mcmaster-circ/canids-v2-backend > canids-release-v2.0.0/canids-backend-v2.0.0.tar
docker save mcmaster-circ/canids-v2-frontend > canids-release-v2.0.0/canids-frontend-v2.0.0.tar
cp ./deploy-prod.yml ./canids-release-v2.0.0/docker-compose.yml
cp ./Caddyfile ./canids-release-v2.0.0/Caddyfile
cp -R ./cert ./canids-release-v2.0.0/cert
cp -R ./config ./canids-release-v2.0.0/config

tar -czvf canids-release-v2.0.0.tar.gz canids-release-v2.0.0/

# build & package ingestion client
cd ingestion
docker-compose build
cd ..

mkdir canids-ingest-v2.0.0
docker save mcmaster-circ/canids-v2-ingestion > canids-ingest-v2.0.0/canids-ingestion-v2.0.0.tar
cp ./ingestion/docker-compose.yml ./canids-ingest-v2.0.0/docker-compose.yml
cp -R ./ingestion/cert ./canids-ingest-v2.0.0/cert

tar -czvf canids-ingest-v2.0.0.tar.gz canids-ingest-v2.0.0/
```

## Running in Server Production
```sh
tar xzf canids-release-v2.0.0.tar.gz
cd canids-release-v2.0.0

docker load --input canids-backend-v2.0.0.tar
docker load --input canids-frontend-v2.0.0.tar

docker-compose up -d --no-build
```

Opens ports 80 (HTTP), 443 (HTTPS) & 50000 (Ingestion client upload)

## Running in Ingestion Client Production
```sh
tar xzf canids-ingest-v2.0.0.tar.gz
cd canids-ingest-v2.0.0

docker load --input canids-ingestion-v2.0.0.tar

docker-compose up -d --no-build

# logs will be read from the ./logs directory
```

## Generating SSH keys
Navigate to `/certfiles`. `gen-certs.sh` contains the commands to run to generate the keys.

### Parameters
Please ensure `config/config.env` and `config/secret.env` are present. Any
adjustable parameters may be found in these files.

## Backend Documentation
The `docker-compose` command launches Swagger. API documentation is listening on
`http://localhost:8080`.
