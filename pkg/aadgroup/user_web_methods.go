package aadgroup

import (
	"net/http"

	"github.com/dfds/manage-aadgroup-members/pkg/config"
	"github.com/dfds/manage-aadgroup-members/pkg/logging"
	"github.com/gin-gonic/gin"
)

func GetIndex(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/aadgroup/swagger/index.html")
}

// Health check verifies connectivity and login to Microsoft Graph API.
// GraphClient will kill the main process if it is not working.
func GetHealth(c *gin.Context) {
	log := logging.GetLogger()
	client, err := getGraphClient()
	_ = client
	statusCode := 200
	message := "OK"
	if err != nil {
		log.Error(err.Error())
		statusCode = 500
		message = "ERROR"
	}
	c.JSON(statusCode, gin.H{
		"message": message,
	})
}

// GetUser godoc
// @BasePath 		/aadgroup/api/v1
// @Summary 		Get user
// @Schemes
// @Description 	Return a single user based on User Principal Name
// @Tags 			azuread group user
// @Accept 			json
// @Produce 		json
// @Param 			upn	path	string	true	"User Principal Name"
// @Success 		200 {object} employee
// @Router 			/user/{upn} [get]
func GetUser(c *gin.Context) {
	log := logging.GetLogger()
	upn := c.Param("upn")
	employee, err := getUserFromAAD(upn)
	if err != nil {
		log.Error(err.Error())
	}
	c.IndentedJSON(http.StatusOK, employee)
}

// AddUsers godoc
// @BasePath 		/aadgroup/api/v1
// @Summary 		Add users to group
// @Schemes
// @Description 	Add a list of users to a group
// @Tags 			azuread group user
// @Accept 			multipart/form-data
// @Produce 		json
// @Param			groupMembers		formData	GroupMembers	true	"The group members with group ObjectId"
// @Success 		201 {boolean} true
// @Router 			/users [post]
func AddUsers(c *gin.Context) {
	log := logging.GetLogger()
	retVal := true
	var group GroupMembers
	err := c.Bind(&group)
	if err != nil {
		log.Error(err.Error())
		retVal = false
	}

	// If GroupObjectId is not supplied through the API, then read it from environment variable
	groupObjectId := group.GroupObjectId
	if groupObjectId == "" {
		groupObjectId = config.GroupId()
	}

	err = addUsersToGroup(groupObjectId, group.UserPrincipalNames)
	if err != nil {
		log.Error(err.Error())
		retVal = false
	}

	c.IndentedJSON(http.StatusCreated, retVal)
}

// RemoveUser godoc
// @BasePath 		/aadgroup/api/v1
// @Summary 		Remove user from group
// @Schemes
// @Description 	Remove a single user from a group
// @Tags 			azuread group user
// @Accept 			multipart/form-data
// @Produce 		json
// @Param			groupMember		formData	GroupMember	true	"The group member with group ObjectId"
// @Success 		204 {string} No Content
// @Router 			/user [delete]
func RemoveUser(c *gin.Context) {
	log := logging.GetLogger()
	statusCode := 204
	message := "No Content"
	var group GroupMember
	err := c.Bind(&group)
	if err != nil {
		log.Error(err.Error())
		statusCode = 500
		message = err.Error()
	}

	// If GroupObjectId is not supplied through the API, then read it from environment variable
	groupObjectId := group.GroupObjectId
	if groupObjectId == "" {
		groupObjectId = config.GroupId()
	}

	err = removeUserFromGroup(groupObjectId, group.UserPrincipalName)
	if err != nil {
		log.Error(err.Error())
		statusCode = 500
		message = err.Error()
	}

	c.JSON(statusCode, gin.H{
		"message": message,
	})
}
