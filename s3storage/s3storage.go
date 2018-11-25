package s3

import (
	"bytes"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

const ginS3Key = "_s3Tools"

// Config is used to configure the injector and tool
type Config struct {
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Region   string `json:"region"`
}

// File is a descriptor for a file. This will essetially simply file access
type File struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

// Tool is used to provide functionality to the dev environment
type Tool struct {
	config         *aws.Config
	storageSession *session.Session
	client         *s3.S3
}

func (t *Tool) setup(Endpoint, Region, Key, Secret string) {
	t.storageSession = session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(Key, Secret, ""),
		Endpoint:    &Endpoint,
		Region:      &Region,
	})
	t.client = s3.New(t.storageSession)
}

// ListFiles returns a list of files from a bucket and path
func (t *Tool) ListFiles(bucket, path string) ([]File, error) {
	out, err := t.client.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucket,
		Prefix: &path,
	})

	if err != nil {
		return nil, err
	}

	var ret []File
	for _, v := range out.Contents {
		ret = append(ret, File{
			Name:         *v.Key,
			Size:         *v.Size,
			LastModified: *v.LastModified,
		})
	}
	return ret, nil
}

// GetPresignedURL returns a url for a file to be viewed for a duration
func (t *Tool) GetPresignedURL(bucket, path string, duration time.Duration) (string, error) {
	req, _ := t.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &path,
	})
	url, err := req.Presign(duration)
	if err != nil {
		return "", err
	}
	return url, nil
}

// UploadFile will take a byte array, and upload it to your bucket
func (t *Tool) UploadFile(bucket, path string, file []byte) error {
	reader := bytes.NewReader(file)
	_, err := t.client.PutObject(&s3.PutObjectInput{
		Bucket:        &bucket,
		Key:           &path,
		Body:          reader,
		ContentLength: aws.Int64(reader.Size()),
	})
	return err
}

// GetStorage will get the instance of S3Tool, returns nil if it is unavailable
func GetStorage(c *gin.Context) *Tool {
	v, exists := c.Get(ginS3Key)
	if !exists {
		return nil
	}
	return v.(*Tool)
}

// InjectStorage will the S3Tool into gin for use
func InjectStorage(config *Config) gin.HandlerFunc {
	var tool Tool
	return func(c *gin.Context) {
		c.Set(ginS3Key, &tool)
	}
}
