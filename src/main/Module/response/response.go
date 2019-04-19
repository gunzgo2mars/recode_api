package response

func ResponseMessage(title string, message string, err bool, status int) map[string]interface{} {

	return map[string]interface{}{
		"status":  status,
		"title":   title,
		"message": message,
		"error":   err,
	}
}

func ResponseSingleProduct() map[string]interface{} {

	return nil

}

func ResponsePayloadProducts() map[string]interface{} {

	return nil

}

func ResponseSinglePlace() map[string]interface{} {

	return nil

}

func ResponsePayloadPlaces() map[string]interface{} {

	return nil

}
