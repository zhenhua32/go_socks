## 简介

代码 fork 自 https://github.com/gwuhaolin/lightsocks
以及文档 https://github.com/gwuhaolin/blog/issues/12

## windows 上构建

```powershell
# 构建服务器端 for linux
$env:GOOS="linux"; go build .\cmd\server\

## 构建本地客户端 for windows
$env:GOOS="windows"; go build .\cmd\server\
```


