# Yendo

![License](https://img.shields.io/github/license/jamiefdhurst/yendo.svg)
[![Build Status](https://ci.jamiehurst.co.uk/buildStatus/icon?job=github%2fyendo%2Fmaster)](https://ci.jamiehurst.co.uk/job/github/job/yendo/job/master/)
[![Latest Version](https://img.shields.io/github/release/jamiefdhurst/yendo.svg)](https://github.com/jamiefdhurst/yendo/releases)

A small and simple MySQL database package for Go that includes support for automatic migrations.

Currently supports Golang 1.13 - 1.15.

## Usage Example

A detailed example is present in the `example/` folder. 

## Running Tests

The tests require a MySQL connection available - the following environment variables are used to establish the 
connection:

* `DB_HOST` - Hostname or IP address, e.g. `localhost`
* `DB_PORT` - Port number, usually 3306 - this can be ommitted
* `DB_USER` - Username to connect as
* `DB_PASSWORD` - Password to connect using
* `DB_NAME` - Database name to use when testing

If you check out the code into the standard GOPATH-expected location (`src/github.com/jamiefdhurst/yendo`), you can run 
the tests immediately from within that location.

Once these conditions are met, the tests can be performed simply by running:

```bash
go test
```
