package db

import (
	//	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//向用户信息表插入记录
func MemberInfoInsert(member MemberInfo) bool {
	//插入数据
	stmt, err := db.Prepare("INSERT INTO membersInfo(ip, clientPort, peerPort, name, id) VALUES(?,?,?,?,?)")
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	res, err := stmt.Exec(member.IP, member.ClientPort, member.PeerPort, member.Name, member.ID)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	_, err = res.LastInsertId()
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	return true
}

//取出用户信息表的所有记录
func MembersInfoQuery() ([]MemberInfo, error) {
	rows, err := db.Query("SELECT * FROM membersinfo")
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var membersInfo []MemberInfo = make([]MemberInfo, 0)
	for rows.Next() {
		var seq int
		var ip string
		var clientPort string
		var peerPort string
		var name string
		var id string
		err := rows.Scan(&seq, &ip, &clientPort, &peerPort, &name, &id)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}

		data := MemberInfo{
			IP:         ip,
			ClientPort: clientPort,
			PeerPort:   peerPort,
			Name:       name,
			ID:         id,
		}
		membersInfo = append(membersInfo, data)
	}

	return membersInfo, nil
}

func MemberInfoDelete(member MemberInfo) bool {
	//插入数据
	stmt, err := db.Prepare("DELETE FROM membersInfo WHERE ip=? AND clientPort=?")
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	res, err := stmt.Exec(member.IP, member.ClientPort)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	_, err = res.LastInsertId()
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	return true
}
