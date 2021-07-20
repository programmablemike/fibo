# fibo
A technical interview project for [Reserve Trust]().

## prompt
Expose a Fibonacci sequence generator through a web API that memoizes intermediate values.
The web API should expose operations to
- [ ] fetch the Fibonacci number given an ordinal (e.g. Fib(11) == 89, Fib(12) == 144),
- [ ] fetch the number of memoized results less than a given value (e.g. there are 12 intermediate results less than 120), and
- [ ] clear the data store.

The web API:
- [x] must be written in Go
- [ ] Postgres must be used as the data store for the memoized results.

Please include:
- [ ] tests for your solution, and
- [ ] a README.md describing how to build and run it.

Bonus points:
- [ ] Use dockertest.
- [ ] Include a Makefile.
- [ ] Include some data on performance.

## license
MIT licensed. See [licenses](./licenses) for more details.
