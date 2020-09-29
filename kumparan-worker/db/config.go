package db

import (
	"fmt"

	"kumparan/kumparan-worker/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	Dbcon     *gorm.DB
	Errdb     error
	dbuser    string
	dbpass    string
	dbname    string
	dbaddres  string
	dbport    string
	sslmode   string
	dbtimeout string
)

func init() {

	fmt.Println("DB POSTGRES")
	dbuser = utils.GetEnv("DB_POSTGRES_USER", "kumparanuser")
	dbpass = utils.GetEnv("DB_POSTGRES_PASS", "kumparanpassword")
	dbname = utils.GetEnv("DB_POSTGRES_NAME", "kumparan")
	dbaddres = utils.GetEnv("DB_POSTGRES_HOST", "localhost")
	dbport = utils.GetEnv("DB_POSTGRES_PORT", "5432")
	sslmode = utils.GetEnv("DB_POSTGRES_SSLMODE", "disable")
	dbtimeout = utils.GetEnv("DB_POSTGRES_TIMEOUT", "5")

	if DbOpen() != nil {
		fmt.Println("Can't Open db POSTGRES")
	}
	Dbcon = GetDbCon()
	Dbcon = Dbcon.LogMode(true)

	Dbcon.AutoMigrate(&News{})

	Dbcon.DB().SetConnMaxLifetime(5)
	Dbcon.DB().SetMaxIdleConns(25)
	Dbcon.DB().SetMaxOpenConns(10)

}

// DbOpen ..
func DbOpen() error {
	args := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%s", dbaddres, dbport, dbuser, dbpass, dbname, sslmode, dbtimeout)
	Dbcon, Errdb = gorm.Open("postgres", args)

	if Errdb != nil {
		fmt.Println("open db Err ", Errdb)
		return Errdb
	}

	if errping := Dbcon.DB().Ping(); errping != nil {
		return errping
	}
	return nil
}

// GetDbCon ..
func GetDbCon() *gorm.DB {
	//TODO looping try connection until timeout
	// using channel timeout
	if errping := Dbcon.DB().Ping(); errping != nil {
		fmt.Println("Db Not Connect test Ping :", errping)
		errping = nil
		if errping = DbOpen(); errping != nil {
			fmt.Println("try to connect again but error :", errping)
		}
	}
	Dbcon.LogMode(true)
	return Dbcon
}
