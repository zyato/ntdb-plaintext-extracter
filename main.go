package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	sql "github.com/FloatTech/sqlite"

	"github.com/zyato/ntdb-plaintext-extracter/model"
)

const (
	debug   = false                      // NOTE：debug模式会输出到控制台，而不是output.csv
	groupID = 0                          // NOTE: 这里替换成预期qq群
	dbPath  = "./yato/nt_msg.decrypt.db" // NOTE：这里替换成已经解密的nt_msg db路径

	table = "group_msg_table"
)

type Service struct {
	DB     sql.Sqlite
	File   *os.File
	Writer *csv.Writer
}

func NewService(dbPath string) Service {
	db := sql.New(dbPath)
	err := db.Open(time.Hour)
	if err != nil {
		panic("Open" + err.Error())
	}
	file := os.Stdout
	if !debug {
		file, err = os.OpenFile("output.csv", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic("OpenFile: " + err.Error())
		}
	}
	writer := csv.NewWriter(file)
	// 写入 UTF-8 BOM，防止中文在 Excel 中显示乱码
	// BOM 是字节序标记：0xEF,0xBB,0xBF
	_, err = file.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		panic("Write" + err.Error())
	}

	return Service{
		DB:     db,
		File:   file,
		Writer: writer,
	}
}

func (srv *Service) Run(condition string) error {
	// 写入标题
	header := []string{"序号", "发送人", "发送时间", "消息内容"}
	if err := srv.Writer.Write(header); err != nil {
		panic("Write: " + err.Error())
	}

	iter := &model.Row{}
	return srv.DB.FindFor(table, iter, condition, func() error {
		srv.Writer.Flush()
		err := srv.Writer.Write(iter.MakeCSVRecord())
		if err != nil {
			return fmt.Errorf("Write: " + err.Error())
		}
		if debug {
			srv.Writer.Flush()
		}
		return nil
	})
}

func (srv *Service) Close() {
	srv.Writer.Flush()
	_ = srv.File.Close()
	_ = srv.DB.Close()
}

func main() {
	// 获取指定群组全量的聊天信息
	srv := NewService(dbPath)
	defer srv.Close()

	// 存在部分消息40021列为空、40027/40030不为空的情况，这里需要OR一下
	condition := fmt.Sprintf("WHERE [40027]=%d OR [40030]=%d OR [40021]='%d' ORDER BY [40003] DESC", groupID, groupID, groupID)
	err := srv.Run(condition)
	if err != nil {
		panic("Run: " + err.Error())
	}
}
