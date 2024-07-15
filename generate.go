//go:generate go run cmd/treenode/main.go -name=TreeNode -fields=Data/T -g=T/any -o=generic.go
//go:generate go run cmd/treenode/main.go -name=BoolNode -fields=Data/bool -o=bool.go
//go:generate go run cmd/treenode/main.go -name=ByteNode -fields=Data/byte -o=byte.go
//go:generate go run cmd/treenode/main.go -name=Complex64Node -fields=Data/complex64 -o=complex64.go
//go:generate go run cmd/treenode/main.go -name=Complex128Node -fields=Data/complex128 -o=complex128.go
//go:generate go run cmd/treenode/main.go -name=ErrorNode -fields=Data/error -o=error.go
//go:generate go run cmd/treenode/main.go -name=Float32Node -fields=Data/float32 -o=float32.go
//go:generate go run cmd/treenode/main.go -name=Float64Node -fields=Data/float64 -o=float64.go
//go:generate go run cmd/treenode/main.go -name=IntNode -fields=Data/int -o=int.go
//go:generate go run cmd/treenode/main.go -name=Int8Node -fields=Data/int8 -o=int8.go
//go:generate go run cmd/treenode/main.go -name=Int16Node -fields=Data/int16 -o=int16.go
//go:generate go run cmd/treenode/main.go -name=Int32Node -fields=Data/int32 -o=int32.go
//go:generate go run cmd/treenode/main.go -name=Int64Node -fields=Data/int64 -o=int64.go
//go:generate go run cmd/treenode/main.go -name=RuneNode -fields=Data/rune -o=rune.go
//go:generate go run cmd/treenode/main.go -name=StringNode -fields=Data/string -o=string.go
//go:generate go run cmd/treenode/main.go -name=UintNode -fields=Data/uint -o=uint.go
//go:generate go run cmd/treenode/main.go -name=Uint8Node -fields=Data/uint8 -o=uint8.go
//go:generate go run cmd/treenode/main.go -name=Uint16Node -fields=Data/uint16 -o=uint16.go
//go:generate go run cmd/treenode/main.go -name=Uint32Node -fields=Data/uint32 -o=uint32.go
//go:generate go run cmd/treenode/main.go -name=Uint64Node -fields=Data/uint64 -o=uint64.go
//go:generate go run cmd/treenode/main.go -name=UintptrNode -fields=Data/uintptr -o=uintptr.go
//go:generate go run cmd/treenode/main.go -name=StatusNode -fields=Status/S,Data/T -g=S/common.Enumer,T/any -o=status.go
package treenode
