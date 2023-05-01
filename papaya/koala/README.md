# KOALA

### dynamic data structure

### TODOs
- new console
```go
package console

import (
  "github.com/mattn/go-colorable"
  "github.com/mattn/go-isatty"
	
)
```
- Any, Array, Map
```go
package typed

type Any interface{} // auto
type Array []interface{} // static
type Map map[string]interface{} // objective
```
- Collection
```go
package collection

type ArrayList struct {} // dynamic
type MapList struct {} // dynamic
```