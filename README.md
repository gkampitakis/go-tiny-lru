# Go Tiny LRU

This go module is a replication of [tiny-lru](https://www.npmjs.com/package/tiny-lru) written in js.

## Description

Trying to learn basic concepts of golang, implement tests and benchmarks.

## Usage

```bash
go get -u github.com/gkampitakis/go-tiny-lru
```
Then you can use it in your code 

```go

import "github.com/gkampitakis/go-tiny-lru"

LRU.New(10,1000)
// First parameter is the max items that LRU can store
// and the second is the ttl number. 
//You can set them both to zero for no max capacity and no ttl.
```

## Methods

### New

Creates a new LRU instance. 
<br>
**Parameters:** `max int`,`ttl int64`. max and ttl can be set to zero for 
specifying "infinite" capacity and no expiration time. 
<br>
**Returns:** `error,*LRU`

### Clear

Clears LRU and resets all values. 
<br>
**Parameters:** -
<br>
**Returns:** `*LRU`

### Delete

Deletes key from LRU if existent.
<br>
**Parameters:** `key string`
<br>
**Returns:** `*LRU`

### Keys

Returns a slice containing all keys in LRU.
<br>
**Note:** the order of insertion is not guaranteed in the output slice as internally LRU uses maps for storing items.
<br>
**Parameters:** -
<br>
**Returns:** `[]string`

### Set

Adds a new item in LRU.
<br>
**Parameters:** `key string`,`value interface{}`
<br>
**Returns:** `*LRU`

### Get

Returns the value of a key if existent
<br>
**Parameters:** `key string`
<br>
**Returns:** `interface{} || nil`