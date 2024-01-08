package api

import (
	"fmt"
	"genericAPI/api/api_config"
	"genericAPI/api/environment"
	"genericAPI/api/os_specific_constants"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

func ConfigureGinLogger() {
	if !environment.IsTestEnvironment() {
		gin.DisableConsoleColor()
		f, _ := os.Open(api_config.Config.App.Logging.LogFilepath + time.Now().Format("2006-01-02.log"))
		gin.DefaultWriter = io.MultiWriter(f)
	} else {
		wd, _ := os.Getwd()
		f, _ := os.Open(fmt.Sprintf("%s%s%s", wd, os_specific_constants.PATH_SEPERATOR, time.Now().Format("2006-01-02.log")))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
}
