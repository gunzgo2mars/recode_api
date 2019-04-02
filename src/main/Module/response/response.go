package response

func ResponseMessage(title string, message string, err bool, status int) map[string]interface{} {

	return map[string]interface{}{
		"status":  status,
		"title":   title,
		"message": message,
		"error":   err,
	}
}
