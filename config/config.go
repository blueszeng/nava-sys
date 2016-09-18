package config

import (
	"log"
	"encoding/json"
	"os"
)

//TODO: เมื่อรันจริงต้องเปลี่ยนเป็น Docker Network Bridge IP เช่น 172.17.0.3 เป็นต้น
type Config struct {
	DBHost string `json:"db_host"`
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
}
func LoadDSN(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Err Open file %v: Error is: %v", file, err)
	}
	decoder := json.NewDecoder(file)
	c := new(Config)
	err = decoder.Decode(&c)
	if err != nil {
		log.Println("error Decode Json:", err)
	}
	//log.Printf("Test Variable: %s", c.DBHost)
	dsn := c.DBUser + ":" + c.DBPass + "@" + c.DBHost + "/" + c.DBName + "?parseTime=true"
	//dsn := "root:mypass@tcp(nava.work:3306)/sys?parseTime=true"
	log.Println("DSN =",dsn)
	return dsn
}
