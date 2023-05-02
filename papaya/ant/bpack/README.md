# Ant/BPack (Bundle Pack)

## bundle data assets into data.go

```go
// example mock file
file, err := bpack.OpenFile("/data/swag/ui.css") // fallback read in PATH = papaya/ant/bpack/data:data

// example read packet
file, err := bpack.OpenPacket("/data/swag/ui.css")
```