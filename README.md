# go-swf

[![Build Status](https://travis-ci.com/Kelvyne/go-swf.svg?token=4hcDvc75wyCvjsysDuCx&branch=master)](https://travis-ci.com/Kelvyne/go-swf)

Package swf contains utilities to read Shockwave Flash Format files

### Documentation

See [here](https://godoc.org/github.com/kelvyne/go-swf)

### Example

```go
parser := swf.NewParser(r)
swfFile, err := parser.Parse()
fmt.Printf("Tags count : %v\n", len(swfFile.Tags))
```