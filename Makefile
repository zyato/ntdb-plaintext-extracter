# 生成model
proto:
	protoc --go_out=./model/pb ./model/pb/msg_body.proto

# 安装依赖，目前只有protobuf，用于解析聊天消息体
install:
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
