# Ant/BPack (Bundle Pack)

## bundle data assets into data.go

```go
// example implementation
// if not found, check in real path
file, err := bpack.OpenFile("/data/swag/ui.css", bpack.ReadOnly)
file, err := bpack.OpenFile("/data/swag/ui.css", bpack.TemporaryFile)
```