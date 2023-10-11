module github.com/caumette-co/x/xfoundation

go 1.21.0

require (
	github.com/expectedsh/dig v0.0.1-expected
	github.com/samber/lo v1.38.1
	go.uber.org/zap v1.25.0
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
)

// Note on github.com/expectedsh/dig, this is a fork of go.uber.org/dig
// Our fork is needed to add outputs informations for an invoke
// A merge request is opened, waiting for the merge
// When the merge will occur we could remove this fork and use the official dig package
//
// The fork contains two branches :
// - the master branch is the same as the official dig package
// - the expected branch contains the changes we made to add outputs informations for an invoke and replaced go.mod module name
