# Reference:

- https://yourbasic.org/golang/
- https://zetcode.com/all/#go
- https://gobyexample.com/
- https://www.youtube.com/watch?v=SmoM1InWXr0
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
- go clean -i -x github.com/galdor/go-cmdline # remove a package
- go run main.go -database="jmu" -host localhost -user "root" -password 123456789 -shipInfoID 1
- .\report.exe --outDir "C:/Users/nhatv/OneDrive/Desktop/test/output"

# Add icon to exe:

- go get github.com/akavel/rsrc
- rsrc -ico icon.ico
- go build -o chart.exe

# Use this chart library:

- https://github.com/gonum/plot/wiki/Example-plots

# Excel:

- https://www.get-digital-help.com/use-a-map-in-an-excel-chart/
- https://support.microsoft.com/en-us/office/create-a-table-in-excel-bf0ce08b-d012-42ec-8ecf-a2259c9faf3f
- https://support.microsoft.com/en-us/office/overview-of-formulas-in-excel-ecfdc708-9162-49e8-b993-c311f47ca173

# When importing local modules:

Main method:
- In async folder: go mod init github.com/nhatvu148/go-lang/async
- Enter hackerrank folder: go mod edit -replace=github.com/nhatvu148/go-lang/async=../async
- go get github.com/nhatvu148/go-lang/async
- Reference: https://golang.org/doc/tutorial/call-module-code
- import in code: 
```
import (
    "github.com/nhatvu148/go-lang/async"
)
```

Or using vendor folder:
- go mod tidy
- go mod vendor
- copy local package folders to the vendor folder
- go run .

Or if using GOPATH:
- env GO111MODULE=off go run await.go
- env GO111MODULE=off go build await.go

# Check Go Environment:

- go env
