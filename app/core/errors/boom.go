package errors

const (
	InternalError       = "internalError"
	UserNotFound        = "userNotFound"
	InvalidBindingModel = "invalidBindingModel"
	EntityCreationError = "entityCreationError"
	DeplicateUserFound  = "userAlreadyExists"
)

var errorMessage = map[string]string{
	"internalError":       "an internal error occured",
	"userNotFound":        "user could not be found",
	"userAlreadyExists":   "The user with this email is already exists.",
	"invalidBindingModel": "model could not be bound",
	"entityCreationError": "could not create entity",
}

func ErrorText(code string) string {
	return errorMessage[code]
}
