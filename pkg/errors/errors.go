package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// AppError represents application-specific errors with context
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"`
	RequestID  string `json:"request_id,omitempty"`
	Stack      string `json:"-"` // For logging, not exposed to client
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// Wrap wraps an existing error with AppError
func Wrap(err error, code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

// WithStatus sets the HTTP status code
func (e *AppError) WithStatus(status int) *AppError {
	e.HTTPStatus = status
	return e
}

// WithDetails adds additional details
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithRequestID adds request ID for tracing
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// WithStack captures the current stack trace (for logging)
func (e *AppError) WithStack() *AppError {
	e.Stack = captureStack(2) // Skip WithStack and caller
	return e
}

// captureStack captures the call stack
func captureStack(skip int) string {
	var pcs [32]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var sb strings.Builder
	for {
		frame, more := frames.Next()
		// Skip runtime internals
		if strings.Contains(frame.Function, "runtime.") {
			if !more {
				break
			}
			continue
		}
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}

// Clone creates a copy of the error (prevents mutation of shared errors)
func (e *AppError) Clone() *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		Details:    e.Details,
		HTTPStatus: e.HTTPStatus,
		Err:        e.Err,
		RequestID:  e.RequestID,
		Stack:      e.Stack,
	}
}

// LogFields returns a map of fields for structured logging
func (e *AppError) LogFields() map[string]interface{} {
	fields := map[string]interface{}{
		"error_code":    e.Code,
		"error_message": e.Message,
		"http_status":   e.HTTPStatus,
	}
	if e.Details != "" {
		fields["details"] = e.Details
	}
	if e.RequestID != "" {
		fields["request_id"] = e.RequestID
	}
	if e.Err != nil {
		fields["cause"] = e.Err.Error()
	}
	if e.Stack != "" {
		fields["stack"] = e.Stack
	}
	return fields
}

// Error codes - matching v1 patterns
const (
	// Authentication & Authorization
	ErrCodeUnauthorized     = "ERR_UNAUTHORIZED"
	ErrCodeForbidden        = "ERR_FORBIDDEN"
	ErrCodeInvalidToken     = "ERR_INVALID_TOKEN"
	ErrCodeTokenExpired     = "ERR_TOKEN_EXPIRED"
	ErrCodeInvalidCredential = "ERR_INVALID_CREDENTIAL"

	// Validation
	ErrCodeValidation      = "ERR_VALIDATION"
	ErrCodeInvalidInput    = "ERR_INVALID_INPUT"
	ErrCodeMissingRequired = "ERR_MISSING_REQUIRED"

	// Database
	ErrCodeNotFound      = "ERR_NOT_FOUND"
	ErrCodeDuplicate     = "ERR_DUPLICATE"
	ErrCodeDBError       = "ERR_DATABASE"
	ErrCodeTransaction   = "ERR_TRANSACTION"

	// Business Logic
	ErrCodeInsufficientBalance = "ERR_INSUFFICIENT_BALANCE"
	ErrCodeInvalidStatus       = "ERR_INVALID_STATUS"
	ErrCodeExpired             = "ERR_EXPIRED"
	ErrCodeAlreadyUsed         = "ERR_ALREADY_USED"
	ErrCodeLimitExceeded       = "ERR_LIMIT_EXCEEDED"
	ErrCodeNotAvailable        = "ERR_NOT_AVAILABLE"

	// Payment
	ErrCodePaymentFailed    = "ERR_PAYMENT_FAILED"
	ErrCodePaymentPending   = "ERR_PAYMENT_PENDING"
	ErrCodeRefundFailed     = "ERR_REFUND_FAILED"
	ErrCodeInvalidPayment   = "ERR_INVALID_PAYMENT"

	// External Services
	ErrCodeExternalService = "ERR_EXTERNAL_SERVICE"
	ErrCodeTimeout         = "ERR_TIMEOUT"
	ErrCodeRateLimit       = "ERR_RATE_LIMIT"

	// System
	ErrCodeInternal = "ERR_INTERNAL"
	ErrCodeConfig   = "ERR_CONFIG"

	
)

