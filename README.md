# Reference:

- https://awesome-go.com/
- https://golang.org/pkg/

# If using GOPATH (old way):

- go get -u github.com/go-echarts/go-echarts/... if using GOPATH
- Don't use $GOPATH (deprecated way of handling go modules)

# If using go.mod go.sum (new way):

- export GO111MODULE="on" (go module mode by default, requires go.mod to work, if "off" use GOPATH)
- go mod init github.com/nhatvu148/go-lang
- go get github.com/go-echarts/go-echarts/v2/...
- go mod tidy
- go clean -modcache
- go install program@version
- go install golang.org/x/tools/gopls@latest

# Add icon to exe:

- go get github.com/akavel/rsrc
- rsrc -ico icon.ico
- go build -o chart.exe

# Use this chart library:

- https://github.com/gonum/plot/wiki/Example-plots
