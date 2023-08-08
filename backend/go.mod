module github.com/mcmaster-circ/canids-v2/backend

go 1.17

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.7.4
	github.com/joho/godotenv v1.3.0
	github.com/mcmaster-circ/canids-v2/protocol v0.0.0-00010101000000-000000000000
	github.com/olivere/elastic v6.2.27+incompatible
	github.com/olivere/elastic/v7 v7.0.11
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/sendgrid-go v3.5.0+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/tdewolff/minify v2.3.6+incompatible
	github.com/yl2chen/cidranger v1.0.2
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.26.0
)

require (
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/oschwald/maxminddb-golang v1.6.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sendgrid/rest v2.4.1+incompatible // indirect
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/tdewolff/test v1.0.6 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

replace github.com/mcmaster-circ/canids-v2/protocol => ../protocol
