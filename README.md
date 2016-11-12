# Circuit Breaker pattern with NGINX
Let's see how to set up NGINX to implement [the circuit breaker pattern](http://martinfowler.com/bliki/CircuitBreaker.html) so that when upstream servers start to fail, nginx won't keep trying to connect to them, wasting precious time on each request. Instead, after a couple of retries, nginx will automatically send an error response. It'll try to connect to the failing upstream again after a while, to see if it's been recovered from the outage.

## Scenario
In a microservice architecture, we have two different services listening for http requests: `service1` and `service2`. These services always send the same response to every request: either `200 OK` responses that take 1 second to build, or `500 Server Error` responses that take 3 seconds to build if we toggle the error mode. We can toggle the error mode on and off just sending a `POST` request.

To achieve high availability, each of these services has three different servers.
All requests coming from the outside go through our API Gateway front server: an NGINX server routing traffic to our microservices based on the subdomain.

## Testing the scenario
When the service is failing, clients need to wait several seconds to get the error response. But when NGINX *opens the circuit* between itself and the failing upstream server, it sends the error response **in just a few milliseconds**.

### Building our infrastructure
Let's bring our infrastructure up

```bash
$ make up
```

### Sending requests to our services
Let's send requests to test it

```bash
$ time curl http://192.168.99.100:8000/
Error tio
curl http://192.168.99.100:8000/  0,00s user 0,00s system 0% cpu 9,012 total
```

```bash
nginx_1     | 2016/11/12 09:46:47 [warn] 5#5: *22 upstream server temporarily disabled while reading response header from upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://192.168.99.100:8001/", host: "192.168.99.100:8000"
nginx_1     | 2016/11/12 09:46:50 [warn] 5#5: *22 upstream server temporarily disabled while reading response header from upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://192.168.99.100:8001/", host: "192.168.99.100:8000"
nginx_1     | 192.168.99.1 - - [12/Nov/2016:09:46:53 +0000] "GET / HTTP/1.1" 500 10 "-" "curl/7.43.0" "-"
nginx_1     | 2016/11/12 09:46:53 [info] 5#5: *22 client 192.168.99.1 closed keepalive connection
```

```bash
$ time curl http://192.168.99.100:8000/
<html>
<head><title>502 Bad Gateway</title></head>
<body bgcolor="white">
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx/1.11.5</center>
</body>
</html>
curl http://192.168.99.100:8000/  0,00s user 0,00s system 0% cpu 3,010 total
```

```bash
nginx_1     | 2016/11/12 09:47:13 [error] 5#5: *26 no live upstreams while connecting to upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://service1/", host: "192.168.99.100:8000"
nginx_1     | 2016/11/12 09:47:13 [info] 5#5: *26 client 192.168.99.1 closed keepalive connection
nginx_1     | 192.168.99.1 - - [12/Nov/2016:09:47:13 +0000] "GET / HTTP/1.1" 502 173 "-" "curl/7.43.0" "-"
```

```bash
$ time curl http://192.168.99.100:8000/
Error tio
curl http://192.168.99.100:8000/  0,00s user 0,00s system 0% cpu 9,015 total
```

```bash
nginx_1     | 2016/11/12 09:47:19 [warn] 5#5: *28 upstream server temporarily disabled while reading response header from upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://192.168.99.100:8001/", host: "192.168.99.100:8000"
nginx_1     | 2016/11/12 09:47:22 [warn] 5#5: *28 upstream server temporarily disabled while reading response header from upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://192.168.99.100:8001/", host: "192.168.99.100:8000"
nginx_1     | 192.168.99.1 - - [12/Nov/2016:09:47:25 +0000] "GET / HTTP/1.1" 500 10 "-" "curl/7.43.0" "-"
nginx_1     | 2016/11/12 09:47:25 [info] 5#5: *28 client 192.168.99.1 closed keepalive connection
```

```bash
$ time curl http://192.168.99.100:8000/
<html>
<head><title>502 Bad Gateway</title></head>
<body bgcolor="white">
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx/1.11.5</center>
</body>
</html>
curl http://192.168.99.100:8000/  0,00s user 0,00s system 0% cpu 3,009 total
```

```bash
nginx_1     | 192.168.99.1 - - [12/Nov/2016:09:47:30 +0000] "GET / HTTP/1.1" 502 173 "-" "curl/7.43.0" "-"
nginx_1     | 2016/11/12 09:47:30 [error] 5#5: *32 no live upstreams while connecting to upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://service1/", host: "192.168.99.100:8000"
nginx_1     | 2016/11/12 09:47:30 [info] 5#5: *32 client 192.168.99.1 closed keepalive connection
```


```bash
time curl http://192.168.99.100:8000/
<html>
<head><title>502 Bad Gateway</title></head>
<body bgcolor="white">
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx/1.11.5</center>
</body>
</html>
curl http://192.168.99.100:8000/  0,00s user 0,00s system 0% cpu 3,010 total
```

```bash
nginx_1     | 192.168.99.1 - - [12/Nov/2016:09:47:35 +0000] "GET / HTTP/1.1" 502 173 "-" "curl/7.43.0" "-"
nginx_1     | 2016/11/12 09:47:35 [warn] 5#5: *34 upstream server temporarily disabled while reading response header from upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://192.168.99.100:8001/", host: "192.168.99.100:8000"
nginx_1     | 2016/11/12 09:47:35 [error] 5#5: *34 no live upstreams while connecting to upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://service1/", host: "192.168.99.100:8000"
nginx_1     | 2016/11/12 09:47:35 [info] 5#5: *34 client 192.168.99.1 closed keepalive connection
```


```bash
time curl http://192.168.99.100:8000/
<html>
<head><title>502 Bad Gateway</title></head>
<body bgcolor="white">
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx/1.11.5</center>
</body>
</html>
curl http://192.168.99.100:8000/  0,00s user 0,00s system 78% cpu 0,007 total
```

```bash
nginx_1     | 192.168.99.1 - - [12/Nov/2016:09:47:36 +0000] "GET / HTTP/1.1" 502 173 "-" "curl/7.43.0" "-"
nginx_1     | 2016/11/12 09:47:36 [error] 5#5: *36 no live upstreams while connecting to upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://service1/", host: "192.168.99.100:8000"
nginx_1     | 2016/11/12 09:47:36 [info] 5#5: *36 client 192.168.99.1 closed keepalive connection
```

```bash
time curl http://192.168.99.100:8000/
<html>
<head><title>502 Bad Gateway</title></head>
<body bgcolor="white">
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx/1.11.5</center>
</body>
</html>
curl http://192.168.99.100:8000/  0,00s user 0,00s system 78% cpu 0,007 total
```

```bash
nginx_1     | 192.168.99.1 - - [12/Nov/2016:09:47:38 +0000] "GET / HTTP/1.1" 502 173 "-" "curl/7.43.0" "-"
nginx_1     | 2016/11/12 09:47:38 [error] 5#5: *37 no live upstreams while connecting to upstream, client: 192.168.99.1, server: service1.armesto.local, request: "GET / HTTP/1.1", upstream: "http://service1/", host: "192.168.99.100:8000"
nginx_1     | 2016/11/12 09:47:38 [info] 5#5: *37 client 192.168.99.1 closed keepalive connection
```

# Further reading
- [Circuit Breaker](http://martinfowler.com/bliki/CircuitBreaker.html)
- [NGINX Passive Health Monitoring](https://www.nginx.com/resources/admin-guide/load-balancer/#health_passive)