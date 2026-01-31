package response

import (
	"distributed_system/pkg/errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response represents a standard API response
// Matches v1's Kaos response pattern
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo represents error details in response
type ErrorInfo struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// Meta represents pagination metadata
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

// Success sends a success response with data
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage sends a success response with message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Resource created successfully",
		Data:    data,
	})
}

// NoContent sends a 204 no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Paginated sends a paginated response with meta information
func Paginated(c *gin.Context, data interface{}, page, pageSize int, total int64) {
	// Prevent division by zero
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}
	if page <= 0 {
		page = 1 // Default page
	}

	totalPages := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	appErr, ok := errors.As(err)
	if !ok {
		// Wrap unknown errors
		appErr = errors.Wrap(err, errors.ErrCodeInternal, "An unexpected error occurred")
	}

	// Get request ID from context (set by request_id middleware)
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	c.JSON(appErr.HTTPStatus, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:      appErr.Code,
			Message:   appErr.Message,
			Details:   appErr.Details,
			RequestID: requestID,
		},
	})
}

// ErrorWithStatus sends an error response with specific HTTP status
func ErrorWithStatus(c *gin.Context, status int, code, message string) {
	c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeInvalidInput,
			Message: message,
		},
	})
}

// Unauthorized sends a 401 unauthorized response
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeUnauthorized,
			Message: message,
		},
	})
}

// Forbidden sends a 403 forbidden response
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Access denied"
	}
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeForbidden,
			Message: message,
		},
	})
}

// NotFound sends a 404 not found response
func NotFound(c *gin.Context, resource string) {
	message := "Resource not found"
	if resource != "" {
		message = resource + " not found"
	}
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeNotFound,
			Message: message,
		},
	})
}

// Conflict sends a 409 conflict response
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeDuplicate,
			Message: message,
		},
	})
}

// InternalError sends a 500 internal server error response
func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeInternal,
			Message: "Internal server error",
		},
	})
}

// ValidationError sends validation errors
func ValidationError(c *gin.Context, validationErrors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeValidation,
			Message: "Validation failed",
		},
		Data: validationErrors,
	})
}

// Abort aborts the request with an error response
func Abort(c *gin.Context, err error) {
	appErr, ok := errors.As(err)
	if !ok {
		appErr = errors.Wrap(err, errors.ErrCodeInternal, "An unexpected error occurred")
	}

	// Get request ID from context
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	c.AbortWithStatusJSON(appErr.HTTPStatus, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:      appErr.Code,
			Message:   appErr.Message,
			Details:   appErr.Details,
			RequestID: requestID,
		},
	})
}

// AbortUnauthorized aborts with 401 status
func AbortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeUnauthorized,
			Message: "Authentication required",
		},
	})
}

// AbortForbidden aborts with 403 status
func AbortForbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeForbidden,
			Message: "Access denied",
		},
	})
}

// AbortWithMessage aborts with custom status, code, and message
func AbortWithMessage(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// AbortBadRequest aborts with 400 status
func AbortBadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeInvalidInput,
			Message: message,
		},
	})
}

// AbortTooManyRequests aborts with 429 status
func AbortTooManyRequests(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusTooManyRequests, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeRateLimit,
			Message: "Too many requests, please try again later",
		},
	})
}

// AbortServiceUnavailable aborts with 503 status
func AbortServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = "Service temporarily unavailable"
	}
	c.AbortWithStatusJSON(http.StatusServiceUnavailable, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeExternalService,
			Message: message,
		},
	})
}

// BindingError handles binding/validation errors from Gin
// It converts internal validation errors to user-friendly messages
// without exposing internal field names or validation rules
func BindingError(c *gin.Context, err error) {
	// Check if it's a validation error
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		// Convert validation errors to user-friendly messages
		messages := make([]string, 0, len(validationErrs))
		for _, e := range validationErrs {
			messages = append(messages, formatValidationError(e))
		}

		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    errors.ErrCodeValidation,
				Message: "Validation failed",
				Details: strings.Join(messages, "; "),
			},
		})
		return
	}

	// Check for JSON parsing errors
	if strings.Contains(err.Error(), "json:") || strings.Contains(err.Error(), "cannot unmarshal") {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    errors.ErrCodeInvalidInput,
				Message: "Invalid JSON format",
			},
		})
		return
	}

	// Generic invalid request error (don't expose internal error message)
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    errors.ErrCodeInvalidInput,
			Message: "Invalid request format",
		},
	})
}

// formatValidationError converts a validation error to a user-friendly message
func formatValidationError(e validator.FieldError) string {
	field := camelToReadable(e.Field())

	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + e.Param() + " characters"
	case "max":
		return field + " must be at most " + e.Param() + " characters"
	case "len":
		return field + " must be exactly " + e.Param() + " characters"
	case "gt":
		return field + " must be greater than " + e.Param()
	case "gte":
		return field + " must be at least " + e.Param()
	case "lt":
		return field + " must be less than " + e.Param()
	case "lte":
		return field + " must be at most " + e.Param()
	case "oneof":
		return field + " must be one of: " + e.Param()
	case "url":
		return field + " must be a valid URL"
	case "uuid":
		return field + " must be a valid UUID"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "numeric":
		return field + " must be numeric"
	case "eqfield":
		return field + " must match " + camelToReadable(e.Param())
	default:
		return field + " is invalid"
	}
}

// camelToReadable converts CamelCase to readable format
// e.g., "FirstName" -> "first name"
func camelToReadable(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
