# go-swf

[![GoDoc](https://godoc.org/github.com/Kelvyne/swf?status.svg)](https://godoc.org/github.com/Kelvyne/swf)
[![Build
Status](https://travis-ci.org/Kelvyne/swf.svg?branch=master)](https://travis-ci.org/Kelvyne/swf)
[![Go Report Card](https://goreportcard.com/badge/github.com/kelvyne/swf)](https://goreportcard.com/report/github.com/kelvyne/swf)
[![Go Coverage](http://gocover.io/_badge/github.com/kelvyne/swf)](https://gocover.io/github.com/kelvyne/swf)

Package swf contains utilities to read Shockwave Flash Format files

### Documentation

See [here](https://godoc.org/github.com/kelvyne/go-swf)

### Example

```go
parser := swf.NewParser(r)
swfFile, err := parser.Parse()
fmt.Printf("Tags count : %v\n", len(swfFile.Tags))
```
