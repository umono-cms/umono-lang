# UmonoLang
UmonoLang is a domain-specific language to provide Umono's users to create page contents.

## Features
- CommonMark compatibility
- HTML output
- Components able to nested for reusability

```
package main

import (
	"fmt"

	ul "github.com/umono-cms/umono-lang"
	"github.com/umono-cms/umono-lang/converters"
)

func main() {
	umonoLang := ul.New(converters.NewHTML())
	umonoLang.SetGlobalComponent("HEADER", "This is a *global* component example.")

	umonoLangRaw := `HEADER
CONTENT

~ CONTENT
This is a **locale** component example.`

	fmt.Println(umonoLang.Convert(umonoLangRaw))
}

// Output
// <p>This is a <em>global</em> component example.
// This is a <strong>locale</strong> component example.</p>

```
