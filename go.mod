module language-detection-service

go 1.24.2

require (
	github.com/aws/aws-sdk-go v1.55.8
	github.com/hovman/ld-proto v0.1.0
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.9
)

replace github.com/hovman/ld-proto => ./ld_proto

require (
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250908214217-97024824d090 // indirect
)


