package cerror

import "errors"

var (
	// ErrInternalServer error internal server
	ErrInternalServer = errors.New("internal server error")
	// ErrLenderNotFound error when data lender not found
	ErrLenderNotFound = errors.New("lender data not found")
	// ErrRDLLenderNotFound error when main rdl not found
	ErrRDLLenderNotFound = errors.New("user do not have rdl yet, please ask them to initiate one")
	// ErrSubscriptionNotFound error when data subscription not found
	ErrSubscriptionNotFound = errors.New("subscription data not found")
	// ErrIdleFundNotFound error when data subscription not found
	ErrIdleFundNotFound = errors.New("idle fund data not found")
	// ErrBalanceInsuffiecient error when balance user insufficient
	ErrBalanceInsuffiecient = errors.New("user has insufficient balance")
	// ErrUserNotFound error when user data not found
	ErrUserNotFound = errors.New("user data not found")
	// ErrUnsignedSubsAgreementNotFound error when subscription agreement data not found
	ErrUnsignedSubsAgreementNotFound = errors.New("unsigned subscription agreement data not found")
	// ErrInvalidToken error when decrypt token is invalid
	ErrInvalidToken = errors.New("cipher: message authentication failed")
	// ErrInvalidSignature error when signature is invalid
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrUserIDNotPresentInHeader error when user id not in header
	ErrUserIDNotPresentInHeader = errors.New("unauthorized request")
)

// ErrUserRPCNotConnect erro when rpc user not connect
var (
	ErrUserRPCNotConnect     = errors.New("USER GRPC SERVER DOES NOT CONNECT")
	ErrActivityRPCNotConnect = errors.New("ACTIVITY GRPC SERVER DOES NOT CONNECT")
	ErrLoanRPCNotConnect     = errors.New("LOAN GRPC SERVER DOES NOT CONNECT")
)

// ErrUploadInvalid erro when rpc user not connect
var ErrUploadInvalid = errors.New("upload invalid")

// ErrPublishMessage error when failed publish message to kafka
var ErrPublishMessage = errors.New("failed publish message")

// ErrNoRowsMessage error when get data from database
var ErrNoRowsMessage = errors.New("sql: no rows in result set")