// Pre-defined errors for common scenarios
var (
	// Auth errors
	ErrUnauthorized = &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    "Authentication required",
		HTTPStatus: http.StatusUnauthorized,
	}
	ErrForbidden = &AppError{
		Code:       ErrCodeForbidden,
		Message:    "Access denied",
		HTTPStatus: http.StatusForbidden,
	}
	ErrInvalidToken = &AppError{
		Code:       ErrCodeInvalidToken,
		Message:    "Invalid or malformed token",
		HTTPStatus: http.StatusUnauthorized,
	}
	ErrTokenExpired = &AppError{
		Code:       ErrCodeTokenExpired,
		Message:    "Token has expired",
		HTTPStatus: http.StatusUnauthorized,
	}
	ErrInvalidCredentials = &AppError{
		Code:       ErrCodeInvalidCredential,
		Message:    "Invalid email or password",
		HTTPStatus: http.StatusUnauthorized,
	}

	// Validation errors
	ErrValidation = &AppError{
		Code:       ErrCodeValidation,
		Message:    "Validation failed",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrInvalidInput = &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    "Invalid input data",
		HTTPStatus: http.StatusBadRequest,
	}

	// Database errors
	ErrNotFound = &AppError{
		Code:       ErrCodeNotFound,
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
	}
	ErrDuplicate = &AppError{
		Code:       ErrCodeDuplicate,
		Message:    "Resource already exists",
		HTTPStatus: http.StatusConflict,
	}

	// Business errors
	ErrInsufficientBalance = &AppError{
		Code:       ErrCodeInsufficientBalance,
		Message:    "Insufficient balance",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrInvalidStatus = &AppError{
		Code:       ErrCodeInvalidStatus,
		Message:    "Invalid status transition",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrExpired = &AppError{
		Code:       ErrCodeExpired,
		Message:    "Resource has expired",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrAlreadyUsed = &AppError{
		Code:       ErrCodeAlreadyUsed,
		Message:    "Resource already used",
		HTTPStatus: http.StatusBadRequest,
	}

	// Payment errors
	ErrPaymentFailed = &AppError{
		Code:       ErrCodePaymentFailed,
		Message:    "Payment processing failed",
		HTTPStatus: http.StatusBadRequest,
	}

	// System errors
	ErrInternal = &AppError{
		Code:       ErrCodeInternal,
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// NotFound creates a not found error with custom message
func NotFound(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		HTTPStatus: http.StatusNotFound,
	}
}

// IsNotFound checks if an error is a NotFound error.
// This allows usecase layer to check for not found without knowing about gorm.
// Usage: if errors.IsNotFound(err) { return nil, err }
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	appErr, ok := As(err)
	return ok && appErr.Code == ErrCodeNotFound
}

// InvalidInput creates a validation error for invalid input
func InvalidInput(message string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// Duplicate creates a duplicate/conflict error
func Duplicate(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeDuplicate,
		Message:    fmt.Sprintf("%s already exists", resource),
		HTTPStatus: http.StatusConflict,
	}
}

// Validation creates a validation error with details
func Validation(message string) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// Database creates a database error
func Database(err error) *AppError {
	return &AppError{
		Code:       ErrCodeDBError,
		Message:    "Database operation failed",
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

// External creates an external service error
func External(service string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeExternalService,
		Message:    fmt.Sprintf("External service error: %s", service),
		HTTPStatus: http.StatusBadGateway,
		Err:        err,
	}
}

// Is checks if the error is of a specific type
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As attempts to convert error to AppError
func As(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// GetHTTPStatus extracts HTTP status from error
func GetHTTPStatus(err error) int {
	if appErr, ok := As(err); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// ===== Edge Case Handling =====

// IsRetryable checks if an error is transient and can be retried
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	appErr, ok := As(err)
	if !ok {
		// Check for common retryable patterns in error message
		errStr := err.Error()
		retryablePatterns := []string{
			"connection refused",
			"connection reset",
			"timeout",
			"temporary failure",
			"service unavailable",
			"too many connections",
			"network is unreachable",
			"no such host",
			"i/o timeout",
		}
		for _, pattern := range retryablePatterns {
			if containsIgnoreCase(errStr, pattern) {
				return true
			}
		}
		return false
	}

	// Retryable error codes
	retryableCodes := map[string]bool{
		ErrCodeExternalService: true,
		ErrCodeTimeout:         true,
		ErrCodeRateLimit:       true,
		ErrCodeDBError:         true,
	}

	return retryableCodes[appErr.Code]
}

// IsClientError checks if error is due to client (4xx)
func IsClientError(err error) bool {
	status := GetHTTPStatus(err)
	return status >= 400 && status < 500
}

// IsServerError checks if error is due to server (5xx)
func IsServerError(err error) bool {
	status := GetHTTPStatus(err)
	return status >= 500 && status < 600
}

// Timeout creates a timeout error
func Timeout(operation string) *AppError {
	return &AppError{
		Code:       ErrCodeTimeout,
		Message:    fmt.Sprintf("Operation timed out: %s", operation),
		HTTPStatus: http.StatusGatewayTimeout,
	}
}

// RateLimit creates a rate limit error
func RateLimit(retryAfter int) *AppError {
	return &AppError{
		Code:       ErrCodeRateLimit,
		Message:    "Too many requests",
		Details:    fmt.Sprintf("Retry after %d seconds", retryAfter),
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// AccountLocked creates an account locked error
func AccountLocked(until string) *AppError {
	return &AppError{
		Code:       "ERR_ACCOUNT_LOCKED",
		Message:    "Account is temporarily locked",
		Details:    fmt.Sprintf("Try again after %s", until),
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// ServiceUnavailable creates a service unavailable error
func ServiceUnavailable(service string) *AppError {
	return &AppError{
		Code:       ErrCodeExternalService,
		Message:    fmt.Sprintf("Service temporarily unavailable: %s", service),
		HTTPStatus: http.StatusServiceUnavailable,
	}
}

// WithContext adds context information to an error
func (e *AppError) WithContext(key, value string) *AppError {
	if e.Details == "" {
		e.Details = fmt.Sprintf("%s=%s", key, value)
	} else {
		e.Details = fmt.Sprintf("%s; %s=%s", e.Details, key, value)
	}
	return e
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
		 len(s) > 0 && len(substr) > 0 &&
		 (s[0]|32) >= 'a' && (s[0]|32) <= 'z' &&
		 containsIgnoreCaseSlow(s, substr))
}

func containsIgnoreCaseSlow(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1, c2 := s[i+j], substr[j]
			if c1 != c2 && (c1|32) != (c2|32) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
