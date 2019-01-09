package usermanagement

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jurekbarth/pup/worker/internal/project"
	"github.com/jurekbarth/pup/usermanager"
	"github.com/jurekbarth/pup/worker"
)

// CreateResources creates a user and  all necessary groups for the user
func createResources(svc *cognitoidentityprovider.CognitoIdentityProvider, userPoolID string, clientID string, userName string, userEmail string, userPassword string, forceCreation bool, groups []string) error {
	cip := usermanager.Usermanager{
		Svc:        svc,
		UserPoolID: userPoolID,
		ClientID:   clientID,
	}

	userExists, err := cip.CheckIfUserExists(userName)
	if err != nil {
		return err
	}
	if forceCreation && userExists {
		if err := cip.DeleteUser(userName); err != nil {
			return err
		}
	}
	if forceCreation || !userExists {
		if err := cip.CreateUser(userName, userEmail, userPassword); err != nil {
			return err
		}
		if err := cip.ActivateUser(userName, userPassword); err != nil {
			return err
		}
	}
	for _, groupName := range groups {
		groupExists, err := cip.CheckIfGroupExists(groupName)
		if err != nil {
			return err
		}
		if !groupExists {
			err := cip.CreateGroup(groupName)
			if err != nil {
				return err
			}
		}
		err = cip.AddUserToGroup(userName, groupName)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create creates all necessary resources for cognito
func Create(w *worker.Worker, users []project.User) error {
	session, err := worker.MakeSession(w, nil)
	if err != nil {
		return err
	}
	c := *w.Config
	svc := cognitoidentityprovider.New(session)
	userPoolID := c.AWSCognitoPoolID
	clientID := c.AWSCognitoClientBackendClientID
	forceCreation := true
	for _, userresource := range users {
		userName:=userresource.Username
		userEmail:= userName + "@" + c.EmailDomain
		userPassword := userresource.Password
		groups:= userresource.Groups
		err := createResources(svc, userPoolID, clientID, userName, userEmail, userPassword, forceCreation, groups)
		if err !=nil {
			return err
		}
	}
	return nil
}
