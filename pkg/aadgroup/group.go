package aadgroup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dfds/manage-aadgroup-members/pkg/logging"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type GroupMember struct {
	GroupObjectId      string `json:"groupObjectId" form:"groupObjectId"`
	UserPrincipalNames string `json:"userPrincipalName" form:"userPrincipalName" binding:"required"`
}

type GroupMembers struct {
	GroupObjectId      string   `json:"groupObjectId" form:"groupObjectId"`
	UserPrincipalNames []string `json:"userPrincipalNames" form:"userPrincipalNames" binding:"required"`
}

type GraphError struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Date            string `json:"date"`
			RequestID       string `json:"request-id"`
			ClientRequestID string `json:"client-request-id"`
		} `json:"innerError"`
	} `json:"error"`
}

// AddUsersToGroup takes a list of user principal names (UPNs), then places each UPN into an Azure AD Group
func addUsersToGroup(groupObjectId string, userPrincipalNames []string) error {
	log := logging.GetLogger()
	client, err := getGraphClient()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, upn := range userPrincipalNames {
		user, err := getUserFromAAD(upn)
		if err != nil {
			log.Error(err)
			return err
		} else {
			err := addMemberToGroup(client, groupObjectId, user)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

// getGroupName takes a GraphServiceClient and a group object id, and return the name of an Azure AD Group
func getGroupName(client *msgraphsdk.GraphServiceClient, groupObjectId string) (string, error) {
	log := logging.GetLogger()
	grp, err := client.GroupsById(groupObjectId).Get(nil)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return *grp.GetDisplayName(), nil
}

// addMemberToGroup takes a group's object id and a user's (or service principal's) object id
func addMemberToGroup(client *msgraphsdk.GraphServiceClient, groupObjectId string, user *employee) error {
	log := logging.GetLogger()
	groupUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s/members/$ref", groupObjectId)
	memberUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", user.ObjectId)

	values := map[string]string{"@odata.id": memberUrl}

	payload, err := json.Marshal(values)
	if err != nil {
		log.Error(err)
		return err
	}

	bearerToken, err := GetBearerToken()
	if err != nil {
		log.Error(err)
		return err
	}

	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", groupUrl, bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", bearerToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	bytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
		return err
	}

	/*
		The response body will only contain the expected error when something goes wrong.
		Hence we don't care about handling any json.Unmarshal errors
	*/
	var graphError GraphError
	unmarshalErr := json.Unmarshal(bytes, &graphError)
	_ = unmarshalErr

	if graphError != (GraphError{}) {
		log.Debug(string(bytes)) // TODO: Remove debug
		if graphError.Error.Code == "Request_BadRequest" {
			log.Error(err.Error())
			return nil
		}
		return errors.New(graphError.Error.Code)
	}

	if resp.StatusCode == 401 {
		return errors.New(resp.Status)
	}

	groupName, err := getGroupName(client, groupObjectId)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("Added %s to %s group", user.Name, groupName)
	return nil
}
