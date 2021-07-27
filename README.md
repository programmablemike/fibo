# fibo
A technical interview project for [Reserve Trust](https://www.reservetrust.com/).

## setup
Requires:  
* [Go](https://golang.org) (v1.16)
* [Docker](https://www.docker.com) 
* [Docker-Compose](https://docs.docker.com/compose/install/)
* [K6][https://k6.io] (Optional. For running load tests.)
* []

## build
To build from source:  
```bash
# Build using standard Go tooling
> go build . 
# Build using the Makefile. This passes the --mod=vendor flag automatically for more reproducible builds.
> GOOS=darwin GOARCH=aarch64 make build 
```

## running the Docker Compose environment
To run the Docker Compose virtual environment:
```bash
# Build the Docker image
> make docker-build

# Start the application server and Postgres database
> make docker-start

# Check the logs for the application to make sure everything is OK
> make docker-logs-app

# Check the logs for the database
> make docker-logs-db

# Shutdown and clean the environment
> make docker-stop
```

## using the CLI
This project comes with a CLI for calling the API server.


### help documentation
```bash
mike@Mikes-MacBook-Pro fibo % ./fibo_darwin_arm64 help                                                                                      $(git_super_status)
Fibo is an API server and CLI client for generating Fibonacci sequences
It uses dynamic programming techniques (memoization) to speed up processing.

Usage:
  fibo [flags]
  fibo [command]

Available Commands:
  calculate   Calculates the Fibonacci number for the given ordinal N
  clear       Clears the memoizer cache
  completion  generate the autocompletion script for the specified shell
  count       Counts the number of ordinals in the Fibonacci value range (0, number)
  help        Help about any command
  server      Run the API server for memoized Fibonacci generation

Flags:
      --config string   config file (default is $HOME/.fibo.yaml)
      --debug           Turns on debugging mode
  -h, --help            help for fibo
      --host string     HTTP server hostname to bind (default: localhost) (default "localhost")
      --port int        HTTP server port to bind (default: 8080) (default 8080)

Use "fibo [command] --help" for more information about a command..
```


### calculating Fibonacci numbers for a given ordinal
```bash
# We're supporting very large Fibonacci values by using big.Int
mike@Mikes-MacBook-Pro fibo % ./fibo_darwin_arm64 calculate 10000

Fibonacci number: 33644764876431783266621612005107543310302148460680063906564769974680081442166662368155595513633734025582065332680836159373734790483865268263040892463056431887354544369559827491606602099884183933864652731300088830269235673613135117579297437854413752130520504347701602264758318906527890855154366159582987279682987510631200575428783453215515103870818298969791613127856265033195487140214287532698187962046936097879900350962302291026368131493195275630227837628441540360584402572114334961180023091208287046088923962328835461505776583271252546093591128203925285393434620904245248929403901706233888991085841065183173360437470737908552631764325733993712871937587746897479926305837065742830161637408969178426378624212835258112820516370298089332099905707920064367426202389783111470054074998459250360633560933883831923386783056136435351892133279732908133732642652633989763922723407882928177953580570993691049175470808931841056146322338217465637321248226383092103297701648054726243842374862411453093812206564914032751086643394517512161526545361333111314042436854805106765843493523836959653428071768775328348234345557366719731392746273629108210679280784718035329131176778924659089938635459327894523777674406192240337638674004021330343297496902028328145933418826817683893072003634795623117103101291953169794607632737589253530772552375943788434504067715555779056450443016640119462580972216729758615026968443146952034614932291105970676243268515992834709891284706740862008587135016260312071903172086094081298321581077282076353186624611278245537208532365305775956430072517744315051539600905168603220349163222640885248852433158051534849622434848299380905070483482449327453732624567755879089187190803662058009594743150052402532709746995318770724376825907419939632265984147498193609285223945039707165443156421328157688908058783183404917434556270520223564846495196112460268313970975069382648706613264507665074611512677522748621598642530711298441182622661057163515069260029861704945425047491378115154139941550671256271197133252763631939606902895650288268608362241082050562430701794976171121233066073310059947366875
```

### counting the number of ordinals given a max value
```bash
# We can also calculate the ordinals in the range of very large numbers too
mike@Mikes-MacBook-Pro fibo % ./fibo_darwin_arm64 count 33644764876431783266621612005107543310302148460680063906564769974680081442166662368155595513633734025582065332680836159373734790483865268263040892463056431887354544369559827491606602099884183933864652731300088830269235673613135117579297437854413752130520504347701602264758318906527890855154366159582987279682987510631200575428783453215515103870818298969791613127856265033195487140214287532698187962046936097879900350962302291026368131493195275630227837628441540360584402572114334961180023091208287046088923962328835461505776583271252546093591128203925285393434620904245248929403901706233888991085841065183173360437470737908552631764325733993712871937587746897479926305837065742830161637408969178426378624212835258112820516370298089332099905707920064367426202389783111470054074998459250360633560933883831923386783056136435351892133279732908133732642652633989763922723407882928177953580570993691049175470808931841056146322338217465637321248226383092103297701648054726243842374862411453093812206564914032751086643394517512161526545361333111314042436854805106765843493523836959653428071768775328348234345557366719731392746273629108210679280784718035329131176778924659089938635459327894523777674406192240337638674004021330343297496902028328145933418826817683893072003634795623117103101291953169794607632737589253530772552375943788434504067715555779056450443016640119462580972216729758615026968443146952034614932291105970676243268515992834709891284706740862008587135016260312071903172086094081298321581077282076353186624611278245537208532365305775956430072517744315051539600905168603220349163222640885248852433158051534849622434848299380905070483482449327453732624567755879089187190803662058009594743150052402532709746995318770724376825907419939632265984147498193609285223945039707165443156421328157688908058783183404917434556270520223564846495196112460268313970975069382648706613264507665074611512677522748621598642530711298441182622661057163515069260029861704945425047491378115154139941550671256271197133252763631939606902895650288268608362241082050562430701794976171121233066073310059947366875

Ordinals in this range: 10001
```

### clearing the memoizer cache
```bash
mike@Mikes-MacBook-Pro fibo % ./fibo_darwin_arm64 clear

Successfully cleared cache
```

## testing
To run the unit/integration tests. This additional runs basic benchmarks on memoized vs. non-memoized Fibonacci
computation using an in-memory cache for comparison.
```bash
mike@Mikes-MacBook-Pro fibo % make test
go test -v ./... --bench=.
?       github.com/programmablemike/fibo        [no test files]
?       github.com/programmablemike/fibo/cmd    [no test files]
=== RUN   TestCreateCache
INFO[0003] Successfully connected to database.          
INFO[0003] Initializing the database...                 
INFO[0003] Trying to connect to database...             
INFO[0003] Database connection succeeded.               
INFO[0003] Successfully initialized the table schemas.  
INFO[0003] Database initialized successfully!           
INFO[0003] Closing the database connection.             
--- PASS: TestCreateCache (0.04s)
=== RUN   TestReadWriteEntry
INFO[0003] Successfully connected to database.          
INFO[0003] Initializing the database...                 
INFO[0003] Trying to connect to database...             
INFO[0003] Database connection succeeded.               
INFO[0003] Successfully initialized the table schemas.  
INFO[0003] Database initialized successfully!           
INFO[0003] Closing the database connection.             
--- PASS: TestReadWriteEntry (0.05s)
PASS
ok      github.com/programmablemike/fibo/internal/cache 3.686s
=== RUN   TestFibonacciOrdinalCount
INFO[0000] Counting ordinals in range 0 to 120...       
INFO[0000] Counting ordinals in range 0 to 43466557686937456435688527675040625802564660517371780402481729089536555417949051890403879840079255169295922593080322634775209689623239873322471161642996440906533187938298969649928516003704476137795166849228875... 
--- PASS: TestFibonacciOrdinalCount (0.00s)
=== RUN   TestFibonacciNoCache
--- PASS: TestFibonacciNoCache (0.01s)
=== RUN   TestFibonacciCached
--- PASS: TestFibonacciCached (0.00s)
=== RUN   TestFibonacciLargeValue
--- PASS: TestFibonacciLargeValue (0.00s)
=== RUN   TestFibonacciVeryLargeValue
--- PASS: TestFibonacciVeryLargeValue (0.00s)
goos: darwin
goarch: arm64
pkg: github.com/programmablemike/fibo/internal/fibonacci
BenchmarkFibonacciNoCache
BenchmarkFibonacciNoCache-8          130           9022607 ns/op
BenchmarkFibonacciCached
BenchmarkFibonacciCached-8        563758              2148 ns/op
PASS
ok      github.com/programmablemike/fibo/internal/fibonacci     3.436s
?       github.com/programmablemike/fibo/internal/router        [no test files]
```

Note: The PostgreSQL cache tests use [Dockertest](https://github.com/ory/dockertest) and require Docker to be installed and running.

## load tests
There are some very basic load tests to get a feeling for the general throughput of 

```bash
# Start the Docker enviroment
> make docker-start

# Run the load tests
> make k6
```

## performance
### cached vs non-cached
Average ops/second difference between using no caching vs an in-memory cache using ordinals in the range `(0, 20)` inclusive.

As we can see there's a roughly 1,000x difference in the average speed of cached vs. non-cached operation.
```
BenchmarkFibonacciNoCache
BenchmarkFibonacciNoCache-8          130           9022607 ns/op
BenchmarkFibonacciCached
BenchmarkFibonacciCached-8        563758              2148 ns/op
```

### server throughput
Benchmarking of the HTTP server shows a throughput of `~746 requests/sec` calculating ordinals in the range `(0, 30)` inclusive.

```bash
mike@Mikes-MacBook-Pro fibo % make k6                                                                                                       $(git_super_status)
k6 run ./k6/calculate-worker.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: ./k6/calculate-worker.js
     output: -

  scenarios: (100.00%) 1 scenario, 10 max VUs, 1m0s max duration (incl. graceful stop):
           * default: 10 looping VUs for 30s (gracefulStop: 30s)


running (0m30.0s), 00/10 VUs, 22389 complete and 0 interrupted iterations
default ✓ [======================================] 10 VUs  30s

     ✓ is status 200
     ✓ gives correct fibonacci value

     checks.........................: 100.00% ✓ 44778     ✗ 0    
     data_received..................: 3.6 MB  120 kB/s
     data_sent......................: 2.2 MB  72 kB/s
     http_req_blocked...............: avg=5.4µs   min=0s       med=2µs    max=7.84ms   p(90)=3µs     p(95)=3µs    
     http_req_connecting............: avg=344ns   min=0s       med=0s     max=794µs    p(90)=0s      p(95)=0s     
     http_req_duration..............: avg=13.25ms min=704µs    med=7.11ms max=126.94ms p(90)=31.6ms  p(95)=39.49ms
       { expected_response:true }...: avg=13.25ms min=704µs    med=7.11ms max=126.94ms p(90)=31.6ms  p(95)=39.49ms
     http_req_failed................: 0.00%   ✓ 0         ✗ 22389
     http_req_receiving.............: avg=33.94µs min=5µs      med=25µs   max=2.11ms   p(90)=53µs    p(95)=66µs   
     http_req_sending...............: avg=11.64µs min=2µs      med=10µs   max=7.17ms   p(90)=15µs    p(95)=18µs   
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s     max=0s       p(90)=0s      p(95)=0s     
     http_req_waiting...............: avg=13.21ms min=672µs    med=7.07ms max=126.84ms p(90)=31.56ms p(95)=39.45ms
     http_reqs......................: 22389   746.03906/s
     iteration_duration.............: avg=13.39ms min=815.91µs med=7.25ms max=127.06ms p(90)=31.75ms p(95)=39.61ms
     iterations.....................: 22389   746.03906/s
     vus............................: 10      min=10      max=10 
     vus_max........................: 10      min=10      max=10 
```

Average request duration is `~13ms` which is likely a result of. 

Currently we're using a single generator shared between all request instances which constrains overall throughput at the
benefit of limiting the number of backend database connections being spawned. We would expect to see a significant improvement
in overall throughput if we transition to using a generator worker pool instead.

## prompt
Expose a Fibonacci sequence generator through a web API that memoizes intermediate values.

The web API should expose operations to:
- [x] fetch the Fibonacci number given an ordinal (e.g. Fib(11) == 89, Fib(12) == 144),
- [x] fetch the number of memoized results less than a given value (e.g. there are 12 intermediate results less than 120), and
- [x] clear the data store.

The web API:
- [x] must be written in Go
- [x] Postgres must be used as the data store for the memoized results.

Please include:
- [x] tests for your solution, and
- [x] a README.md describing how to build and run it.

Bonus points:
- [x] Use dockertest.
- [x] Include a Makefile.
- [x] Include some data on performance.

## license
MIT licensed. See [licenses](./licenses) for more details.
