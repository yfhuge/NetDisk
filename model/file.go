package model

import (
	"database/sql"
	"filestore-server/global"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// OnFileUpdateFinished 文件上传完成，保存文件元信息
func OnFileUpdateFinished(fileHash, fileName, fileAddr string, fileSize int64) bool {
	stmt, err := global.DB.GetConn().Prepare("insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		log.Error("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		log.Error("Failed to insert, err:" + err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Fatalf("file with hash:%s has been uploaded before\n", fileHash)
		}
		return true
	}
	return false
}

// GetFile 获取文件元信息
func GetFile(fileHash string) (*TableFile, error) {
	stmt, err := global.DB.GetConn().Prepare("select `file_sha1`, `file_addr`, `file_name`, `file_size` from tbl_file where file_sha1=? and status=1 limit 1")
	if err != nil {
		log.Error("Failed to select, err:" + err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := &TableFile{}
	err = stmt.QueryRow(fileHash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return tfile, nil
}

// UpdateFileLocation 更新文件的存储地址
func UpdateFileLocation(fileHash string, fileAddr string) error {
	stmt, err := global.DB.GetConn().Prepare("update tbl_file set `file_addr` = ? where `file_sha1`=? limit 1")
	if err != nil {
		log.Error("Failed to update, err:" + err.Error())
		return err
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileAddr, fileHash)
	if err != nil {
		log.Error("Failed to update, err:" + err.Error())
		return err
	}
	rf, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if rf <= 0 {
		log.Printf("更新文件location失败, filehash:%s\n", fileHash)
		err = fmt.Errorf("更新文件location失败，filehahs:%s", fileHash)
		return err
	}
	return nil
}

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// UpdateFileMetaDB 新增/更新文件元信息到MySQL
func UpdateFileMetaDB(fileMeta FileMeta) bool {
	return OnFileUpdateFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.Location, fileMeta.FileSize)
}

// GetFileMetaDB 获取MySQL中文件元信息
func GetFileMetaDB(fileHash string) (FileMeta, error) {
	tfile, err := GetFile(fileHash)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}
