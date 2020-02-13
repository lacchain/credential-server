package models

// MailRequest is the request to an API call.
type MailRequest struct {

	// folderId
	FolderId int32 `json:"folder_id"`

	// name
	Name string `json:"name"`

	// fileData
	FileData string `json:"file_data"`
}

// SendMailRequest is the request to Email API call
type SendMailRequest struct {

	// emailAddress
	EmailAddress string `json:"email_address"`

	// status
	Status string `json:"status"`

	// mergeFiles
	MergeFields MergeFields `json:"merge_fields"`
}

// MergeFields ...
type MergeFields struct {
	FNAME   string `json:"FNAME"`
	LNAME   string `json:"LNAME"`
	MMERGE5 string `json:"QRURL"`
	MMERGE6 string `json:"HASH"`
}
