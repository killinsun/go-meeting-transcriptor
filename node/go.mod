module github.com/killinsun/go-meeting-transcriptor/node

go 1.19

require (
	github.com/gordonklaus/portaudio v0.0.0-20221027163845-7c3b689db3cc
	github.com/killinsun/go-meeting-transcriptor/backend v0.0.0-00010101000000-000000000000
	github.com/youpy/go-wav v0.3.2
	google.golang.org/grpc v1.54.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/youpy/go-riff v0.1.0 // indirect
	github.com/zaf/g711 v0.0.0-20190814101024-76a4a538f52b // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230323212658-478b75c54725 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/killinsun/go-meeting-transcriptor/backend => ../backend
