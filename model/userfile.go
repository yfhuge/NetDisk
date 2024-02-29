package model

import (
	"filestore-server/global"
	log "github.com/sirupsen/logrus"
	"time"
)

type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// OnUserFileUploadFinished 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := global.DB.GetConn().Prepare("insert ignore into tbl_user_file(`user_name`, `file_sha1`, `file_name`, `file_size`, `upload_at`) values(?,?,?,?,?)")
	if err != nil {
		log.Error("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		log.Error("Failed to insert, err:" + err.Error())
		return false
	}
	return true
}

// QueryUserFileMetas 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := global.DB.GetConn().Prepare("select file_sha1, file_name, file_size, upload_at, last_update from tbl_user_file where user_name = ? limit ?")
	if err != nil {
		log.Error("Failed to select, err:" + err.Error())
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username, limit)
	if err != nil {
		log.Error("Failed to select, err:" + err.Error())
		return nil, err
	}
	var userFiles []UserFile
	for rows.Next() {
		userFile := UserFile{}
		err := rows.Scan(&userFile.FileHash, &userFile.FileName, &userFile.FileSize, &userFile.UploadAt, &userFile.LastUpdated)
		if err != nil {
			log.Println(err)
			break
		}
		userFiles = append(userFiles, userFile)
	}
	return userFiles, nil
}

// 删除文件元信息
func DeleteUserFileMeta(filehash, filename string) bool {
	stmt, err := global.DB.GetConn().Prepare("delete from tbl_user_file where file_sha1 = ? and file_name = ? limit 1")
	if err != nil {
		log.Error("Failed to delete, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(filehash, filename)
	if err != nil {
		log.Error("Failed to delete, err:" + err.Error())
		return false
	}
	return true
}
