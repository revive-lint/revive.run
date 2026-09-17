---
title: "Developing a Custom Formatter"
description: "Implement your own output formatter for revive."
weight: 300
icon: "description"
---

Each formatter needs to implement the following interface:

```go
type Formatter interface {
    Format(<-chan Failure, RulesConfig) (string, error)
    Name() string
}
```

The `Format` method accepts a channel of `Failure` instances and the configuration of the enabled rules. The `Name()` method should return a string different from the names of the already existing rules. This string is used when specifying the formatter when invoking the `revive` CLI tool.

## Example

Let us implement a sample formatter which outputs the rule failures to plain text.

First, you need to clone the repository and create a file under `revive/formatter` called `sample.go`. After that, add the following imports and struct declaration:

```go
package formatter

import (
    "fmt"

    "github.com/mgechev/revive/lint"
)

type Sample struct {
    Metadata lint.FormatterMetadata
}
```

Now, `Sample` has to implement the `Formatter` interface, which means that we need to define two methods: `Name` which returns the name of the formatter, and the method `Format`. As a first argument `Format` receives a channel where we'll receive the failures, and as second argument the method accepts the user's configuration.

```go
func (f *Sample) Name() string {
    return "sample"
}

// Format formats the failures gotten from the lint.
func (f *Sample) Format(failures <-chan lint.Failure, _ lint.RulesConfig) (string, error) {
    for failure := range failures {
        fmt.Printf("%v: %s\n", failure.Position.Start, failure.Failure)
    }
    return "", nil
}
```

In the implementations above, `Name` returns the string `"sample"` and `Format` prints the failures one by one, first displaying the failure's position and after that its description.
