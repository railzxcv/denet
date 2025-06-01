package requests

type PostReferrer struct{
	ReferralCode string `json:"referralCode"`
}

type PostTaskComplete struct {
	TaskType string `json:"taskType"`
}
