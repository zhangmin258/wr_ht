package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"fmt"
)

//图片上传到阿里云
func UploadAliyun(filename, filepath string) (error, string) {
	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		return err, "1"
	}
	bucket, err := client.Bucket(BucketName)
	if err != nil {
		return err, "2"
	}
	path := Upload_dir
	path += filename
	err = bucket.PutObjectFromFile(path, filepath)
	if err != nil {
		fmt.Println(err.Error())
		return err, "3"
	}
	path = Imghost + path
	return err, path
}
