package aadgroup

import "github.com/dfds/manage-aadgroup-members/pkg/logging"

type employee struct {
	ObjectId          string `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	UserPrincipalName string `json:"upn"`
}

// getUserFromAAD returns the ObjectId, DisplayName and Email address of an Azure AD User
func getUserFromAAD(user string) (*employee, error) {
	log := logging.GetLogger()
	client, err := getGraphClient()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	usr, err := client.UsersById(user).Get(nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	employee := employee{ObjectId: *usr.GetId(), Name: *usr.GetDisplayName(), Email: *usr.GetMail(), UserPrincipalName: *usr.GetUserPrincipalName()}
	return &employee, nil
}
