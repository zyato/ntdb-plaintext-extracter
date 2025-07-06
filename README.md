# 提取 QQNT 数据库 group_msg_table 中的聊天记录

## 前提
 
- 已经通过教程导出了 nt_msg.decrypt.db
- 教程链接：https://qq.sbcnm.top/decrypt/decode_db.html

## 用法
- main.go 文件中的 groupID 替换成预期qq群（当然也可以直接改 SQL）。
- main.go 文件中的 pbPath 替换成导出的 nt_msg.decrypt.db 本地路径
- 执行 `go run main.go`

