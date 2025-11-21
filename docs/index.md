# Introduction & Installation

**Dino** is a lightweight HTTP library that contains utilities compatible with the `net/http` package.
If you want to work directly with `net/http` and the standard library, without relying on abstractions
(third-party routers, frameworks, etc.), but don't want to implement common functionalities from
scratch, Dino is the right library for you

This documentation contains guides for functionalities that you might need for your project but
that are not directly implemented in Dino (such as database access, monitoring, dependency injection,
etc.). This way, you'll get an idea of ​​how to develop an application in the Dino/Go style

We also provide some applications made with Dino so you can get inspired!

!!! note

    This documentation assumes you already have some experience with `net/http`. If you don't have it,
    check the `net/http` [cheat sheet](cheat-sheet/net-http.md) for a better understanding and context

## Philosophy

Dino doesn't try to be the magic solution that solves all your problems. Instead, it provides a thin
layer built on top of `net/http` for functionalities commonly needed when building applications

It also follows the Go philosophy, which is to maintain simplicity and not apply breaking changes or
major versions all the time. This makes Dino an easy-to-use and stable library in the long term

## Installation

Want to try Dino? Go get the package below and start reading the next sections of this documentation

```
go get github.com/willpinha/dino
```
