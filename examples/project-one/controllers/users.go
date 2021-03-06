// This is an example that uses the purely "functional" approach (No DI, no structs)
// It's fine for small projects, or if you already use some DI method, or if you want less technical debt

package controllers

import users "project-one/models/users"

// @HttpGet("/api/v1/users/{userId}/info")
// @Response(400, "InvalidUserId: UserId is invalid", models.ErrorResponse)
// @Required("userId", "username")
// @Parameter("userId", {default: 1})
func GetUserById(userId int64) users.GetUserById {
	return users.GetUserById{
		Id:       userId,
		Username: "Test123",
	}
}

// @HttpGet("/api/v1/users/username")
func GetUserByName(username string) users.GetUserById {
	return users.GetUserById{
		Id:       123,
		Username: username,
	}
}
