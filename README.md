# auto changelog

## Feature

根据最近的tag自动生成changelog并写入changelog文件。

## Install

```bash
go install github.com/aak1247/gchangelog@latest
```

## Usage

```bash
gchangelog [-options -arguments]
  -f string
        文件名 (default "changelog.md")
  -mr
        记录mr日志
  -p string
        仓库地址
  -skip string
        跳过消息 (default "skip"), 匹配跳过消息的commit message不会生成到changelog中,多个消息用","分割
```

## 示例

见 [changelog.md](changelog.md) , 里面的内容是自动生成的