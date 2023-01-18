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
	statusCode := http.StatusOK
	message := http.StatusText(http.StatusOK)
	client, err := getGraphClient()
	_ = client
	if err != nil {
		log.Error(err.Error())
		statusCode = http.StatusInternalServerError
		message = err.Error()
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
	statusCode := http.StatusOK
	upn := c.Param("upn")
	employee, err := getUserFromAAD(upn)
	if err != nil {
		log.Error(err.Error())
		statusCode = http.StatusInternalServerError
	}
	c.IndentedJSON(statusCode, employee)
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
	statusCode := http.StatusCreated
	message := http.StatusText(http.StatusCreated)
	var group GroupMembers
	err := c.Bind(&group)
	if err != nil {
		log.Error(err.Error())
		statusCode = http.StatusInternalServerError
		message = err.Error()
	}

	// If GroupObjectId is not supplied through the API, then read it from environment variable
	groupObjectId := group.GroupObjectId
	if groupObjectId == "" {
		groupObjectId = config.GroupId()
	}

	err = addUsersToGroup(groupObjectId, group.UserPrincipalNames)
	if err != nil {
		log.Error(err.Error())
		statusCode = http.StatusInternalServerError
		message = err.Error()
	}

	c.JSON(statusCode, gin.H{
		"message": message,
	})
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
	statusCode := http.StatusNoContent
	message := http.StatusText(http.StatusNoContent)
	var group GroupMember
	err := c.Bind(&group)
	if err != nil {
		log.Error(err.Error())
		statusCode = http.StatusInternalServerError
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
		statusCode = http.StatusInternalServerError
		message = err.Error()
	}

	c.JSON(statusCode, gin.H{
		"message": message,
	})
}

// GetUsers godoc
// @BasePath 		/aadgroup/api/v1
// @Summary 		Get list of users from group
// @Schemes
// @Description 	Return all users from a group
// @Tags 			azuread group user
// @Accept 			json
// @Produce 		json
// @Param 			groupObjectId	path	string	false	"Group ObjectId"
// @Success 		200 {object} employee
// @Router 			/users/{groupObjectId} [get]
func GetUsers(c *gin.Context) {
	log := logging.GetLogger()
	statusCode := http.StatusOK

	groupObjectId := c.Param("groupObjectId")
	group := groupObjectId
	if group == "" {
		group = config.GroupId()
	}

	users, err := getUsersFromGroup(group)
	if err != nil {
		statusCode = http.StatusInternalServerError
		log.Error(err)
	}
	c.JSON(statusCode, gin.H{
		"message": users.GetValue(),
	})
}
