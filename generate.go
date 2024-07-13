//go:generate go run cmd/treenode/main.go -type=TreeNode -fields=Data/T -g=T/any -output=generic.go
//go:generate go run cmd/treenode/main.go -type=Bool -fields=Data/bool -output=bool.go
//go:generate go run cmd/treenode/main.go -type=Byte -fields=Data/byte -output=byte.go
//go:generate go run cmd/treenode/main.go -type=Complex64 -fields=Data/complex64 -output=complex64.go
//go:generate go run cmd/treenode/main.go -type=Complex128 -fields=Data/complex128 -output=complex128.go
//go:generate go run cmd/treenode/main.go -type=Error -fields=Data/error -output=error.go
//go:generate go run cmd/treenode/main.go -type=Float32 -fields=Data/float32 -output=float32.go
//go:generate go run cmd/treenode/main.go -type=Float64 -fields=Data/float64 -output=float64.go
//go:generate go run cmd/treenode/main.go -type=Int -fields=Data/int -output=int.go
//go:generate go run cmd/treenode/main.go -type=Int8 -fields=Data/int8 -output=int8.go
//go:generate go run cmd/treenode/main.go -type=Int16 -fields=Data/int16 -output=int16.go
//go:generate go run cmd/treenode/main.go -type=Int32 -fields=Data/int32 -output=int32.go
//go:generate go run cmd/treenode/main.go -type=Int64 -fields=Data/int64 -output=int64.go
//go:generate go run cmd/treenode/main.go -type=Rune -fields=Data/rune -output=rune.go
//go:generate go run cmd/treenode/main.go -type=String -fields=Data/string -output=string.go
//go:generate go run cmd/treenode/main.go -type=Uint -fields=Data/uint -output=uint.go
//go:generate go run cmd/treenode/main.go -type=Uint8 -fields=Data/uint8 -output=uint8.go
//go:generate go run cmd/treenode/main.go -type=Uint16 -fields=Data/uint16 -output=uint16.go
//go:generate go run cmd/treenode/main.go -type=Uint32 -fields=Data/uint32 -output=uint32.go
//go:generate go run cmd/treenode/main.go -type=Uint64 -fields=Data/uint64 -output=uint64.go
//go:generate go run cmd/treenode/main.go -type=Uintptr -fields=Data/uintptr -output=uintptr.go
//go:generate go run cmd/treenode/main.go -type=StatusNode -fields=Status/S,Data/T -g=S/uc.Enumer,T/any -output=status.go
package treenode
