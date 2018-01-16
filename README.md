# Databox bindings for Go

Go wrapper for [Databox](http://databox.com) - Mobile Executive Dashboard.

[![Build Status][BuildStatus-Image]][BuildStatus-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]

## Installation

```bash
go install github.com/databox/databox-go
go get github.com/databox/databox-go # or this. :)
```

## Usage
```go
package main

import (
	databox "github.com/databox/databox-go"
	"fmt"
)

func main() {
	client := databox.NewClient("<push token>")

	if status, _ := client.Push(&databox.KPI{
		Key:	"temp.ny",
		Value: 	52.0,
		Date: 	"2015-01-01 09:00:00",
	}); len(status.Errors) == 0 {
		fmt.Println("Inserted.")
	}

	if data, _ := client.LastPush(); data != nil {
		fmt.Println(string(data))
	}

	// Additional attributes
	var attributes = make(map[string]interface{})
	attributes["test.number"] = 10
	attributes["test.string"] = "Oto Brglez"

	if status, _ := client.Push(&KPI{
		Key: "testing.this",
		Value: 10.0,
		Date: time.Now().Format(time.RFC3339),
		Attributes: attributes,
	}); len(status.Errors) > 0 {
		fmt.Println("This status contains errors")
	}

	// Retriving last push
	lastPush, err := client.LastPush()
	if err != nil {
		fmt.Println("Error was raised", err)
	}

	fmt.Println("Number of errors in last push", lastPush.NumberOfErrors)
}
```

## Author
- [Oto Brglez](https://github.com/otobrglez)

## Contributions
- [Caitlin Elfring](https://github.com/caitlin615)

[BuildStatus-Url]: https://travis-ci.org/databox/databox-go
[BuildStatus-Image]: https://travis-ci.org/databox/databox-go.svg
[ReportCard-Url]: http://goreportcard.com/report/databox/databox-go
[ReportCard-Image]: http://goreportcard.com/badge/databox/databox-go
