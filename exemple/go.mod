module github.com/caumette-co/x/example

go 1.21.0

replace github.com/caumette-co/x/xfoundation => ../xfoundation

replace github.com/caumette-co/x/xweb => ../xweb

replace github.com/caumette-co/x/xrenderer => ../xrenderer

require (
	github.com/caumette-co/x/xfoundation v0.0.0-00010101000000-000000000000
	github.com/caumette-co/x/xrenderer v0.0.0-00010101000000-000000000000
	github.com/caumette-co/x/xweb v0.0.0-00010101000000-000000000000
)

require (
	github.com/expectedsh/dig v0.0.1-expected // indirect
	github.com/go-chi/chi v1.5.4 // indirect
	github.com/samber/lo v1.38.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.25.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
)
