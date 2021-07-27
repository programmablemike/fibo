# fibo
A technical interview project for [Reserve Trust](https://www.reservetrust.com/).

## setup
requires:  
* Go 1.16

`make build` - Build the CLI/server binary
`make test`  - Run the unit & integration tests

## build
To build the software you can either call the `go build` command line directly or use
the built-in [Makefile](./Makefile) rules. When using the Makefile you should specify
the OS type (ex. `GOOS=darwin` and architecture (ex. `GOARCH=aarch64`) that corresponds
to your run-time environment.

For a list of all supported OS and architectures, run: `make list-build-targets`

```bash
> GOOS=linux GOARCH=amd64 
```

## docker
requires:  
* Docker
* Docker-Compose

This project includes a [docker-compose.yml](./docker-compose.yml) file that will start
a local instance of the API server for testing w/ the CLI.

```bash
# Start the app & database
> docker-compose up

# Connect the CLI
> fibo_darwin_arch64 -host <host> -port <port>
```

## load testing
requires:  
* [K6](https://k6.io/)

```bash
> k6 
```

## prompt
Expose a Fibonacci sequence generator through a web API that memoizes intermediate values.

The web API should expose operations to:
- [x] fetch the Fibonacci number given an ordinal (e.g. Fib(11) == 89, Fib(12) == 144),
- [ ] fetch the number of memoized results less than a given value (e.g. there are 12 intermediate results less than 120), and
- [ ] clear the data store.

The web API:
- [x] must be written in Go
- [x] Postgres must be used as the data store for the memoized results.

Please include:
- [x] tests for your solution, and
- [ ] a README.md describing how to build and run it.

Bonus points:
- [x] Use dockertest.
- [x] Include a Makefile.
- [ ] Include some data on performance.

## license
MIT licensed. See [licenses](./licenses) for more details.
