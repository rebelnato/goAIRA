package isolatedfunctions

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rebelnato/goAIRA/endpoints"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Requests struct {
	RequestId     string    `gorm:"<-:create","primaryKey"`
	ConsumerId    string    `gorm:"<-:create"`
	RequestFlow   string    `gorm:"<-:create"`
	RequestMethod string    `gorm:"<-:create"`
	Request       []byte    `gorm:"type:bytea" ,"<-:create"`
	CreatedAt     time.Time `gorm:"<-:create","autoCreateTime"`
	UpdatedAt     time.Time `gorm:"<-:create","autoUpdateTime"`
}

type Responses struct {
	RequestId string    `gorm:"<-:create"`
	Response  []byte    `gorm:"type:bytea","<-:create"`
	CreatedAt time.Time `gorm:"<-:create","autoCreateTime"`
	UpdatedAt time.Time `gorm:"<-:create","autoUpdateTime"`
}

var Db *gorm.DB

func Initiatedbconnection() {

	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai", endpoints.DbHost, os.Getenv("db_user"), os.Getenv("db_pass"), os.Getenv("db_name"))
	var err error
	Db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Println("Unable to connect to db")
	}
	log.Println("connected to db successfully")

	// Automatically creates/updates tables
	Db.AutoMigrate(&Requests{})
	Db.AutoMigrate(&Responses{})
}

func DbPing() bool {
	var ping bool
	sqlDb, err := Db.DB()
	if err != nil {
		log.Println("Unable to get db context from GORM")
	}

	if err := sqlDb.Ping(); err == nil {
		ping = true
	}
	return ping
}

func CreateRequestEntry(requestId, consumerId, requestFlow, requestMethod string, request []byte) error {
	result := Db.Create(&Requests{
		RequestId:     requestId,
		ConsumerId:    consumerId,
		RequestFlow:   requestFlow,
		RequestMethod: requestMethod,
		Request:       request,
	})
	if result.Error != nil {
		log.Println("Failed to insert in db due to ", result.Error)
		return result.Error
	}
	if result.RowsAffected == 1 {
		log.Println("Data was inserted using CreateRequestEntry")
	} else {
		log.Println("More than 1 rows were affected , actual affected row count ", result.RowsAffected)
	}

	return result.Error
}

func CreateResponseEntry(requestId string, response []byte) error {
	result := Db.Create(&Responses{
		RequestId: requestId,
		Response:  response,
	})
	if result.Error != nil {
		log.Println("Failed to insert in db due to ", result.Error)
		return result.Error
	}
	if result.RowsAffected == 1 {
		log.Println("Data was inserted using CreateResponseEntry")
	} else {
		log.Println("More than 1 rows were affected , actual affected row count ", result.RowsAffected)
	}

	return result.Error
}
