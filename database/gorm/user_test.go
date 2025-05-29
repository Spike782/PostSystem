package database

import (
	"PostSystem/util"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func init() {
	util.InitSlog("../../log/post.log")
	CreateConnection("../../conf", "db", "yaml", "../../log")
}

func hash(pass string) string {
	hasher := md5.New()
	hasher.Write([]byte(pass))
	digest := hasher.Sum(nil)
	return hex.EncodeToString(digest)
}

func TestRegistUser(t *testing.T) {
	uid, err := RegistUser("xiaoming", hash("123456"))
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("注册成功！uid=%d\n", uid)
	}

}

func TestLogOffUser(t *testing.T) {
	err := LogOffUser(2)
	if err != nil {
		t.Fatal(err)
	}
}
