module github.com/mcmaster-circ/canids-v2/backend

go 1.20

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.5.1
	github.com/mcmaster-circ/canids-v2/protocol v0.0.0-00010101000000-000000000000
	github.com/olivere/elastic v6.2.37+incompatible
	github.com/olivere/elastic/v7 v7.0.32
	github.com/oschwald/geoip2-golang v1.9.0
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/sendgrid-go v3.12.0+incompatible
	github.com/sirupsen/logrus v1.9.3
	github.com/tdewolff/minify v2.3.6+incompatible
	github.com/yl2chen/cidranger v1.0.2
	golang.org/x/crypto v0.11.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/oschwald/maxminddb-golang v1.12.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/tdewolff/test v1.0.6 // indirect
	golang.org/x/net v0.13.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
)

replace github.com/mcmaster-circ/canids-v2/protocol => ../protocol
