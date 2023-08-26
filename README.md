# Token Bucket List

This page explains how the code works, if you wanna know what is the pourpose of a rate limiter and its detail you should checkout my text on medium or dev.to

## Setup

To execute the bucket list script all you need to have is Golang.

I am using `go1.18.3 darwin/arm64` not sure how it will perform on other Go versions and OS but I think it is not going to be a problem. Let me know if works on differents system configuration.

To execute the whole environment you'll need:

- Docker
- Docker Compose
- K6 CLI
- Go Compiler


## How to

First run `docker-compose up`to run the environment. You should have running 3 services:

- `Token Bucket Service`
- `Prometheus`
- `Grafana`

Make requests on `http://localhost:8080` (it can be from your very browser) you should recieve a `200` and if you spam requests should start receiving some `429`. This is great!

`http://localhost:9090` should render the prometheus page, in this page you can query the Metrics you created on prometheus

`http://localhost:3000` goes to the graphana webpage, for loggin in, use username `admin` and password `admin`. 
Add prometheus as data source and create a dashboard to see your metrics :)

To stress test the application simply run `k6 run stressTest.js`
ß


## Code Analysis

### Dependencies

The only dependency is the client of `prometheus` for the sake of developing metrics, the prometheus client comes with other dependencies that can be analyized on the `go.mod` file

### Prometheus Metrics

There are only 2 metrics that the program generates:

- `http_request_conumed_tokens`: Indicates how much tokens are consumed therefore indicating a successful request on the rate limiter

- `http_request_denied_requests`: Indicates requests that are denied due to lack of tokens.

### Details
Basically the code runs a `goroutine` that keeps adding tokens, if the tokens added are above the maxToken then it overflow the tokens to be added. It exposes an endpoint that is only responsible to consume a token if possible and return `Http Status Code 200` otherwise, return a `429 Too Many Requests`




