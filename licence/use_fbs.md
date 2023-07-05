# Go 使用 FlatBuffers

### 编译 flatc 工具
- 将 flatbuffers.rar 文件解压
- cd flatbuffers/
- cmake -G "Unix Makefiles"
- make -j8 && make install
- flatc --version

#### Go 使用 FlatBuffers
  - 编译 scheme 文件为 go 文件，会生成 xxxx.go文件
- 示例如下：
```shell
$ cd fbs/
$ ls
lic_proto.fbs
$ flatc --go lic_proto.fbs
$ ls
proto           lic_proto.fbs
$ cd proto/
$ ls
AuthType.go     CentreInfo.go   LicenceProto.go LocalInfo.go    TimeInfo.go
```

- go get -v github.com/google/flatbuffers/go
```shell 
 go get -v github.com/google/flatbuffers/go
go: downloading github.com/google/flatbuffers v1.12.1
go: downloading github.com/google/flatbuffers v2.0.6+incompatible
github.com/google/flatbuffers/go
go get: added github.com/google/flatbuffers v2.0.6+incompatible

```
   