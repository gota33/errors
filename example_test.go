package errors

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func ExampleAnnotate() {
	err := Annotate(
		context.DeadlineExceeded,
		DeadlineExceeded,
		StackTrace("heavy job"),
		RequestInfo{RequestId: "<uuid>"},
		LocalizedMessage{Local: "en-US", Message: "Background task timeout"},
		LocalizedMessage{Local: "zh-CN", Message: "后台任务超时"},
	)

	fmt.Printf("%+v", err)
	// Example Output:
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
}

func ExampleEncoder() {
	err := Annotate(
		context.DeadlineExceeded,
		DeadlineExceeded,
		StackTrace("heavy job"),
		RequestInfo{RequestId: "<uuid>"},
		LocalizedMessage{Local: "en-US", Message: "Background task timeout"},
		LocalizedMessage{Local: "zh-CN", Message: "后台任务超时"},
	)

	jEnc := json.NewEncoder(os.Stdout)
	jEnc.SetIndent("  ", "  ")
	enc := NewEncoder(jEnc)
	enc.Mappers = []DetailMapper{HideDebugInfo}
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
}

func ExampleCode() {
	err := WithNotFound(sql.ErrNoRows, ResourceInfo{
		ResourceType: "pet.com/pet.v1.Cat",
		ResourceName: "cat123@pet.com",
		Owner:        "user123",
		Description:  "cat123 not found",
	})
	code := Code(err)
	fmt.Println(code)
	// Output: 404 NOT_FOUND
}

func ExampleDetails() {
	client := &http.Client{}
	client.Get("")
	err := WithNotFound(sql.ErrNoRows, ResourceInfo{
		ResourceType: "pet.com/pet.v1.Cat",
		ResourceName: "cat123@pet.com",
		Owner:        "user123",
		Description:  "cat123 not found",
	})

	for _, detail := range Details(err) {
		fmt.Printf("%+v", detail)
	}
	// Output:
	// type: "type.googleapis.com/google.rpc.ResourceInfo"
	// resource_type: "pet.com/pet.v1.Cat"
	// resource_name: "cat123@pet.com"
	// owner: "user123"
	// description: "cat123 not found"
}