package responses

import (
	"time"

	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
)

type BeneficiaryResponse struct {
	ID             string  `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	FullName       string  `json:"full_name"`
	Phone          *string `json:"phone,omitempty"`
	Email          *string `json:"email,omitempty"`
	CurrentAddress *string `json:"current_address,omitempty"`
	ImageURL       *string `json:"image_url,omitempty"`
}

type AssignedToResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CreatedByResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UploadedByResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CaseResponse struct {
	ID                string              `json:"id"`
	CaseNumber        string              `json:"case_number"`
	Title             string              `json:"title"`
	Description       string              `json:"description"`
	Status            string              `json:"status"`
	Priority          string              `json:"priority"`
	UrgencyLevel      *string             `json:"urgency_level,omitempty"`
	CaseType          string              `json:"case_type"`
	BeneficiaryID     string              `json:"beneficiary_id"`
	Beneficiary       BeneficiaryResponse `json:"beneficiary"`
	AssignedToID      string              `json:"assigned_to_id"`
	AssignedTo        AssignedToResponse  `json:"assigned_to"`
	DueDate           *string             `json:"due_date,omitempty"`
	EstimatedDuration *string             `json:"estimated_duration,omitempty"`
	BudgetAllocated   *string             `json:"budget_allocated,omitempty"`
	Tags              []string            `json:"tags"`
	NotesCount        int                 `json:"notes_count"`
	DocumentsCount    int                 `json:"documents_count"`
	LastActivity      string              `json:"last_activity"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
}

type CaseNoteResponse struct {
	ID        string            `json:"id"`
	CaseID    string            `json:"case_id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Tags      []string          `json:"tags"`
	NoteType  string            `json:"note_type"`
	Important bool              `json:"is_important"`
	CreatedBy CreatedByResponse `json:"created_by"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type CaseDocumentResponse struct {
	ID           string             `json:"id"`
	CaseID       string             `json:"case_id"`
	DocumentName string             `json:"document_name"`
	FileName     string             `json:"file_name"`
	DocumentType string             `json:"document_type"`
	FileSize     int64              `json:"file_size"`
	MimeType     string             `json:"mime_type"`
	Description  string             `json:"description"`
	Tags         []string           `json:"tags"`
	UploadedBy   UploadedByResponse `json:"uploaded_by"`
	CreatedAt    string             `json:"created_at"`
	DownloadURL  string             `json:"download_url"`
}

type CaseStatsResponse struct {
	TotalCases        int `json:"total_cases"`
	OpenCases         int `json:"open_cases"`
	InProgressCases   int `json:"in_progress_cases"`
	OverdueCases      int `json:"overdue_cases"`
	ClosedThisMonth   int `json:"closed_this_month"`
	AvgResolutionDays int `json:"avg_resolution_days"`
}

type CaseListResponse struct {
	Count int            `json:"count"`
	Data  []CaseResponse `json:"data"`
}

type CaseNoteListResponse struct {
	Count int                `json:"count"`
	Data  []CaseNoteResponse `json:"data"`
}

type CaseDocumentListResponse struct {
	Count int                    `json:"count"`
	Data  []CaseDocumentResponse `json:"data"`
}

// Conversion functions
func NewCaseResponse(c entities.Case) CaseResponse {
	// Handle beneficiary name parsing
	var firstName, lastName string
	if c.Beneficiary.FullName != "" {
		names := utils.ParseFullName(c.Beneficiary.FullName)
		firstName = names.FirstName
		lastName = names.LastName
	}

	// Handle optional fields
	var phone, email, currentAddress, imageURL *string
	if c.Beneficiary.Email != "" {
		email = &c.Beneficiary.Email
	}
	if len(c.Beneficiary.Phones) > 0 {
		phone = &c.Beneficiary.Phones[0]
	}
	if c.Beneficiary.Address.AddressLine1 != "" {
		addr := c.Beneficiary.Address.AddressLine1 + ", " + c.Beneficiary.Address.City
		currentAddress = &addr
	}
	if c.Beneficiary.ImageURL != "" {
		imageURL = &c.Beneficiary.ImageURL
	}

	// Handle optional case fields
	var dueDate, estimatedDuration, budgetAllocated, urgencyLevel *string
	if c.DueDate != nil {
		date := c.DueDate.Format(time.RFC3339)
		dueDate = &date
	}
	if c.EstimatedDuration != "" {
		estimatedDuration = &c.EstimatedDuration
	}
	if c.BudgetAllocated != "" {
		budgetAllocated = &c.BudgetAllocated
	}
	if c.UrgencyLevel != "" {
		urgencyLevel = &c.UrgencyLevel
	}

	return CaseResponse{
		ID:            c.ID,
		CaseNumber:    c.CaseNumber,
		Title:         c.Title,
		Description:   c.Description,
		Status:        c.Status,
		Priority:      c.Priority,
		UrgencyLevel:  urgencyLevel,
		CaseType:      c.CaseType,
		BeneficiaryID: c.BeneficiaryID,
		Beneficiary: BeneficiaryResponse{
			ID:             c.Beneficiary.ID,
			FirstName:      firstName,
			LastName:       lastName,
			FullName:       c.Beneficiary.FullName,
			Phone:          phone,
			Email:          email,
			CurrentAddress: currentAddress,
			ImageURL:       imageURL,
		},
		AssignedToID: c.AssignedToID,
		AssignedTo: AssignedToResponse{
			ID:        c.AssignedTo.ID,
			FirstName: c.AssignedTo.FirstName,
			LastName:  c.AssignedTo.LastName,
			Email:     c.AssignedTo.Email,
		},
		DueDate:           dueDate,
		EstimatedDuration: estimatedDuration,
		BudgetAllocated:   budgetAllocated,
		Tags:              c.Tags,
		NotesCount:        c.NotesCount,
		DocumentsCount:    c.DocumentsCount,
		LastActivity:      c.LastActivity.Format(time.RFC3339),
		CreatedAt:         c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         c.UpdatedAt.Format(time.RFC3339),
	}
}

func NewCaseNoteResponse(n entities.CaseNote) CaseNoteResponse {
	return CaseNoteResponse{
		ID:        n.ID,
		CaseID:    n.CaseID,
		Title:     n.Title,
		Content:   n.Content,
		Tags:      n.Tags,
		NoteType:  n.NoteType,
		Important: n.IsImportant,
		CreatedBy: CreatedByResponse{
			ID:   n.CreatedBy.ID,
			Name: n.CreatedBy.FullName,
		},
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
		UpdatedAt: n.UpdatedAt.Format(time.RFC3339),
	}
}

func NewCaseDocumentResponse(d entities.CaseDocument) CaseDocumentResponse {
	return CaseDocumentResponse{
		ID:           d.ID,
		CaseID:       d.CaseID,
		DocumentName: d.DocumentName,
		FileName:     d.FileName,
		DocumentType: d.DocumentType,
		FileSize:     d.FileSize,
		MimeType:     d.MimeType,
		Description:  d.Description,
		Tags:         d.Tags,
		UploadedBy: UploadedByResponse{
			ID:   d.UploadedBy.ID,
			Name: d.UploadedBy.FullName,
		},
		CreatedAt:   d.CreatedAt.Format(time.RFC3339),
		DownloadURL: d.DownloadURL,
	}
}
