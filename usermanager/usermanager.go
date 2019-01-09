package usermanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// CheckIfUserExists checks if a user exists for a given username
func (cip Usermanager) CheckIfUserExists(username string) (bool, error) {
	adminGetUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(cip.UserPoolID),
		Username:   aws.String(username),
	}
	_, err := cip.Svc.AdminGetUser(adminGetUserInput)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == cognitoidentityprovider.ErrCodeUserNotFoundException {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

// CreateUser creates a user
func (cip Usermanager) CreateUser(username string, email string, password string) error {
	userAttributes := []*cognitoidentityprovider.AttributeType{
		&cognitoidentityprovider.AttributeType{
			Name:  aws.String("email"),
			Value: aws.String(email),
		},
		&cognitoidentityprovider.AttributeType{
			Name:  aws.String("email_verified"),
			Value: aws.String("true"),
		},
	}
	newUser := &cognitoidentityprovider.AdminCreateUserInput{
		Username:          aws.String(username),
		TemporaryPassword: aws.String(password),
		MessageAction:     aws.String("SUPPRESS"),
		UserPoolId:        aws.String(cip.UserPoolID),
		UserAttributes:    userAttributes,
	}
	_, err := cip.Svc.AdminCreateUser(newUser)
	return err
}

// ActivateUser *activates* a user in cognito
func (cip Usermanager) ActivateUser(username string, password string) error {
	adminInitiateAuthInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   aws.String(cognitoidentityprovider.AuthFlowTypeAdminNoSrpAuth),
		ClientId:   aws.String(cip.ClientID),
		UserPoolId: aws.String(cip.UserPoolID),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
	}
	adminInitiateAuthResponse, err := cip.Svc.AdminInitiateAuth(adminInitiateAuthInput)
	if err != nil {
		return err
	}

	adminRespondToAuthChallengeInput := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
		ClientId:      aws.String(cip.ClientID),
		UserPoolId:    aws.String(cip.UserPoolID),
		ChallengeResponses: map[string]*string{
			"USERNAME":     aws.String(username),
			"NEW_PASSWORD": aws.String(password),
		},
		Session: adminInitiateAuthResponse.Session,
	}
	_, err = cip.Svc.AdminRespondToAuthChallenge(adminRespondToAuthChallengeInput)
	return err
}

// DeleteUser deletes a user with a given username
func (cip Usermanager) DeleteUser(username string) error {
	adminDeleteUserInput := &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(cip.UserPoolID),
		Username:   aws.String(username),
	}
	_, err := cip.Svc.AdminDeleteUser(adminDeleteUserInput)
	return err
}

// CheckIfGroupExists for a given groupsname
func (cip Usermanager) CheckIfGroupExists(groupname string) (bool, error) {
	getGroupInput := &cognitoidentityprovider.GetGroupInput{
		GroupName:  aws.String(groupname),
		UserPoolId: aws.String(cip.UserPoolID),
	}
	_, err := cip.Svc.GetGroup(getGroupInput)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == cognitoidentityprovider.ErrCodeResourceNotFoundException {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

// CreateGroup creates a group in cognito
func (cip Usermanager) CreateGroup(groupname string) error {
	createGroupInput := &cognitoidentityprovider.CreateGroupInput{
		GroupName:  aws.String(groupname),
		UserPoolId: aws.String(cip.UserPoolID),
	}
	_, err := cip.Svc.CreateGroup(createGroupInput)
	return err
}

// AddUserToGroup adds a user to a group
func (cip Usermanager) AddUserToGroup(username string, groupname string) error {
	adminAddUserToGroupInput := &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String(groupname),
		UserPoolId: aws.String(cip.UserPoolID),
		Username:   aws.String(username),
	}
	_, err := cip.Svc.AdminAddUserToGroup(adminAddUserToGroupInput)
	return err
}

// Usermanager is a tool to create, delecte and activate cognito user pool users
type Usermanager struct {
	Svc        *cognitoidentityprovider.CognitoIdentityProvider
	UserPoolID string
	ClientID   string
}
