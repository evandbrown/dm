go clean -i
gox -osarch="linux/amd64 darwin/amd64" -output="out/{{.OS}}_{{.Arch}}/dm" -ldflags "-X github.com/evandbrown/dm/commands.version $(git describe --tags)"
go install -ldflags "-X github.com/evandbrown/dm/commands.version $(git describe --tags)"
