package restful

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcError(c *fiber.Ctx, st *status.Status) error {
	httpCode := 500
	message := "sorry, internal server error try again later"

	switch st.Code() {
	case codes.InvalidArgument:
		httpCode = 400
		message = st.Message()

	case codes.FailedPrecondition:
		httpCode = 400
		message = st.Message()

	case codes.Unauthenticated:
		httpCode = 401
		message = st.Message()

	case codes.PermissionDenied:
		httpCode = 403
		message = st.Message()

	case codes.NotFound:
		httpCode = 404
		message = st.Message()

	case codes.AlreadyExists:
		httpCode = 409
		message = st.Message()
	}


	return c.Status(httpCode).JSON(fiber.Map{"errors": message})
}
