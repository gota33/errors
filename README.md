# errors

[![Go Doc](https://godoc.org/github.com/nathany/looper?status.svg)](https://pkg.go.dev/github.com/gota33/errors) [![Go Report](https://goreportcard.com/badge/github.com/gota33/errors)](https://goreportcard.com/report/github.com/gota33/errors) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gota33/errors/blob/master/LICENSE)

Package error provides a simple error info binding in RESTful service. error definition follow the [Google API design](https://cloud.google.com/apis/design/errors).

`go get github.com/gota33/errors`

## Adding status & error_detail to an error

``` go
// import . "github.com/gota33/errors"

_, err := db.QueryRow("SELECT * FROM cats WHERE name = 'white' LIMIT 1")

if err != nil {
    detail := ResourceInfo{
        ResourceType: "pet.com/pet.v1.Cat",
        ResourceName: "cat123@pet.com",
        Owner:        "user123",
        Description:  "cat123 not found",
    }
    
    // Standard way
    err = Annotate(err, NotFound, detail)
    
    // Or use predifined wrapper
    err = WithNotFound(cause, detail)
    return
}

```

## Retrieving the code of an error

``` go
// ...
code := Code(err)
fmt.Println(code)
// Output: 404 NOT_FOUND
```

## Check temporary

``` go
// ...
fmt.Println(Temporary(OK))
fmt.Println(Temporary(Internal))
fmt.Println(Temporary(Unavailable))
fmt.Println(Temporary(&net.DNSError{IsTemporary: true}))

// Output:
// false
// false
// true
// true
```

## Retrieving the details of an error

``` go
// ...
for _, detail := range Details(err) {
    fmt.Printf("%+v", detail)
}

// Output:
// type: "type.googleapis.com/google.rpc.ResourceInfo"
// resource_type: "pet.com/pet.v1.Cat"
// resource_name: "cat123@pet.com"
// owner: "user123"
// description: "cat123 not found"
```

## Print formatted message

 ``` go
 // ...
 err := Annotate(
     context.DeadlineExceeded,
     DeadlineExceeded,
     StackTrace("heavy job"),
     RequestInfo{RequestId: "<uuid>"},
     LocalizedMessage{Local: "en-US", Message: "Background task timeout"},
     LocalizedMessage{Local: "zh-CN", Message: "后台任务超时"},
 )
 
 fmt.Printf("%+v", err)
 
 // Output:
 // status: "504 DEADLINE_EXCEEDED"
// message: "context deadline exceeded"
// detail[0]:
// 	type: "type.googleapis.com/google.rpc.DebugInfo"
// 	detail: "heavy job"
// 	stack:
// 		goroutine 1 [running]:
// 		runtime/debug.Stack(0xc00005e980, 0x40, 0x40)
// 			/home/user/go/src/runtime/debug/stack.go:24 +0xa5
// 		github.com/gota33/errors.StackTrace.Annotate(0xfe36af, 0x9, 0x1056490, 0xc00005e980)
// 			/home/user/github/gota33/errors/detail.go:368 +0x2d
// 		github.com/gota33/errors.Annotate(0x1051780, 0x1257e60, 0xc00010fc00, 0x5, 0x5, 0xc00010fba8, 0x10)
// 			/home/user/github/gota33/errors/errors.go:79 +0x97
// 		github.com/gota33/errors.ExampleAnnotate()
// 			/home/user/github/gota33/errors/example_test.go:10 +0x251
// 		testing.runExample(0xfe589a, 0xf, 0xfff6c0, 0xfead08, 0x1a, 0x0, 0x0)
// 			/home/user/go/src/testing/run_example.go:63 +0x222
// 		testing.runExamples(0xc00010fed0, 0x120aee0, 0x3, 0x3, 0x0)
// 			/home/user/go/src/testing/example.go:44 +0x185
// 		testing.(*M).Run(0xc000114100, 0x0)
// 			/home/user/go/src/testing/testing.go:1419 +0x27d
// 		main.main()
// 			_testmain.go:71 +0x145
//
// detail[1]:
// 	type: "type.googleapis.com/google.rpc.RequestInfo"
// 	request_id: "<uuid>"
// 	serving_data: ""
// detail[2]:
// 	type: "type.googleapis.com/google.rpc.LocalizedMessage"
// 	local: "en-US"
// 	message: "Background task timeout"
// detail[3]:
// 	type: "type.googleapis.com/google.rpc.LocalizedMessage"
// 	local: "zh-CN"
// 	message: "后台任务超时"
 ```

## Encode error to JSON

``` go
// ...
// Use any JSON encoder you prefered
jEnc := json.NewEncoder(os.Stdout)
jEnc.SetIndent("  ", "  ")
enc := NewEncoder(jEnc)
// Hide DebugInfo before send to client
enc.Filters = []DetailFilter{HideDebugInfo}
_ = enc.Encode(err)

// Example Output:
// {
//    "error": {
//      "code": 504,
//      "message": "context deadline exceeded",
//      "status": "DEADLINE_EXCEEDED",
//      "details": [
//        {
//          "@type": "type.googleapis.com/google.rpc.RequestInfo",
//          "requestId": "\u003cuuid\u003e"
//        },
//        {
//          "@type": "type.googleapis.com/google.rpc.LocalizedMessage",
//          "local": "en-US",
//          "message": "Background task timeout"
//        },
//        {
//          "@type": "type.googleapis.com/google.rpc.LocalizedMessage",
//          "local": "zh-CN",
//          "message": "后台任务超时"
//        }
//      ]
//    }
//  }
```

## Decode error from JSON

### Decode manually

``` go
// import . "github.com/gota33/errors"
client := &http.Client{}
resp, err := client.Get("https://localhost:8080/users/1")
if err != nil {
    return err
}
defer resp.Body.Close()

// Decode from response body
r := strings.NewReader(resp.Body)
dec := NewDecoder(json.NewDecoder(r))

err := dec.Decode()

fmt.Println(Code(err))
// Output: 404 NOT_FOUND

```

### Decode automatically

``` go
// import . "github.com/gota33/errors"
client := http.Client{
    Transport: &RoundTripper{
        Parent: http.DefaultTransport,
    },
}

resp, err := client.Get("https://localhost:8080/users/1")
if err == nil {
    defer resp.Body.Close()
}

fmt.Println(Code(err))
// Output: 404 NOT_FOUND
```

