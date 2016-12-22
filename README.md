# go-swf

[![Build
Status](https://travis-ci.org/Kelvyne/go-swf.svg?branch=master)](https://travis-ci.org/Kelvyne/go-swf)
[![Go Report Card](https://goreportcard.com/badge/github.com/kelvyne/go-swf)](https://goreportcard.com/report/github.com/kelvyne/go-swf)
[![Go Coverage](http://gocover.io/_badge/github.com/kelvyne/go-swf)](https://gocover.io/github.com/kelvyne/go-swf)

Package swf contains utilities to read Shockwave Flash Format files

### Documentation

See [here](https://godoc.org/github.com/kelvyne/go-swf)

### Example

```go
parser := swf.NewParser(r)
swfFile, err := parser.Parse()
fmt.Printf("Tags count : %v\n", len(swfFile.Tags))
```
