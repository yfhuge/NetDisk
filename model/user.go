package model

import (
	"database/sql"
	"filestore-server/global"
	log "github.com/sirupsen/logrus"
)

type User struct {
	UserName     string `json:"UserName"`
	Email        string `json:"Email"`
	Phone        string `json:"Phone"`
	SignupAt     string `json:"SignupAt"`
	LastActiveAt string `json:"LastActiveAt"`
	Status       int    `json:"Status"`
}

type UserInfo struct {
	Location string `json:"location"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

// UserSignup 通过用户名即密码完成user注册
func UserSignup(username string, passwd string) bool {
	stmt, err := global.DB.GetConn().Prepare("insert into tbl_user(`user_name`, `user_pwd`) values(?, ?)")
	if err != nil {
		log.Error("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		log.Error("Failed to insert, err:" + err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil && rf > 0 {
		return true
	}
	return false
}

// UserSignIn 判断登录密码是否一致
func UserSignIn(username string, encPwd string) bool {
	stmt, err := global.DB.GetConn().Prepare("select * from tbl_user where user_name = ? limit 1")
	if err != nil {
		log.Error("Failed to select, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		log.Error("Failed to select, err:" + err.Error())
		return false
	} else if rows == nil {
		log.Error("username not found:" + username)
		return false
	}
	pRows := ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encPwd {
		return true
	}
	return false
}

// UpdateToken 更新用户的Token
func UpdateToken(username string, token string) bool {
	stmt, err := global.DB.GetConn().Prepare("replace into tbl_user_token(`user_name`, `user_token`) values(?, ?)")
	if err != nil {
		log.Error("Failed to replace token, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		log.Error("Failed to replace token, err:" + err.Error())
		return false
	}
	return true
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := global.DB.GetConn().Prepare("select user_name, signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Error("Failed to select, err:" + err.Error())
		return user, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.UserName, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
