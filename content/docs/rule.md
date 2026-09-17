---
title: "Developing a Custom Rule"
description: "Implement your own linting rule on top of the revive framework."
weight: 200
icon: "code"
---

Each rule needs to implement the `lint.Rule` interface:

```go
type Rule interface {
    Name() string
    Apply(*File, Arguments) []Failure
}
```

The `Arguments` type is an alias of the type `[]interface{}`. The arguments of the rule are passed from the configuration file.

## Example

Let us develop a rule called `ArgumentsLimitRule` which sets a restriction on the number of arguments that a function may accept. As we explained above, our rule needs to implement the `Rule` interface:

```go
package rule

import (
    "fmt"
    "go/ast"

    "github.com/mgechev/revive/lint"
)

// ArgumentsLimitRule lints given else constructs.
type ArgumentsLimitRule struct{}

// Apply applies the rule to given file.
func (r *ArgumentsLimitRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
    // ...
    return failures
}

// Name returns the rule name.
func (r *ArgumentsLimitRule) Name() string {
    return "argument-limit"
}

type lintArgsNum struct {
    total     int
    onFailure func(lint.Failure)
}

func (w lintArgsNum) Visit(n ast.Node) ast.Visitor {
    // ...
}
```

In the snippet above, we first import the packages `fmt`, `go/ast` and the `lint` package from `revive`. After that we define a struct called `ArgumentsLimitRule` which has two methods `Apply` and `Name`. `Apply` accepts a file and the rule's argument, and `Name` simply returns the name of the rule.

Notice how we also have the declaration of the `lintArgsNum` struct. We're going to use it to traverse the program's AST.

Now let us fill the blanks!

```go
// ...
// Apply applies the rule to given file.
func (r *ArgumentsLimitRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
    if len(arguments) != 1 {
        panic(`invalid configuration for "argument-limit"`)
    }

    total, ok := arguments[0].(int64) // Alt. non panicking version
    if !ok {
        panic(`invalid value passed as argument number to the "argument-list" rule`)
    }

    var failures []lint.Failure

    walker := lintArgsNum{
        total: int(total),
        onFailure: func(failure lint.Failure) {
            failures = append(failures, failure)
        },
    }

    ast.Walk(walker, file.AST)

    return failures
}
// ...
```

In the example above we perform the following algorithm:

- We first validate the arguments. The rule should accept a single argument of type `int64`, specifying the total number of allowed arguments
- After that, we define a slice of type `[]lint.Failure` which we'll fill during the AST traversal
- As the next step, we create a new instance of the struct `lintArgsNum`. The struct has two fields: `total`, which specifies the rule's arguments, and `onFailure`. The struct's instance will invoke the `onFailure` callback when it finds a new failure. Inside of this callback, we'll push the failure to the list of failures
- Next, we start the AST traversal
- Finally, we return the slice of collected failures

That's pretty much everything in the rule's declaration!

Now let us explore the AST traversal algorithm:

```go
type lintArgsNum struct {
    total     int
    onFailure func(lint.Failure)
}

func (w lintArgsNum) Visit(n ast.Node) ast.Visitor {
    node, ok := n.(*ast.FuncDecl)
    if ok {
        num := 0
        for _, l := range node.Type.Params.List {
            for range l.Names {
                num++
            }
        }
        if num > w.total {
            w.onFailure(lint.Failure{
                Confidence: 1,
                Failure:    fmt.Sprintf("maximum number of arguments per function exceeded; max %d but got %d", w.total, num),
                Node:       node.Type,
            })
            return w
        }
    }
    return w
}
```

The `lintArgsNum` struct implements the `ast.Visitor` interface, which means that it has a method called `Visit` which accepts an `ast.Node` and returns an `ast.Visitor`.

Inside of the `Visit` method, we check if the current node is of type `*ast.FuncDecl`, and if it is, we validate the number of arguments the given function accepts. If it accepts more than the value of the property `total`, then we invoke the `onFailure` callback with the failure object.

## Conclusion

In this guide we saw how we can implement a new rule in the framework for static code analysis for Go, revive.

We first explained how revive passes arguments defiled in the config file to the rules. After that, we looked at the `lint.Rule` interface and implemented it in the rule `ArgumentsLimitRule`. In the rule, we traversed the AST using the `lintArgsNum` visitor which collected all the warnings it found in the process.
