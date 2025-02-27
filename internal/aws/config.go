package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsV2Config "github.com/aws/aws-sdk-go-v2/config"
	creds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"gopkg.in/ini.v1"
)

type Profiles struct {
	Data  []string
	Error string
}
type credentialProvider struct {
	awsV2.Credentials
}

func (c credentialProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey, SessionToken: os.Getenv("AWS_SESSION_TOKEN")}, nil
}

func (c credentialProvider) IsExpired() bool {
	return c.Expired()
}

func GetCfg(profile, region string) (awsV2.Config, error) {
	cfg, err := awsV2Config.LoadDefaultConfig(
		context.TODO(),
		awsV2Config.WithSharedConfigProfile(profile),
		awsV2Config.WithRegion(region),
	)
	if err != nil {
		fmt.Printf("failed to load config")
		return awsV2.Config{}, err
	}
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Printf("failed to read credentials")
		return awsV2.Config{}, err
	}

	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		fmt.Println("Credentials have expired")
		return awsV2.Config{}, errors.New("AWS Credentials expired")
	}
	return cfg, err
}

func GetCfgUsingEnvVariables(profile, region string) (awsV2.Config, error) {
	akid := aws.String(os.Getenv(AWS_ACCESS_KEY_ID))
	secKey := aws.String(os.Getenv(AWS_SECRET_ACCESS_KEY))
	cfg, err := awsV2Config.LoadDefaultConfig(
		context.TODO(),
		awsV2Config.WithSharedConfigProfile(profile),
		awsV2Config.WithRegion(region),
		config.WithCredentialsProvider(
			creds.NewStaticCredentialsProvider(*akid, *secKey, ""),
		),
	)
	if err != nil {
		fmt.Printf("failed to load config")
		return awsV2.Config{}, err
	}
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Printf("failed to read credentials")
		return awsV2.Config{}, err
	}

	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		fmt.Println("Credentials have expired")
		return awsV2.Config{}, errors.New("AWS Credentials expired")
	}
	return cfg, err
}

func GetProfiles() (profiles []string, err error) {
	fpCred := defaults.SharedCredentialsFilename()
	_, errCred := os.Stat(fpCred)
	fpConf := defaults.SharedConfigFilename()
	_, errConf := os.Stat(fpConf)
	if os.IsNotExist(errCred) && os.IsNotExist(errConf) {
		return nil, errConf
	}
	var ret []string
	defaultReturn := &Profiles{Data: nil, Error: ""}
	fp := defaults.SharedCredentialsFilename()
	_, err = os.Stat(fp)
	if os.IsNotExist(err) {
		fp = defaults.SharedConfigFilename()
	}
	f, err := ini.Load(fp) // Load ini file
	if err != nil {
		defaultReturn.Error = err.Error()
	} else {
		arr := []string{}
		for _, v := range f.Sections() {
			if len(v.Keys()) != 0 {
				arr = append(arr, v.Name())
			}
		}
		defaultReturn.Data = arr
	}
	for i := 0; i < len(defaultReturn.Data); i++ {
		spltiArr := strings.Split(defaultReturn.Data[i], " ")
		if len(spltiArr) == 1 {
			ret = append(ret, spltiArr[len(spltiArr)-1])
		} else if len(spltiArr) > 1 && spltiArr[0] == "profile" {
			ret = append(ret, spltiArr[len(spltiArr)-1])
		}
	}
	return ret, nil
}

func GetLocalstackCfg(region string) (awsV2.Config, error) {
	customResolver := awsV2.EndpointResolverFunc(func(service, region string) (awsV2.Endpoint, error) {
		return awsV2.Endpoint{
			URL:           "http://localhost:4566",
			SigningRegion: region,
		}, nil
	})

	awsLSCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}
	return awsLSCfg, nil
}
