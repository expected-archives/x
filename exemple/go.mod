module github.com/caumette-co/x/example

go 1.21.0

replace github.com/caumette-co/x/xfoundation => ../xfoundation

replace github.com/caumette-co/x/xweb => ../xweb

require (
	github.com/caumette-co/x/xfoundation v0.0.0-00010101000000-000000000000 // indirect
	github.com/caumette-co/x/xweb v0.0.0-00010101000000-000000000000 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.25.0 // indirect
)
