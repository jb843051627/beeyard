package model

import "errors"

// 哨兵错误：跨层用 errors.Is / errors.As 区分缺失与空。
var (
	ErrApiaryNotFound     = errors.New("apiary not found")
	ErrHiveNotFound       = errors.New("hive not found")
	ErrQueenNotFound      = errors.New("queen not found")
	ErrAlertNotFound      = errors.New("alert not found")
	ErrInspectionNotFound = errors.New("inspection not found")
	ErrMaintenanceNotFound = errors.New("maintenance task not found")
	ErrHarvestNotFound    = errors.New("harvest not found")
	ErrTransferNotFound   = errors.New("transfer not found")
	ErrReadingNotFound    = errors.New("reading not found")

	ErrInvalidStatus  = errors.New("invalid status transition")
	ErrInvalidInput   = errors.New("invalid input")
	ErrDuplicateCode  = errors.New("duplicate hive code")
)

// ValidationError 可被 errors.As 识别的业务校验错误。
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func NewValidationError(field, msg string) *ValidationError {
	return &ValidationError{Field: field, Message: msg}
}
