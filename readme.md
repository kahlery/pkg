# github.com/kahleryasla/pkg

`pkg` bag of packages for:

-   golang
-   ts/react

planning to add my Unity Engine packages/plugins here as well soon.

## Table of Contents

-   [Installation](#installation)
-   [Usage](#usage)
    -   [Using with NPM](#using-with-npm)
    -   [Using with Go](#using-with-go)
-   [Structure](#structure)

## Installation

### Using with NPM

To install a specific package from the TypeScript modules, you can use `npm` with the specific module you want to install. Here's how:

1. Ensure that you have `Node.js`, `npm` & especially `tiged` (tiged is an npm package allows you to install specific subdirectory in a whole project) installed on your machine.
2. Use the following command to install a specific module. For example, to install the `admin` module:
    ```bash
    tiged github.com/kahleryasla/pkg/ts/react/modules/core
    ```

To use a Go module, you can install it using go get. For example:

### Using with Go

Run the go get command with the package you want to install. For example, to install the auth package:

```bash
go get github.com/kahleryasla/pkg/go/auth
```

### Structure

#### Go Modules:

-   `go/auth`: Authentication utilities (JWT, hashing, etc.)

-   `go/aws`: AWS S3 service utilities

-   `go/env`: Environment-related utilities

-   `go/image`: Image processing utilities

-   `go/log`: Logging middleware and utilities

#### TypeScript Modules:

-   `ts/react/admin`: Admin panel components and utilities

-   `ts/react/auth`: Authentication components and utilities

-   `ts/react/common`: Common utilities shared across modules

-   `ts/react/core`: Core utilities shared across modules

-   `ts/react/layout`: Layout components for your React app

-   `ts/react/static`: Static components like images, CSS, etc.
