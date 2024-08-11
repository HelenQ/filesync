package filesync_test

import (
	"fmt"
	"testing"

	"github.com/HelenQ/filesync"
	config "github.com/HelenQ/filesync/config/oss"
)

func TestSyncLocal2Oss(t *testing.T) {
	ossSync := config.OssSync{
		Endpoint:        "https://oss-cn-hangzhou.aliyuncs.com",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "hy-cloudreve",
	}
	err := filesync.SyncLocal2Oss(ossSync, "/Users/qinjinyang/Downloads/es", "es")
	fmt.Println(err)
}
