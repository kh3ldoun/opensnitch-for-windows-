module github.com/evilsocket/opensnitch-windows/daemon

go 1.24.3

require (
	github.com/evilsocket/opensnitch-windows/daemon/proto v0.0.0-00010101000000-000000000000
	github.com/tailscale/wf v0.0.0-20240214030419-6fbb0a674ee6
	golang.org/x/sys v0.39.0
	google.golang.org/grpc v1.79.1
)

require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	go4.org/netipx v0.0.0-20220725152314-7e7bdc8411bf // indirect
	golang.org/x/exp/typeparams v0.0.0-20220218215828-6cf2b201936e // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	golang.org/x/tools/go/expect v0.1.1-deprecated // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	honnef.co/go/tools v0.3.2 // indirect
)

replace github.com/evilsocket/opensnitch/daemon => ../../daemon // Point to original daemon if needed, or local fork

replace github.com/evilsocket/opensnitch-windows/daemon/proto => ./proto
