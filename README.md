Go Verify Apple ID Token
======

![](https://img.shields.io/badge/golang-1.17-blue.svg?style=flat)

A library for verify Apple ID token.

## Installation
```
go get github.com/pikann/go-verify-apple-id-token
import "github.com/pikann/go-verify-apple-id-token/apple"
```

## Usage
There is example file based on your particular use case which can be found below:
- [Verify an Apple ID token](examples/main.go)

### Example
While it is recommended to look at the specific example file, here is validating an app token:
``` golang
import "github.com/pikann/go-verify-apple-id-token/apple"

...

resp, err := apple.VerifyAppleIdToken(
    "token_test",
    "client_id",
)
```

## Contact
Van Hai Huynh - hvhai22@gmail.com

Project Link: https://github.com/pikann/CrawlGoogleImage