package errx

import "errors"

// custom errors
var ErrReferrerNotFound = errors.New("referrer not found")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidUserID = errors.New("invalid user id")
var ErrEmptyRequestBody = errors.New("empty request body")

var ErrInvalidReferralCode = errors.New("invalid referral code")
var ErrIncorrectReferralCode = errors.New("incorrect referral code")

var ErrInvalidTaskType = errors.New("invalid task type")

var ErrTaskNotFound = errors.New("task not found")
var ErrNoChange = errors.New("no change")
// errors string representations
var ReferrerNotFound = ErrReferrerNotFound.Error()
var UserNotFound = ErrUserNotFound.Error()
var InvalidUserID = ErrInvalidUserID.Error()
var EmptyRequestBody = ErrEmptyRequestBody.Error()

var InvalidReferralCode = ErrInvalidReferralCode.Error()
var IncorrectReferralCode = ErrIncorrectReferralCode.Error()
var InvalidTaskType = ErrInvalidTaskType.Error()

var TaskNotFound = ErrTaskNotFound.Error()





type ErrorResponse struct {
	Error string `json:"error"`
}
