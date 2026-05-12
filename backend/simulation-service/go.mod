module github.com/ykshvn/reactive-dsr/simulation-service

go 1.26.3

replace github.com/ykshvn/reactive-dsr/shared => ../shared

require (
	github.com/ykshvn/reactive-dsr/shared v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.28.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)
