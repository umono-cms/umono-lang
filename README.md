# UmonoLang
A domain-specific language to provide Umono's users to create page contents.

[![Go Report Card](https://goreportcard.com/badge/github.com/umono-cms/umono-lang)](https://goreportcard.com/report/github.com/umono-cms/umono-lang)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## Features
- CommonMark compatibility
- HTML output
- Components able to nested for reusability

```go
package main

import (
	ul "github.com/umono-cms/umono-lang"
	"github.com/umono-cms/umono-lang/converters"
)

func main() {
	umonoLang := ul.New(converters.NewHTML())
	umonoLang.SetGlobalComponent("HEADER", "This is a *global* component example.")

	input := `HEADER

CONTENT

{{ LINK text="click me!" url="https://github.com/umono-cms/umono" new-tab=true }}

~ CONTENT
This is a **local** component example.`

  output := umonoLang.Convert(input)
  // <p>This is a <em>global</em> component example.</p><p>This is a <strong>local</strong> component example.</p><p><a href="https://github.com/umono-cms/umono" target="_blank" rel="noopener noreferrer">click me!</a></p>
}
```
