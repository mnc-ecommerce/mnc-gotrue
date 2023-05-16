package conf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type Cognito struct {
	CognitoClient *cognito.CognitoIdentityProvider
}

func ConnectCognito(config *GlobalConfiguration) (*Cognito, error) {
	conf := &aws.Config{Region: aws.String("us-southeast-1")}
	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}

	conn := Cognito{
		CognitoClient:   cognito.New(sess),
		UserPoolID:      "ap-southeast-1_nwv3wd7Hy",
		AppClientID:     "3074gf3454m87ptt8vi28v1qdl",
		AppClientSecret: "17qcsitu42gnqodf8kpe47ft37m6plqqcjahnkc5hv5qr0hs5va",
	}

	return &conn, nil
}
