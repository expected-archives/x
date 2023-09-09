module github.com/caumette-co/x/xrenderer

go 1.21.0

replace github.com/caumette-co/x/xfoundation => ../xfoundation

require (
	github.com/caumette-co/x/xfoundation v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.25.0
)

require go.uber.org/multierr v1.11.0 // indirect
