# If using GOPATH (old way):

- go get -u github.com/go-echarts/go-echarts/... if using GOPATH
- Don't use $GOPATH (deprecated way of handling go modules)

# If using go.mod go.sum (new way):

- export GO111MODULE="on" (default)
- go mod init github.com/nhatvu148/go-lang
- go get github.com/go-echarts/go-echarts/v2/...
- go mod tidy
- go clean -modcache
