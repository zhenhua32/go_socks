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

## 运行方式

在服务器上运行 server 文件, 会打印出配置信息.

在本地运行 local.exe 文件, 第一次运行会失败, 但是会生成配置文件,
然后按照服务器显示的配置信息更新配置文件. 最后重新运行即可.

示例的配置信息如下:

```json
{
 "listen": ":7448",
 "remote": "[2001:19f0:5:3d3f:5400:01ff:fe4a:a55f]:41757",
 "password": "mECSTEmIUj/BWdOmgqmwH6f5xrIikAFP78L04ljPC/GeNxKtSoXQvwohtsvh/iYZG2z3aXBj6LhiBrERfX6GNk2Pt5tToaRoHb1b+g++rOgQOJvGR7VFZCKtiLQ/+ub74c9jPHTtdVyZG6ONad/5ZflVz4Zfws8107MplRAC4jeHVFjg6l7fWoLekHbeVeWuxHILOrdH+KE9k8jB6NPt131dHabvsEcno1MQ=="
}
```

配置文件的名字是 `.socks.json`, 位置应该是在用户主目录下.
listen 是本地监听的地址, remote 是服务器的地址和端口, password 是密码.

注意, 如果是 ipv6 的地址, 需要使用 `[]` 包含地址.

别忘记填写端口号.
