package main

//go:generate protoc --go_out=proto --proto_path proto proto/idl/config/FetchMailConfig.proto proto/idl/email/EMail.proto
//go:generate protoc --go_out=proto --proto_path proto proto/idl/archive/Envelope.proto

func main() {}
