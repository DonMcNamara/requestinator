# Requestinator

A *standalone executable* that has functionality inspired by [RequestBin](https://requestb.in/).

You create a "bin" that is represented by a URL. All HTTP interaction with that "bin" (URL) is recorded and can be interrogated using the API.

A good use of this service would be for end-to-end testing for your third party integrations. Or if you were developing a third party interaction and wanted to test with a local development enviroment instead of sending possibly sensitive information to an external internet service.

## Important differences

This project is in the spirit of RequestBin, but is not 100% compatible. The major differences include

1. Recorded traffic is sent to a slightly different URL path. RequestBin records traffic at a root level (`https://servername/[ID]`), this service puts all recording URLs down one level (`http://servername/bin/[ID]`).
2. This API returns recorded requests in chronological order. RequestBin gives request back in *reverse* chronological order.
3. Each request header value is an array of strings. HTTP headers can have multiple values. This project stores header values in an array. RequestBin joins multiple values so each header value is one string.
4. No HTML interface, yet.
5. No external persistence. All bins are stored in memory and lost when the server exits.

## Stability

This project is brand new and the API interface may change slightly.

## Installing Requestinator

For now, you can install using go. In the future I will probably provide precompiled binaries and probably a docker image.

```bash
go get -u github.com/DonMcNamara/requestinator
```

## Starting requestinator
```bash
# runs by default on 8080
requestinator

# run on port 9999
requestinator -p 9999  
```

## Usage
There is a API that you can hit with any http client. For example, curl.

### Create a bin
```bash
curl -X POST localhost:8080/api/v1/bins
```

Example response:
```
{"name":"akl5p55kim","request_count":0}
```

### Record an http action
```bash
curl -i -H "some-header: some-value" localhost:8080/bin/akl5p55kim?somequery=somequeryvalue
```

Response
```
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Sun, 01 Jan 2017 17:29:24 GMT
Content-Length: 2

ok
```

### Get bin information
```
curl localhost:8080/api/v1/bins/akl5p55kim
```

Response:
```
{"name":"akl5p55kim","request_count":1}
```

### Get request information
```
curl localhost:8080/api/v1/bins/akl5p55kim/requests

# pipe output through jq
curl localhost:8080/api/v1/bins/akl5p55kim/requests | jq
```

Response made pretty using jq:
```
[
  {
    "content_length": 0,
    "content_type": "",
    "time": 1483291764.0247948,
    "method": "GET",
    "body": "",
    "headers": {
      "Accept": [
        "*/*"
      ],
      "Some-Header": [
        "some-value"
      ],
      "User-Agent": [
        "curl/7.47.0"
      ]
    },
    "query_string": {
      "somequery": [
        "somequeryvalue"
      ]
    },
    "form_data": null
  }
]
```

## FAQ

### How is this different from RequestBin?
This is a standalone executable with nothing else to install. You run it locally rather than using a hosted service.

### Why did you write this?
We were using a self hosted RequestBin server for some integration tests. This seemed like a simpler solution to me.

### Wait, why don't you just stub your http client for your tests?
We do. Additionally, we want tests that exercise the HTTP client for coverage and regression detection. Third party integrations are important and we want end-to-end test coverage.

### You could just run RequestBin in a docker image.
That is not a question, but I get your point. We use OS X/macOS for local development. At one point we tried using docker for our test dependencies, but it was unwieldy. A standalone executable is simpler.

### How is the performance of this service?
100% untested. Don't use this if you require high performance.

### Should I run this in production?
No.

### Should I run this on the public internet?
No.

### What happens if I run this service for a long time?
The bins are stored in memory and do not time out. Each request to the service will be recorded and consume a little more memory. Eventually the service will run out of memory and crash.

### Is there an SBT plugin to run this as part of my build?
That's an odd question. Are you reading my mind? No, not yet.

### Can I mock an external service by having requestinator respond with predetermined HTTP responses?
Not yet, but I like the idea. If you're feeling up to it, open a pull request.

### Why did you use this non-idiomatic way of doing something with golang?
Because I don't know go very well. Let me know if something makes your golang brain angry.
