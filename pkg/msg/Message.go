package msg

import (
	"strings"
	"time"
)

type Message struct {
	ID          string
	MsgID       string
	Hash        string
	Area        string
	From        string
	To          string
	Subject     string
	Content     string
	UnixTime    int64
	DateWritten *time.Time
	ViewCount   int
}

func NewMessage() *Message {
	msg := new(Message)
	return msg
}

func (self *Message) GetContent() string {
	return self.Content
}

func (self *Message) SetID(id string) {
	self.ID = id
}

func (self *Message) SetMsgID(msgId string) {
	self.MsgID = msgId
}

func (self *Message) SetArea(area string) {
	self.Area = strings.TrimRight(area, "\x00")
}

func (self *Message) SetFrom(from string) {
	self.From = strings.TrimRight(from, "\x00")
}

func (self *Message) SetTo(to string) {
	self.To = strings.TrimRight(to, "\x00")
}

func (self *Message) SetSubject(subject string) {
	self.Subject = strings.TrimRight(subject, "\x00")
}

func (self *Message) SetContent(content string) {
	self.Content = strings.TrimRight(content, "\x00")
}

func (self *Message) SetUnixTime(unixTime int64) {
	self.UnixTime = unixTime
	tm := time.Unix(unixTime, 0)
	self.DateWritten = &tm
}

func (self *Message) SetTime(ptm *time.Time) {
	self.DateWritten = ptm
	if ptm != nil {
		var tm time.Time = *ptm
		self.UnixTime = tm.Unix()
	}
}

func (self *Message) SetViewCount(count int) {
	self.ViewCount = count
}

func (self *Message) SetMsgHash(hash string) {
	self.Hash = hash
}
