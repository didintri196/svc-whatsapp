package libraries

import (
	"fmt"

	_ "github.com/lib/pq"
	waLog "go.mau.fi/whatsmeow/util/log"

	"go.mau.fi/whatsmeow/store/sqlstore"
)

type WhatsmeowLibrary struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBName     string
}

func (wm WhatsmeowLibrary) Connect() (container *sqlstore.Container, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		wm.DBHost, wm.DBPort, wm.DBUser, wm.DBPassword, wm.DBName)
	fmt.Println(psqlInfo)
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err = sqlstore.New("postgres", psqlInfo, dbLog)
	if err != nil {
		return container, err
	}
	return container, nil
}
