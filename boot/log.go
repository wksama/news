package boot

import (
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

func initLog() {
	var f io.Writer
	f = os.Stdout
	logFile := viper.GetString("log.file")
	if logFile != "" {
		if _, err := os.Stat(logFile); err != nil {
			f, err = os.Create(logFile)
			if err != nil {
				panic(err)
			}
		}
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	log.SetOutput(f)
}
