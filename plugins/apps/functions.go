package apps

import (
	"fmt"
	"os"
	"github.com/dokku/dokku/plugins/common"
)

func Apps_create(appName string){
	DOKKU_ROOT := os.Getenv("DOKKU_ROOT")
	err := common.VerifyAppName(appName)
	fmt.Println(err)
	out := os.MkdirAll(DOKKU_ROOT+"/"+appName,0755)
	fmt.Println(out)
	common.LogInfo1Quiet("Creating "+appName+"... done")
	common.PlugnTrigger("post-create",appName)
}

