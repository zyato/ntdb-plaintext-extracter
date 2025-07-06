package model

import (
	"fmt"
	"strconv"
	"time"
)

// https://qq.sbcnm.top/view/db_file_analysis/nt_msg.db.html#_40600%E5%80%BC%E8%A7%A3%E8%AF%BB
type Row struct {
	MsgID    int64 `db:"40001"` // 消息id。具有唯一性
	Useless1 int64 `db:"40002"` // 消息随机值。用于对消息去重
	MsgSeq   int64 `db:"40003"` // 群聊消息id。在每个聊天中一次递增
	ChatType int64 `db:"40010"` // 聊天类型。私聊1、群聊2、频道4、公众号103、企业客服102、临时会话100

	// 无消息0（消息损坏？多见于已退出群聊且时间久远）、1消息空白（msgid存在，应该是没加载出来）、text文本消息2、群文件3、我的聊天记录里没有4~大佬带带＞︿＜
	// 系统（灰字）消息5、语音消息6、视频文件7、合并转发消息8、回复类型消息9、红包10、应用消息11
	MsgType int64 `db:"40011"` // 疑似用于区分消息类型。
	// 非常规text消息0、普通文本消息1、群文件其他类型消息1、图片消息2、群文件图片消息2、群公告3、群文件视频消息4、撤回消息提醒4、群文件音频消息8、原创表情包8
	// 设精消息11、拍一拍消息12、群文件docx消息16、平台文本消息32、群文件pptx消息32、回复类型消息33、群文件xlsx消息64
	// 存在链接161、群文件zip消息512、群文件exe消息2048、表情消息4096
	SubMsgType int64 `db:"40012"` // 疑似用于区分protobuf消息类型。1 is normal
	/* 40011与40012组合可判断消息类型。以下是一些常见的组合
	由于优先级问题（特别是2类别的信息），部分消息不满足以下规则
	空消息：0，0
	已撤回消息：1，0
	普通文本类消息：2，1
	图片消息：2，2
	只带图片的纯文本消息：2，3
	纯表情消息：2，16
	带表情的纯文本消息：2，17
	带图片带表情的纯文本消息：2，19
	纯链接消息：2，129
	带表情链接消息：2，145
	机器人消息：2，577
	机器人Markdown消息：2，65
	@消息：2，35
	回复引用消息（不带表情）：2，33
	回复引用消息（带表情）：2，49
	收藏表情：2，2
	收藏表情包：2，4096
	群文件其他类型消息：3，1
	群文件图片（png，jpg）消息：3，2
	群文件视频消息：3，4
	群文件音频（mp3，flac）消息：3，8
	群文件docx消息：3，16
	群文件pptx消息：3，32
	群文件xlsx消息：3，64
	群文件zip消息：3，512
	群文件exe消息：3，2048
	拍一拍消息：5，12
	撤回消息提醒：5，4
	amr语音文件消息：6，0
	视频文件消息：7，0
	合并转发消息：8，0
	回复消息：9，33
	回复带图片消息（无@）：9，34
	回复带图片消息（有@）：9，35
	回复带图片带@：9，35
	回复卡片引用消息：9，49
	带表情回复：9，49
	带表情带图片带@：9，51
	回复存在链接的消息：9，161
	红包：10，0
	应用消息（如小程序）：11，0
	群公告：11，3
	表情包：17，8
	原创表情：17，8
	*/

	SendType         int64  `db:"40013"` // 发送标志。本机发送的消息1、其他客户端发送的消息2、别人发的消息0、转发消息5、在已退出后者被封禁的消息中为当日整点时间戳
	SenderUID        string `db:"40020"` // nt_uid。对应nt_uid_mapping_table
	Useless2         int64  `db:"40026"`
	PeerUID          string `db:"40021"` // 会话id。
	PeerUin          int64  `db:"40027"` // 会话id。
	Useless3         int64  `db:"40040"`
	SendStatus       int64  `db:"40041"` // 发送状态。成功2、发送被阻止0（比如不是对方好友）、尚未发送陈工1（比如网络问题）、消息被和谐3
	MsgTime          int64  `db:"40050"` // 发送消息时的完整时间戳。UTC+8:00
	Useless4         int64  `db:"40052"`
	SenderMemberName string `db:"40090"` // 发送者群名片。旧版qq迁移数据中格式为name(12345)或者name<i@example.com>，QQNT中为群名片
	SenderNickName   string `db:"40093"` // 发送着昵称。旧版qq此字段为空，QQNT中未设置群名片时才有此字段
	MsgBody          []byte `db:"40800"` // 聊天消息。文档中是40080，应该写错了。结构最为复杂，尚未解析完全
	MsgAdditional    []byte `db:"40900"` // 补充消息。不同情况下存在不一样的数据（以MsgType-40011列为区分）。MsgType=8时存储转发聊天的缓存、MsgType=9时存储引用的消息
	Useless5         int64  `db:"40105"`
	Useless6         int64  `db:"40005"` // 不清楚。只知道自己发的消息一定概率存在数值，正常情况为0
	UselessTimestamp int64  `db:"40058"` // 当日0时整的时间戳格式。时区为GMT+0800
	Useless7         int64  `db:"40006"` // 不清楚用处。elem id?
	AtStatus         int64  `db:"40100"` // @状态。有人@我为6；有人@他人为2；此条消息不包含@为0
	// 当40600（16进制）值为14 00时，为回复消息
	// 此时：40100的值：
	// 为6代表有人回复自己，为2代表他人回复他人
	// 当40600（16进制）值为c2e91304a8d114****时（不唯一），为撤回消息
	Useless8           []byte `db:"40600"` // 状态标志？protobuf格式
	DissolvedFlag      int64  `db:"40060"` // 已退出或者已解散的群聊标志
	RepliedMsgID       int64  `db:"40850"` // 回复消息序号。该消息所回复的消息的序号。todo 确认对应上面的哪个值
	Useless9           int64  `db:"40851"`
	Useless10          []byte `db:"40601"` // blob
	Useless11          []byte `db:"40801"` // 不清楚。是个protobuf
	Useless12          []byte `db:"40605"` // blob
	GroupID            int64  `db:"40030"` // QQNT保存的群号
	SenderQQ           int64  `db:"40033"` // QQNT保存的发送者qq号
	ExpressInfo        []byte `db:"40062"` // 存储详细表态信息（包括表态表情和表态数量）。其数字与QQBOT中表情编号对应（超级表情不在此列表中）
	ExpressEmojiTotal1 *int64 `db:"40083"` // 表态表情数量总和
	ExpressEmojiTotal2 *int64 `db:"40084"` // 表态表情数量总和
	Useless13          *int64 `db:"40008"`
	Useless14          *int64 `db:"40009"`
}

var (
	// 旧的记录没有SenderMemberName和SendNickName，这里存储新消息中使用的名字
	// 需要是按时间倒序查询，这里才有用
	qqToName = map[int64]string{}
)

func (r *Row) MakeCSVRecord() []string {
	pb := &MsgBody{}
	err := pb.Unmarshal(r.MsgBody)
	if err != nil {
		panic("Unmarshal" + err.Error())
	}
	memberName := r.SenderMemberName // 群昵称
	if memberName == "" {
		memberName = r.SenderNickName // 昵称
	}
	if memberName == "" {
		memberName = qqToName[r.SenderQQ] // 使用过的昵称
	}
	qqToName[r.SenderQQ] = memberName
	// header := []string{"序号", "发送人", "发送时间", "消息内容"}
	return []string{
		strconv.FormatInt(r.MsgSeq, 10),
		fmt.Sprintf("%s[%d]", memberName, r.SenderQQ),
		time.Unix(r.MsgTime, 0).Format("2006-01-02 15:04:05"),
		pb.String(),
	}
}
