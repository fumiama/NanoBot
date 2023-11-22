package nano

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// FileType 媒体类型
type FileType int

const (
	FileTypeImage = iota + 1 // png/jpg
	FileTypeVideo            // mp4
	FileTypeAudio            // silk
	FileTypeFile             // 暂不开放
)

func (ft FileType) String() string {
	switch ft {
	case FileTypeImage:
		return "图片"
	case FileTypeVideo:
		return "视频"
	case FileTypeAudio:
		return "语音"
	case FileTypeFile:
		return "文件"
	default:
		return "未知类型" + strconv.Itoa(int(ft))
	}
}

// FilePost QQ 富媒体消息发送请求参数
//
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/rich-text-media.html
type FilePost struct {
	Type       FileType `json:"file_type"`
	URL        string   `json:"url"`
	IsPositive bool     `json:"srv_send_msg"` // IsPositive
	// file_data		否	【暂未支持】
}

func (fp *FilePost) String() string {
	sb := strings.Builder{}
	sb.WriteString("[v2.")
	sb.WriteString(fp.Type.String())
	sb.WriteString("]")
	if fp.URL == "" {
		sb.WriteString("无链接")
	} else {
		sb.WriteString("链接: ")
		sb.WriteString(fp.URL)
	}
	return sb.String()
}

// PostFileToQQUser 发送文件到 QQ 用户的 openid
//
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/rich-text-media.html#%E5%8F%91%E9%80%81%E5%88%B0%E5%8D%95%E8%81%8A
func (bot *Bot) PostFileToQQUser(id string, content *FilePost) (*Message, error) {
	logrus.Infoln(getLogHeader(), "<= [Q]单:", id+",", content)
	return bot.postOpenAPIofMessage("/v2/users/"+id+"/files", "", WriteBodyFromJSON(content))
}

// PostFileToQQGroup 发送文件到 QQ 群的 openid
//
// https://bot.q.qq.com/wiki/develop/api-231017/server-inter/message/send-receive/rich-text-media.html#%E5%8F%91%E9%80%81%E5%88%B0%E7%BE%A4%E8%81%8A
func (bot *Bot) PostFileToQQGroup(id string, content *FilePost) (*Message, error) {
	logrus.Infoln(getLogHeader(), "<= [Q]群:", id+",", content)
	return bot.postOpenAPIofMessage("/v2/groups/"+id+"/files", "", WriteBodyFromJSON(content))
}

// MessageMedia used in MessagePost
type MessageMedia struct {
	FileInfo string `json:"file_info"`
}
