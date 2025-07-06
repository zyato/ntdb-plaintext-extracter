package model

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/zyato/ntdb-plaintext-extracter/model/pb/pbmodel"
)

type MsgBody pbmodel.MsgBody

func (m *MsgBody) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, (*pbmodel.MsgBody)(m))
}

func (m *MsgBody) String() string {
	sb := strings.Builder{}
	for _, elem := range m.Elements {
		sb.WriteString(((*MsgElem)(elem)).String())
	}
	return sb.String()
}

type MsgElem pbmodel.MsgElem

func (m *MsgElem) String() string {
	var mp = map[int32]func() string{
		1:  m.stringText,
		2:  m.stringPic,
		3:  m.stringFile,
		4:  m.stringPTT,
		5:  m.stringVideo,
		6:  m.stringFace,
		7:  m.stringReply,
		8:  m.stringGrayTip,
		9:  m.stringWallet,
		10: m.stringARK,
		11: m.stringMarketFace,
		14: m.stringXML,
		17: m.stringInLineKeyBoard,
		21: m.stringAVRecord,
		27: m.stringFaceBubble,
		28: m.stringShareLocation,
	}
	f := mp[int32(m.ElementType)]
	if f != nil {
		return f()
	}
	return fmt.Sprintf("stringDefault_%d", m.ElementType)
}

func (m *MsgElem) stringText() string {
	return m.Content
}

func (m *MsgElem) stringPic() string {
	picType := "null"
	if m.PicType == 1000 {
		picType = "static"
	}
	if m.PicType == 2000 {
		picType = "gif"
	}
	return "stringPic_" + picType
}

func (m *MsgElem) stringFile() string {
	return "stringFile_" + m.FileName
}

func (m *MsgElem) stringPTT() string {
	return "stringPTT_" + m.VoiceText
}

func (m *MsgElem) stringVideo() string {
	return "stringVideo_" + m.Content
}

func (m *MsgElem) stringFace() string {
	return "stringFace_" + m.FaceText
}

func (m *MsgElem) stringReply() string {
	return "reference[" + m.ReplyContent + "]" + m.Content
}

func (m *MsgElem) stringGrayTip() string {
	return "stringGrayTip_" + m.Content
}

func (m *MsgElem) stringWallet() string {
	return "stringWallet_" + m.Content
}

func (m *MsgElem) stringARK() string {
	return "stringARK_" + m.Content
}

func (m *MsgElem) stringMarketFace() string {
	return "stringMarketFace_" + m.Content
}

func (m *MsgElem) stringMarkDown() string {
	return "stringMarkDown_" + m.Content
}

func (m *MsgElem) stringXML() string {
	return "stringXML_" + m.Content
}

func (m *MsgElem) stringInLineKeyBoard() string {
	return "stringInLineKeyBoard_" + m.Content
}

func (m *MsgElem) stringAVRecord() string {
	return "stringAVRecord_" + m.Content
}

func (m *MsgElem) stringFaceBubble() string {
	return "stringFaceBubble_" + m.Content
}

func (m *MsgElem) stringShareLocation() string {
	return "stringShareLocation_" + m.Content
}
