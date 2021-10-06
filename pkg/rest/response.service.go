package rest

import (
	"idaman.id/storage/pkg/app"
)

func createSuccessResponse(param ResponseDto) ResponseEntity {

	response := ResponseEntity{
		Message: app.STATUS_OK,
	}

	if param.Data != nil {
		response.Data = param.Data
	}

	if param.Message != "" {
		response.Message = param.Message
	}

	if param.Translator != nil {
		translation := param.Translator(response.Message)
		response.Message = translation
	}

	return response
}

func createFailedResponse(param ResponseDto) ResponseEntity {

	response := ResponseEntity{
		Message: app.STATUS_ERROR,
	}

	if param.Error != nil {
		response.Error = param.Error
	}

	if param.Message != "" {
		response.Message = param.Message
	}

	if param.Translator != nil {
		translation := param.Translator(response.Message)
		response.Message = translation
	}

	return response
}
