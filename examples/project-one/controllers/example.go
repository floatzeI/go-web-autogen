package controllers

import usersmodels "project-one/models/users"

// @HttpGet("/api/v1/users/{userId}/info")
// @Response(400, "InvalidUserId: UserId is invalid", models.ErrorResponse)
// @Required("userId", "username")
// @Parameter("userId", {default: 1})
func GetUserById(userId int64) usersmodels.GetUserById {
	return usersmodels.GetUserById{
		Id:       userId,
		Username: "Test123",
	}
}

// @HttpGet("/api/v1/users/username")
func GetUserByName(username string) usersmodels.GetUserById {
	return usersmodels.GetUserById{
		Id:       123,
		Username: username,
	}
}
