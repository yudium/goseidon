package builtin_app

import (
	"github.com/gofiber/fiber/v2"
	app_error "idaman.id/storage/internal/error"
	response "idaman.id/storage/internal/response"
	"idaman.id/storage/internal/retrieving"
	"idaman.id/storage/internal/uploading"
)

func NewFileGetDetailHandler(rService retrieving.FileGetter) Handler {
	return func(ctx *Context) error {
		fileDetail, err := rService.GetFile(ctx.Params("identifier"))
		if err != nil {
			var statusCode int
			var resBody *response.ResponseEntity

			switch err.(type) {
			case *app_error.NotfoundError:
				notFoundError := err.(*app_error.NotfoundError)
				statusCode = fiber.StatusNotFound
				resBody = response.NewErrorResponse(&response.ResponseParam{
					Message: notFoundError.Error(),
				})
			default:
				statusCode = fiber.StatusBadRequest
				resBody = response.NewErrorResponse(&response.ResponseParam{
					Message: err.Error(),
				})
			}

			return ctx.Status(statusCode).JSON(resBody)
		}

		fileEntity := &FileDetailEntity{
			UniqueId:  fileDetail.UniqueId,
			Name:      fileDetail.Name,
			Extension: fileDetail.Extension,
			Size:      fileDetail.Size,
			Mimetype:  fileDetail.Mimetype,
			Url:       fileDetail.Url,
			Path:      fileDetail.Path,
			CreatedAt: fileDetail.CreatedAt,
			UpdatedAt: fileDetail.UpdatedAt,
		}
		resBody := response.NewSuccessResponse(&response.ResponseParam{
			Data: fileEntity,
		})
		return ctx.JSON(resBody)
	}
}

func NewGetResourceHandler(rService retrieving.FileRetriever) Handler {
	return func(ctx *Context) error {
		result, err := rService.RetrieveFile(ctx.Params("identifier"))

		if err != nil {
			var responseEntity *response.ResponseEntity
			var statusCode int

			switch err.(type) {
			case *app_error.NotfoundError:
				notFoundError := err.(*app_error.NotfoundError)
				statusCode = fiber.StatusNotFound
				responseEntity = response.NewErrorResponse(&response.ResponseParam{
					Message: notFoundError.Error(),
				})
			default:
				statusCode = fiber.StatusBadRequest
				responseEntity = response.NewErrorResponse(&response.ResponseParam{
					Message: err.Error(),
				})
			}

			return ctx.Status(statusCode).JSON(responseEntity)
		}

		ctx.Set("Content-Type", result.File.Mimetype)
		return ctx.Send(result.FileData)
	}
}

func NewUploadFileHandler(uService uploading.UploadService) Handler {
	return func(ctx *Context) error {

		form, err := ctx.MultipartForm()

		if err != nil {
			responseEntity := response.NewErrorResponse(&response.ResponseParam{
				Message: err.Error(),
			})
			return ctx.Status(fiber.StatusBadRequest).JSON(responseEntity)
		}

		uploadData := uploading.UploadFileParam{
			Files:    form.File["files"],
			Provider: ctx.FormValue("provider"),
		}
		result, err := uService.UploadFile(uploadData)

		isUploadSuccess := err == nil
		if isUploadSuccess {
			responseEntity := response.NewSuccessResponse(&response.ResponseParam{
				Data: result.Items,
			})
			return ctx.JSON(responseEntity)
		}

		var responseEntity *response.ResponseEntity
		var status int

		switch err.(type) {
		case *app_error.ValidationError:
			status = fiber.StatusUnprocessableEntity
			validationError := err.(*app_error.ValidationError)
			responseEntity = response.NewErrorResponse(&response.ResponseParam{
				Message: validationError.Error(),
				Error:   validationError.Items,
			})
		default:
			status = fiber.StatusBadRequest
			responseEntity = response.NewErrorResponse(&response.ResponseParam{
				Message: err.Error(),
			})
		}

		return ctx.Status(status).JSON(responseEntity)
	}
}
