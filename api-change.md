## general

- Get2 are all renamed to Get and always returns tow value:
  a struct related type and a `bool` indicating if the value is found or not.

## package trie

For most cases, leave the first argument to be nil.
Thus this does not affect a lot.

```
- func NewCompactedTrie(c array.EltConverter) *CompactedTrie
+ func NewSlimTrie(m marshal.Marshaler, keys []string, values interface{}) (*SlimTrie, error) {
```

```
- CompactedTrie.SearchStringEqual(key string) interface{}
+ SlimTrie.Get(key string) interface{}, bool
```

Internal changes:

SlimTrie is now based on predefined array structure and do not need to defines
array itself.

```
type SlimTrie struct {
    Children array.ArrayU32
    Steps    array.ArrayU16
    Leaves   array.Array
}
```

## package array

Change uint32 to int32:

```
message Array32 {
-   uint32 Cnt             ;
+    int32 Cnt             ;
    repeated uint64 Bitmaps;
-   repeated uint32 Offsets;
+   repeated int32  Offsets;
    bytes  Elts            ;
}
```

Rename struct:

```
- array.CompactedArray
+ array.Array
```

API changes according to prototype changes:

```
- array.CompactedArray.Init(index []uint32, elts interface{}) error
+ array.Array.Init(index []int32, elts interface{}) error

// also Get now return two value: interface{} and bool
- array.CompactedArray.Get(index uint32) interface{}
+ array.Array.Get(index []int32) (interface{}, bool)

- array.CompactedArray.Has(index uint32) bool
+ array.Array.Has(index []int32) bool
```

Add new APIs:

```
array.New(indexes []int32, elts interface{}) (*array.Array, error)
```

Other added APIs are bound to `array.ArrayBase`:
`array.ArrayBase` provides underlying functionalities for fixed type array like
`ArrayU16` and auto-type `Array`.

```
+ array.ArrayBase.InitIndex(index []int32) error
+ array.ArrayBase.InitElts(elts interface{}, m marshal.Marshaler) (int, error)
+ array.ArrayBase.GetEltIndex(idx int32) (int32, bool)
+ array.ArrayBase.Has(idx int32) bool
+ array.ArrayBase.Init(idx int32, elts interface{}) error
+ array.ArrayBase.GetTo(idx int32, elts interface{}) bool
+ array.ArrayBase.GetBytes(idx int32) ([]byte, bool)
```

Adds predefined array data structures:

```
+ array.ArrayU16
+ array.ArrayU32
+ array.ArrayU64
```


## Converter is moved to marshal.Marshaler

Marshaler requires another method `GetSize(v interface{})` that Converter does
not require.

```
- array.Converter
+ marshal.Marshaler
```
