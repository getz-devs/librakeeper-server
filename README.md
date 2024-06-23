# librakeeper-server

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/search.proto
```

```shell
protoc --go-grpc_out=. librakeeper-server/pkg/pb/book.proto
```