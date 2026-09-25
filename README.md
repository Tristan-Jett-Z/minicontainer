# MiniContainer
一个使用 Go 从 Linux Namespace、Cgroups、rootfs 等底层机制实现的微型容器运行时

## 已实现
- PID Namespace
- Mount Namespace
- UTS Namespace
- Cgroups
- rootfs
- /proc
- chroot
- exec

## 环境
- Linux
- Go

## 运行
首先需要在minicontainer里准备一个rootfs:
minicontainer/
├── main.go
├── go.mod
├── .gitignore
└── rootfs/

然后：
go build -o minicontainer
sudo ./minicontainer

## 注意
rootfs/ 未包含在 Git 仓库中，需要自行准备。

## 项目目的
通过亲手实现一个最小容器运行时，理解 Linux Namespace、Cgroups、rootfs、进程隔离等容器底层机制
