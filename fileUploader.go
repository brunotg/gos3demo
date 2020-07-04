package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	us_west2 string = "us-west-2"
	us_east1 string = "us-east-1"
)

//
func getSession(awsRegion string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	fmt.Println("AWS Session Created")
	return sess, err
}

func listBuckets(awsRegion string) {

	sess, _ := getSession(awsRegion)
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)

	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets: ")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created %s \n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func listBucketItems(awsRegion string, bucket string) {

	sess, _ := getSession(awsRegion)
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in buckets %q %v", bucket, err)
	}
	for _, item := range resp.Contents {
		fmt.Println("Name:			", *item.Key)
		fmt.Println("Last modified:	", *item.LastModified)
		fmt.Println("Size:			", *item.Size)
		fmt.Println("Storage class:	", *item.StorageClass)
		fmt.Println("")
	}
}

func createBucket(awsRegion string, bucket string) {

	sess, _ := getSession(awsRegion)
	svc := s3.New(sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucket, err)
	}

	fmt.Printf("Waiting for bucket %q to be created ...\n", bucket)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	fmt.Printf("Bucket %s created", bucket)

}

func uploadFile(awsRegion string, bucketName string, fileName string) {

	file, err := os.Open(fileName)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)

	}
	defer file.Close()

	sess, _ := getSession(awsRegion)

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		exitErrorf("Unable to upload %q to %q, %v", fileName, bucketName, err)
	}

	fmt.Printf("Succesfully uploaded %q to %q\n", fileName, bucketName)

}

func main() {
	bucketName := "bruno-golang-5"
	createBucket(us_west2, bucketName)
	listBuckets(us_west2)
	listBucketItems(us_west2, bucketName)
	uploadFile(us_west2, bucketName, "baby.JPG")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
