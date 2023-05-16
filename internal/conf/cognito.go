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
	conf := &aws.Config{Region: aws.String("ap-southeast-1")}
	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}

	conn := Cognito{
		CognitoClient: cognito.New(sess),
	}

	return &conn, nil
}
