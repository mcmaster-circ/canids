if [ "$1" = "-h" ] ; then
    echo "Usage: `basename $0` [Country] [State or Province] [City] [Backend Name] [Frontend Name] [Name] [Email Address (can be blank)]"
    echo "For example: gen.sh CA Ontario Hamilton FYELABS Ingestion Backend host.docker.internal"
    echo ""
    exit 0
fi

if [[ "$8" != "" ]]; then
    email="$8"
else
    email="."
fi

# Remove all key files generated previously
rm -rf ./temp
rm -rf ../cert
rm -rf ../ingestion/cert

# make required directories
mkdir -p ../ingestion/cert
mkdir -p ../cert
mkdir -p temp

# Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 9999 -nodes -keyout ./temp/ca-key.pem -out ../ingestion/cert/ca-cert.pem -subj "/C=$1/ST=$2/L=$3/O=$4/OU=$5/CN=$7/emailAddress=$email"

# Generate web server's private key and CSR
openssl req -newkey rsa:4096 -nodes -keyout ../cert/server-key.pem -out ./temp/server-req.pem -subj "/C=$1/ST=$2/L=$3/O=$4/OU=$6/CN=$7/emailAddress=$email"

# # Use CA's private key to sign web server's CSR and get the signed certificate
openssl x509 -req -in ./temp/server-req.pem -days 9999 -CA ../ingestion/cert/ca-cert.pem -CAkey ./temp/ca-key.pem -CAcreateserial -out ../cert/server-cert.pem -extfile ./server-ext.cnf

# remove temp folder
rm -rf ./temp
