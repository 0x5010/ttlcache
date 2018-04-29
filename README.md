This is a ttlcache in the [Go](http:golang.org).

[![Build Status](https://travis-ci.org/0x5010/ttlcache.png?branch=master)](https://travis-ci.org/0x5010/ttlcache)

Installation
-----------

	go get github.com/0x5010/ttlcache


Usage
-----------

Get and set with default ttlcache:
```go
// cache bytes
ttlcache.Set("key", []byte("value"), time.Duration(30*time.Second))

// get
b, ok := ttlcache.Get("key")
```
```go
// cache other with gob
type Test struct {
	V  [][]string
}

var test Test
...
var serialize bytes.Buffer
encoder := gob.NewEncoder(&serialize)
err := encoder.Encode(testStruct)
if err != nil {
	return nil, err
}
ttlcache.Set("key", serialize.Bytes(), time.Duration(30*time.Second))

// get
b, ok := ttlcache.Get("key")

t := &Test{}
decoder := gob.NewDecoder(bytes.NewReader(b))
err := decoder.Decode(t)
```
new cache:
```go
myCache := ttlcache.New(time.Duration(20 * time.Minute))

myCache.Set("key", []byte("value"), time.Duration(30*time.Second))

b, ok := myCache.Get("key")

```

