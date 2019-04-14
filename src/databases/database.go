package databases

import (
	"database/sql"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

//DataManager will provide client to manage Redis and MySQL databases
type DataManager struct {
	RedisClient *redis.Client
	MySQLClient *sql.DB
}

//SetClients will intialize the clients for redis and mysql
func (ins *DataManager) SetClients() {
	datasourcestring := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PW") + "@/" + os.Getenv("MYSQL_DB")
	db, err := sql.Open("mysql", datasourcestring)
	checkErr(err)
	checkErr(db.Ping())
	ins.MySQLClient = db

	ins.RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PW"), // no password set
		DB:       0,                     // use default DB
	})

	ins.createTableIfNotExists()
}

//createTableIfNotExists will intialize the mysql DB and create table if not exist
func (ins *DataManager) createTableIfNotExists() {
	stmtIns, err := ins.MySQLClient.Prepare(`
		CREATE TABLE IF NOT EXISTS sites(
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			url VARCHAR(500) NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	checkErr(err)
	defer stmtIns.Close()
	_, err = stmtIns.Exec()
	checkErr(err)
}

//RegisterURLToMysql save a map of string (URL) to MYSQL
func (ins DataManager) RegisterURLToMysql(external map[string]int) {
	stmtIns, err := ins.MySQLClient.Prepare("INSERT IGNORE INTO sites (url) VALUES( ? )")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()
	for key := range external {
		_, err = stmtIns.Exec(key)
		checkErr(err)
	}
}

//RegisterURLToRedis save a slice of string (URL) to Redis
func (ins DataManager) RegisterURLToRedis(external []string) {
	for _, element := range external {
		err := ins.RedisClient.SetNX(element, "0", 0).Err()
		checkErr(err)
	}
}

//RemoveFromMysql will remove a given URL from mysql
func (ins DataManager) RemoveFromMysql(url string) {
	stmtIns, err := ins.MySQLClient.Prepare("DELETE FROM sites WHERE url = ?")
	checkErr(err)
	defer stmtIns.Close()
	_, err = stmtIns.Exec(url)
	checkErr(err)
}

//SetURLOnRedis will set a given URL as already worked on Redis
func (ins DataManager) SetURLOnRedis(current string) {
	checkErr(ins.RedisClient.Set(current, "1", 0).Err())
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
