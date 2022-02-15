module github.com/mcmaster-circ/canids-v2/ingestion

go 1.17

require (
	github.com/golang/protobuf v1.5.1
	github.com/google/uuid v1.1.2
	github.com/mcmaster-circ/canids-v2/protocol v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.36.0
	gopkg.in/urfave/cli.v1 v1.20.0
)

require (
	golang.org/x/net v0.0.0-20190311183353-d8887717615a // indirect
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

replace github.com/mcmaster-circ/canids-v2/protocol => ../protocol
