# utils
go语言常用函数集，做为第三方包使用

## go交叉编译
GOOS：目标平台的操作系统（darwin、freebsd、linux、windows）
GOARCH：目标平台的体系架构（386、amd64、arm）
交叉编译不支持 CGO 所以要禁用它
### Mac 下编译 Linux 和 Windows 64位可执行程序
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
### Linux 下编译 Mac 和 Windows 64位可执行程序
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
### Windows 下编译 Mac 和 Linux 64位可执行程序
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go

## go条件编译
### 文件命名(文件名后缀)
filename_GOOS_GOARCH.go
myfunc.go
myfunc_linux_amd64.go
myfunc_linux_test.go
myfunc_darwin.go            苹果系统
myfunc_windows.go
### 构建标签(build tags) 编译标签
// +build linux,cgo darwin,cgo
// +build darwin dragonfly freebsd linux netbsd openbsd plan9 solaris
// +build !darwin !dragonfly !freebsd !linux !netbsd !openbsd !plan9 !solaris
// +build windows
