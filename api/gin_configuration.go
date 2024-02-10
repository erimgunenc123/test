package api

import (
	"fmt"
	"genericAPI/api/api_config"
	"genericAPI/api/environment"
	"genericAPI/api/os_specific_constants"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"time"
)

func ConfigureGin(app *gin.Engine) {
	if !environment.IsTestEnvironment() {
		log.Print("Running on prod environment. Disabled: Console logging, Console colors")
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
		f, _ := os.OpenFile(api_config.Config.App.Logging.LogFilepath+time.Now().Format("2006-01-02.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		gin.DefaultWriter = io.MultiWriter(f)
	} else {
		log.Print("Running on test environment. Enabled: Console logging, Console colors")
		app.SetTrustedProxies(nil)
		wd, _ := os.Getwd()
		f, _ := os.OpenFile(fmt.Sprintf("%s%s%s", wd, os_specific_constants.PATH_SEPERATOR, time.Now().Format("2006-01-02.log")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
}
