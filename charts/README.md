# Command:

- go get -u github.com/go-echarts/go-echarts/... if using GOPATH
- Don't use $GOPATH (deprecated way of handling go modules)
- export GO111MODULE="on"

- If using mod sum (new way):
- go mod init github.com/nhatvu148/go-lang
- go get github.com/go-echarts/go-echarts/v2/...
- go mod tidy
- go clean -modcache
