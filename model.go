package common

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm"
    "github.com/sirupsen/logrus"
    "time"
)

type BaseModel struct {
    ID           int       `gorm:"primary_key" json:"id"`
    CreationTime time.Time `json:"creation_time"`
    ModifiedTime time.Time `json:"modified_time"`
}

type BaseModelNormal struct {
    ID int `gorm:"primary_key" json:"id"`
}

type BaseModelName struct {
    ID   int    `gorm:"primary_key" json:"id"`
    Name string `json:"name"`
}

type BaseModelCreate struct {
    ID           int       `gorm:"primary_key" json:"id"`
    CreationTime time.Time `json:"creation_time"`
}

type AuthModel struct {
    BaseModel
    UserId    int    `json:"user_id"`
    AppKey    string `json:"app_key"`
    AppSecret string `json:"app_secret"`
    Status    int    `json:"status";gorm:"DEFAULT:0"`
}

func (AuthModel) TableName() string {
    return "auth"
}

var DB *gorm.DB

func InitDB(DbUri string) {
    var err error
    DB, err = gorm.Open("mysql", DbUri)
    if err != nil {
        logrus.Error(err)
        panic(err)
    }
    DB.SingularTable(true)
    DB.DB().SetMaxIdleConns(10)
    DB.DB().SetMaxOpenConns(500)
    DB.DB().SetConnMaxLifetime(2 * time.Second)
}

func CloseDB() {
    defer DB.Close()
}

type AuthRedis struct {
    UserId int    `json:"user_id"`
    Secret string `json:"secret"`
    Key    string `json:"key"`
}

type AuthGlobal struct {
    UserId int `json:"user_id"`
}
