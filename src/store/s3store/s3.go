package s3store

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	JSONUPLOADER_RAND_STRLEN = 10
)

type S3Store struct {
	bucket string
	region string
	svc    *s3.S3
}

// initializes the store by creating S3 handler
func Init(bucket string, region string) (*S3Store, error) {
	var s3st S3Store

	s3st.bucket = bucket
	s3st.region = region

	//create session object which contains aws configuration information  like region, credentials etc
	var sess *session.Session
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(s3st.region)}, //override default aws.config
		Profile: "test",
	}))

	s3st.svc = s3.New(sess)
	return &s3st, nil
}

func (s3st *S3Store) UploadToStore(prefix string, key string, data []byte) (*s3.PutObjectOutput, error) {
	// rand string
	randstr := make([]byte, JSONUPLOADER_RAND_STRLEN)
	// reads 10 cryptographically secure pseudorandom numbers from rand.Reader and writes them to randstr
	if _, err := rand.Read(randstr); err != nil {
		return nil, err
	}

	input := &s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket:               aws.String(s3st.bucket),
		Key:                  aws.String(fmt.Sprintf("%s/%s-%X", prefix, key, randstr)), //201509151400/20150915140413-6CBDA28FDBB15E7AECC9
		ServerSideEncryption: aws.String("AES256"), //256-bit Advanced Encryption Standard
	}

	result, err := s3st.svc.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return result, aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return result, aerr
		}
	} else {
		return result, err
	}

	return result, nil
}
