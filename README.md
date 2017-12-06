Testing FDK-based functions
===========================

Function development kit (FDK) as a piece of software that helps to write hot functions by encapsulating logic of processing incoming requests with respect to defined protocol (HTTP, JSON).
Particular testing framework help developers to identify if any programming language-specific FDK compatible with Fn hot formats.


Prerequisites
-------------

This testing framework allows to run FDK tests against live Fn service, to let tests know of where Fn service is hosted please set following environment variable:
```bash
    export FN_API_URL=http://fn.io:8080
```

Test suite requires general purpose programming language-specific FDK-based function image that must be developed specifically for this test suite, following environment variable must be set:
```bash
    export FDK_FUNCTION_IMAGE="..."
```
This environment variable should contain a reference to the particular docker image.

At this moment Fn server supports 3 formats:

 - default
 - http
 - json

To let developers freedom of choice test suite allows to configure list of formats to test against FDK.
By default tests will use two hot formats: `json` and `http`, it's possible to override those formats using following environment format:
```bash
    export FDK_FORMATS=json,http
```
Test suite details
------------------

Test suite contains following tests:

1. `TestFDKFormatSmallBody`

`TestFDKFormatSmallBody` test
--------------------------

FDK should support following formats:

 - HTTP (subtest: `test-fdk-http-small-body`)
 - JSON (subtest: `test-fdk-json-small-body`)

Request input body:
```json
{
  "name": "John"
}
```

Response output body:
```text
Hello John
```

How to run tests?
-----------------

There are couple options: from source code, from release binary

From source code:
```bash
    go test -v ./...
```

How to build test binary executable?
------------------------------------

Regular `go build` does not work with tests, so following command will create a binary executable for this particular test suite:
```bash
    go test -c -i
```

Sample FDK-based functions
--------------------------

As an example test suite supplied with general purpose test [function written with FDK Python](./functions/python).
