package boot

import (
	"log"
)

func initLog() {
	//f, err := os.OpenFile(viper.GetString("log.file"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	panic(errors.Errorf("error opening log file: %v", err))
	//}
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	//log.SetOutput(f)
}
