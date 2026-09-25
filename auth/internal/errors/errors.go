package errors

const (
	MsgUserNotFound       = "user not found"
	MsgInvalidCredentials = "incorrect login details"

	MsgRefreshToken = "invalid or expired refresh token"
	MsgAuthRequired = "authorization required"

	MsgTokenGeneration = "token generation error"

	MsgDatabaseOperation = "database operation error"

	MsgInvalidData  = "incorrect data format"
	MsgUserCreation = "error creating a user"

	MsgUserRegistered  = "user successfully registered"
	MsgLoginSuccess    = "user successfully logged in"
	MsgTokensRefreshed = "tokens successfully updated"
	MsgUserUpdated     = "user successfully updated"
	MsgUserDeleted     = "user successfully deleted"
)
