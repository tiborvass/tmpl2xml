<a href="https://pkg.go.dev/github.com/tiborvass/tmpl2xml#section-documentation" rel="nofollow"><img src="https://pkg.go.dev/badge/github.com/tiborvass/tmpl2xml" alt="Documentation"></a>
<a href="https://opensource.org/licenses/Apache-2.0">
    <img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg">
</a>

# tmpl2xml

I find Go templates hard to read and debug because it prioritizes text layout over logic.

This utility allows to convert them to XML which funnily makes it more readable by reprioritizing logic over text layout.

## CLI

```sh
$ go install github.com/tiborvass/tmpl2xml/cmd/tmpl2xml@latest
$ tmpl2xml go_template.tmpl
```

## Library

```sh
$ go get github.com/tiborvass/tmpl2xml
```

```go
package main

import (
	"fmt"

	"github.com/tiborvass/tmpl2xml"
)

func main () {
	out, err := tmpl2xml.String("Hello {{if .Cond}}world{{else}}friend{{end}}!")
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
```
