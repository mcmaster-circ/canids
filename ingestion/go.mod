module github.com/mcmaster-circ/canids-v2/ingestion

go 1.20

require (
	github.com/google/uuid v1.3.0
	github.com/mcmaster-circ/canids-v2/protocol v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/urfave/cli.v1 v1.20.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.13.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto v0.0.0-20230726155614-23370e0ffb3e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
)

replace github.com/mcmaster-circ/canids-v2/protocol => ../protocol
