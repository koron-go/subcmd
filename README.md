# koron-go/subcmd

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/subcmd)](https://pkg.go.dev/github.com/koron-go/subcmd)
[![Actions/Go](https://github.com/koron-go/subcmd/workflows/Go/badge.svg)](https://github.com/koron-go/subcmd/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/subcmd)](https://goreportcard.com/report/github.com/koron-go/subcmd)

koron-go/subcmd is very easy and very simple sub-commander library.
It focuses solely on providing a hierarchical subcommand mechanism.
It does not provide any flags nor options.

## Getting Started

Install or update:

```console
$ go install github.com/koron-go/subcmd@latest
```

The basic usage consists of just 3 steps.

1. Define commands with function.

    ```go
    package item

    import "github.com/koron-go/subcmd"

    var List = subcmd.DefineCommand("list", "list items", func(ctx context.Context, args []string) error {
        // TODO: list your items (A)
        return nil
    })

    var Add = subcmd.DefineCommand("add", "add a new item", func(ctx context.Context, args []string) error {
        // TODO: add a your item (B)
        return nil
    })
    ```

2. Define command sets by bundling the defined commands.

    ```go
    package item

    var CommandSet = subcmd.DefineSet("item", List, Add)
    ```

    Similarly, assume that multiple commands (List/Add/Delete) have been defined for `user` package.

    ```go
    package user

    var CommandSet = subcmd.DefineSet("user", List, Add, Delete)
    ```

3. Call function `subcmd.Run()` with the defined set from `main()`.

    ```go
    package main

    import (
        "log"
        "os"

        "github.com/koron-go/subcmd"

        "{some/pacakge/path}/item"
        "{some/package/path}/user"
    )

    var rootSet = subcmd.DefineRootSet(
        item.CommandSet,
        user.CommandSet,
    )

    func main() {
        if err := subcmd.Run(rootSet, os.Args[1:]...); err != nil {
            log.Fatal(err)
        }
    }
    ```

This will enable you to use subcommands such as the following:

```console
# This call (A)
$ myprog item list

# This call (B)
$ myprog item add

# List subcommands in "item" command set
$ myprog item --help

# Same about "user" command set.
$ myprog user list
$ myprog user add
$ myprog user delete
$ myprog user --help
```

## Concept

* Want to write each command as a function.
* Want to associate command name and description with a function.
* Want to make commands and command sets reusable, and re-assemblable.
* Don't want to use types that are not in the standard library in the function signature.
* Don't want to declare a `struct` or define a type instead of a function.
