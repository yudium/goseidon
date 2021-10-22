package rest

import (
	"github.com/gofiber/fiber/v2"
	"idaman.id/storage/pkg/app"
	"idaman.id/storage/pkg/deleting"
	"idaman.id/storage/pkg/retrieving"
	"idaman.id/storage/pkg/translation"
	"idaman.id/storage/pkg/uploading"
)

func createGetDetailHandler(dependency *Dependency) Handler {
	return func(ctx Context) Result {
		localizer := dependency.getLocalizer(ctx)
		translator := translation.NewGoI18nService(localizer).Translate

		fileEntity, err := retrieving.GetFile(ctx.Params("identifier"))
		isFileAvailable := err == nil
		if isFileAvailable {
			response := createSuccessResponse(ResponseParam{
				Translator: translator,
				Data:       fileEntity,
			})
			return ctx.JSON(response)
		}

		var statusCode int
		var response ResponseEntity

		switch err.(type) {
		case *app.NotFoundError:
			notFoundError := err.(*app.NotFoundError)
			statusCode = fiber.StatusNotFound
			response = createFailedResponse(ResponseParam{
				Message:    notFoundError.Error(),
				Translator: translator,
				TranslationData: map[string]interface{}{
					"context": notFoundError.Context,
				},
			})
		default:
			statusCode = fiber.StatusBadRequest
			response = createFailedResponse(ResponseParam{
				Message:    err.Error(),
				Translator: translator,
			})
		}

		return ctx.Status(statusCode).JSON(response)
	}
}

func createDeleteFileHandler(dependency *Dependency) Handler {
	return func(ctx Context) Result {
		localizer := dependency.getLocalizer(ctx)
		translator := translation.NewGoI18nService(localizer).Translate

		err := deleting.DeleteFile(ctx.Params("identifier"))
		isSuccessDelete := err == nil
		if isSuccessDelete {
			response := createSuccessResponse(ResponseParam{
				Translator: translator,
			})
			return ctx.JSON(response)
		}

		var statusCode int
		var response ResponseEntity

		switch err.(type) {
		case *app.NotFoundError:
			notFoundError := err.(*app.NotFoundError)
			statusCode = fiber.StatusNotFound
			response = createFailedResponse(ResponseParam{
				Message:    notFoundError.Error(),
				Translator: translator,
				TranslationData: map[string]interface{}{
					"context": notFoundError.Context,
				},
			})
		default:
			statusCode = fiber.StatusBadRequest
			response = createFailedResponse(ResponseParam{
				Message:    err.Error(),
				Translator: translator,
			})
		}

		return ctx.Status(statusCode).JSON(response)
	}
}

func createGetResourceHandler(dependency *Dependency) Handler {
	return func(ctx Context) Result {
		localizer := dependency.getLocalizer(ctx)
		translator := translation.NewGoI18nService(localizer).Translate

		result, err := retrieving.RetrieveFile(ctx.Params("identifier"))
		isFileAvailable := err == nil

		if isFileAvailable {
			ctx.Set("Content-Type", result.File.Mimetype)
			return ctx.Send(result.FileData)
		}

		var response ResponseEntity
		var statusCode int

		switch err.(type) {
		case *app.NotFoundError:
			notFoundError := err.(*app.NotFoundError)
			statusCode = fiber.StatusNotFound
			response = createFailedResponse(ResponseParam{
				Message:    notFoundError.Error(),
				Translator: translator,
				TranslationData: map[string]interface{}{
					"context": notFoundError.Context,
				},
			})
		default:
			statusCode = fiber.StatusBadRequest
			response = createFailedResponse(ResponseParam{
				Message:    err.Error(),
				Translator: translator,
			})
		}

		return ctx.Status(statusCode).JSON(response)
	}
}

func createUploadFileHandler(dependency *Dependency) Handler {
	return func(ctx Context) Result {
		locale := dependency.getLocale(ctx)
		localizer := dependency.getLocalizer(ctx)
		translator := translation.NewGoI18nService(localizer).Translate

		form, err := ctx.MultipartForm()

		isFormInvalid := err != nil
		if isFormInvalid {
			response := createFailedResponse(ResponseParam{
				Message:    err.Error(),
				Translator: translator,
			})
			return ctx.Status(fiber.StatusBadRequest).JSON(response)
		}

		uploadData := uploading.UploadFileParam{
			Files:    form.File["files"],
			Provider: ctx.FormValue("provider"),
			Locale:   locale,
		}
		result, err := uploading.UploadFile(uploadData)

		isUploadSuccess := err == nil
		if isUploadSuccess {
			response := createSuccessResponse(ResponseParam{
				Data:       result.Items,
				Translator: translator,
			})
			return ctx.JSON(response)
		}

		var response ResponseEntity
		var status int

		switch err.(type) {
		case *app.ValidationError:
			status = fiber.StatusUnprocessableEntity
			validationError := err.(*app.ValidationError)
			response = createFailedResponse(ResponseParam{
				Message:    validationError.Error(),
				Error:      validationError.Items,
				Translator: translator,
			})
		default:
			status = fiber.StatusBadRequest
			response = createFailedResponse(ResponseParam{
				Message:    err.Error(),
				Translator: translator,
			})
		}

		return ctx.Status(status).JSON(response)
	}
}
