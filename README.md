# 导出网易云音乐 NCM 格式 与 qq音乐 qmc格式

## 简介

遍历目录找出网易云加密格式与qq音乐加密格式转化为普通格式

## 如何使用？


```shell
go build ./cmd/
cp cmd $GOBIN/musicdump
cp cmd/qmcdump $GOBIN/qmcdump
musicdump -i 输入目录 -o 输出目录

```
