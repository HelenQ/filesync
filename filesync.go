package filesync

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	config "gitlab.xinc818.com/qinjinyang/filesync/config/oss"
)

func SyncLocal2Oss(ossSync config.OssSync, localPath string, remotePath string) error {
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return errors.New("SyncLocal2Oss fail, " + err.Error())
	}
	if !fileInfo.IsDir() {
		return errors.New("SyncLocal2Oss fail, localPath is not directory")
	}
	bucket, err := ossSync.GetBucket()
	if err != nil {
		return errors.New("SyncLocal2Oss fail, " + err.Error())
	}
	return syncDir(bucket, localPath, remotePath)
}

func syncDir(bucket *oss.Bucket, localPath string, remotePath string) error {
	localKeyMap, err := getFileKeyMap(localPath)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(remotePath, "/") {
		remotePath = remotePath + "/"
	}
	// 判断远程目录是否存在
	isExist, err := bucket.IsObjectExist(remotePath)
	if err != nil {
		return err
	}
	if !isExist {
		bucket.PutObject(remotePath, bytes.NewReader([]byte("")))
	}
	remoteKeyMap := make(map[string]oss.ObjectProperties)
	prefix := oss.Prefix(remotePath)
	continuationToken := ""
	maxKey := oss.MaxKeys(100)
	for {
		lsRes, err := bucket.ListObjectsV2(prefix, maxKey, oss.ContinuationToken(continuationToken))
		if err != nil {
			return err
		}
		for _, object := range lsRes.Objects {
			if remotePath != object.Key {
				remoteKeyMap[object.Key] = object
			}
		}
		if lsRes.IsTruncated {
			continuationToken = lsRes.NextContinuationToken
		} else {
			break
		}
	}
	for k, p := range localKeyMap {
		fileInfo, err := os.Stat(p)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			remoteObject, ok := remoteKeyMap[remotePath+k]
			if !ok || fileInfo.ModTime().After(remoteObject.LastModified) {
				bucket.PutObjectFromFile(remotePath+k, p)
			}
		}
		delete(remoteKeyMap, remotePath+k)
	}
	for k := range remoteKeyMap {
		bucket.DeleteObject(k)
	}
	return nil
}

// 平铺文件 [相对路径:绝对路径]
func getFileKeyMap(dirPath string) (map[string]string, error) {
	fileKeyMap := make(map[string]string)
	entryList, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entryList {
		if entry.IsDir() {
			// 文件夹
			fileKeyMap[entry.Name()] = dirPath + "/" + entry.Name() + "/"
			// 子目录
			subDirPath := strings.Join([]string{dirPath, entry.Name()}, "/")
			subFileKeyMap, err := getFileKeyMap(subDirPath)
			if err != nil {
				return nil, err
			}
			for fileKey, filePath := range subFileKeyMap {
				fileKeyMap[strings.Join([]string{entry.Name(), fileKey}, "/")] = filePath
			}
		} else {
			fileKeyMap[entry.Name()] = strings.Join([]string{dirPath, entry.Name()}, "/")
		}
	}
	return fileKeyMap, nil
}
