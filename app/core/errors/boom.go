package errors

const (
	InternalError       = "internalError"
	UserNotFound        = "userNotFound"
	InvalidBindingModel = "invalidBindingModel"
	EntityCreationError = "entityCreationError"
	EntityDeletionError = "entityDeletionError"
	DeplicateUserFound  = "userAlreadyExists"
)

var errorMessage = map[string]string{
	"internalError":       "an internal error occured",
	"userNotFound":        "user could not be found",
	"userAlreadyExists":   "The user with this email is already exists.",
	"invalidBindingModel": "model could not be bound",
	"entityCreationError": "could not create entity",
	"entityDeletionError": "Could not delete record.",
}

func ErrorText(code string) string {
	return errorMessage[code]
}
