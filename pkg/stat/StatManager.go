package stat

import (
	"database/sql"
	"fmt"
	"github.com/vit1251/golden/pkg/storage"
	"log"
	"time"
)

type StatManager struct {
	conn *sql.DB
}

type Stat struct {
	TicReceived      int
	TicSent          int
	EchomailReceived int
	EchomailSent     int
	NetmailReceived  int
	NetmailSent      int

	PacketReceived   int
	PacketSent       int

	MessageReceived  int
	MessageSent      int

	SessionIn        int
	SessionOut       int
}

func NewStatManager(sm *storage.StorageManager) *StatManager {
	statm := new(StatManager)
	statm.conn = sm.GetConnection()
	statm.createStat()
	return statm
}

func (self *StatManager) RegisterInFile(filename string) (error) {

	self.createStat()

	query1 := "UPDATE `stat` SET `statFileRXcount` = `statFileRXcount` + 1 WHERE `statDate` = ?"
	statDate := self.makeToday()
	self.conn.Exec(query1, statDate)

	return nil
}

func (self *StatManager) RegisterOutFile(filename string) (error) {
	self.createStat()
	return nil
}

type SummaryRow struct {
	Date string
	Value int
}

func (self *StatManager) GetStatRow(statDate string) (*Stat, error) {

	var result *Stat = new(Stat)

	/* Step 2. Start SQL transaction */
	ConnTransaction, err := self.conn.Begin()
	if err != nil {
		return nil, err
	}

	sqlStmt := "SELECT `statMessageRXcount`, `statMessageTXcount`, `statSessionIn`, `statSessionOut`, `statFileRXcount`, `statFileTXcount`, `statPacketIn`, `statPacketOut` FROM `stat` WHERE `statDate` = $1"
	log.Printf("sql = %q echoTag = %q", sqlStmt, statDate)
	rows, err1 := ConnTransaction.Query(sqlStmt, statDate)
	if err1 != nil {
		return nil, err1
	}
	defer rows.Close()
	for rows.Next() {

		var statMessageInCount int64
		var statMessageOutCount int64
		var statSessionInCount int64
		var statSessionOutCount int64
		var statFileInCount int64
		var statFileOutCount int64
		var statPacketInCount int64
		var statPacketOutCount int64

		err2 := rows.Scan(&statMessageInCount, &statMessageOutCount, &statSessionInCount, &statSessionOutCount, &statFileInCount, &statFileOutCount, &statPacketInCount, &statPacketOutCount)
		if err2 != nil{
			return nil, err2
		}

		result.MessageReceived = int(statMessageInCount)
		result.MessageSent = int(statMessageOutCount)
		result.SessionIn = int(statSessionInCount)
		result.SessionOut = int(statSessionOutCount)
		result.TicReceived = int(statFileInCount)
		result.TicSent = int(statFileOutCount)
		result.PacketReceived = int(statPacketInCount)
		result.PacketSent = int(statPacketOutCount)

	}

	ConnTransaction.Commit()

	return result, nil
}

func (self *StatManager) GetStat() (*Stat, error) {
	statDate := self.makeToday()
	stat, err := self.GetStatRow(statDate)
	return stat, err
}

func (self *StatManager) RegisterInPacket() error {
	self.createStat()
	query1 := "UPDATE `stat` SET `statPacketIn` = `statPacketIn` + 1 WHERE `statDate` = ?"
	statDate := self.makeToday()
	self.conn.Exec(query1, statDate)
	return nil
}

func (self *StatManager) RegisterOutPacket() error {
	self.createStat()
	query1 := "UPDATE `stat` SET `statPacketOut` = `statPacketOut` + 1 WHERE `statDate` = ?"
	statDate := self.makeToday()
	self.conn.Exec(query1, statDate)
	return nil
}

func (self *StatManager) RegisterDupe() error {
	self.createStat()
	return nil
}

func (self *StatManager) RegisterInMessage() error {
	self.createStat()
	query1 := "UPDATE `stat` SET `statMessageRXcount` = `statMessageRXcount` + 1 WHERE `statDate` = ?"
	statDate := self.makeToday()
	self.conn.Exec(query1, statDate)
	return nil
}

func (self *StatManager) RegisterOutMessage() error {
	self.createStat()
	query1 := "UPDATE `stat` SET `statMessageTXcount` = `statMessageTXcount` + 1 WHERE `statDate` = ?"
	statDate := self.makeToday()
	self.conn.Exec(query1, statDate)
	return nil
}

func (self *StatManager) makeToday() string {
	currentTime := time.Now()
	//result := currentTime.Format("2006-01-02")
	result := fmt.Sprintf("%04d-%02d-%02d", currentTime.Year(), currentTime.Month(), currentTime.Day())
	return result
}

func (self *StatManager) createStat() {
	query1 := "INSERT INTO `stat` (`statDate`) VALUES ( ? )"
	statDate := self.makeToday()
	log.Printf("Create stat on %s", statDate)
	self.conn.Exec(query1, statDate)
}

func (self *StatManager) RegisterInSession() error {

	self.createStat()

	query1 := "UPDATE `stat` SET `statSessionIn` = `statSessionIn` + 1 WHERE `statDate` = ?"
	statDate := self.makeToday()
	self.conn.Exec(query1, statDate)

	return nil
}

func (self *StatManager) RegisterOutSession() error {

	self.createStat()

	query1 := "UPDATE `stat` SET `statSessionOut` = `statSessionOut` + 1 WHERE `statDate` = ?"
	statDate := self.makeToday()
	self.conn.Exec(query1, statDate)

	return nil
}
