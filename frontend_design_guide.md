# Complete Frontend Design Guide - Case Management System

## 1. **🏠 Main Dashboard/Cases Overview**

### **Screen Layout:**
- **Header:** Search bar, filters, "New Case" button
- **Stats Cards:** Total cases, open, overdue, closed this month
- **Cases Table/Grid:** Main content area

### **Data Needed:**
```json
GET /api/cases?page=1&limit=20&status=open&assigned_to=me
{
  "stats": {
    "total_cases": 145,
    "open_cases": 23, 
    "overdue_cases": 5,
    "closed_this_month": 12
  },
  "cases": [
    {
      "id": "case-123",
      "case_number": "CASE-2024-001",
      "title": "Housing Support Case",
      "priority": "HIGH",
      "status": "IN_PROGRESS", 
      "beneficiary": {
        "first_name": "John",
        "last_name": "Doe"
      },
      "assigned_to": {
        "first_name": "Jane",
        "last_name": "Smith"
      },
      "due_date": "2024-02-15",
      "notes_count": 5,
      "documents_count": 3,
      "last_activity": "2024-01-10T14:30:00Z",
      "created_at": "2024-01-01T09:00:00Z"
    }
  ]
}
```

### **Components to Build:**
- **📊 Stats Cards Component**
- **🔍 Search & Filter Bar**
  - Search by: case number, beneficiary name, title
  - Filters: Status, Priority, Case Type, Assigned To
- **📋 Cases Table/Grid Component**
  - Sortable columns
  - Status badges (color-coded)
  - Priority indicators
  - Quick actions (view, edit, close)
- **➕ New Case Button** → Opens create case modal

---

## 2. **➕ Create New Case Flow**

### **Step 1: Select Beneficiary Modal**
```json
GET /api/beneficiaries?search=john
{
  "beneficiaries": [
    {
      "id": "ben-456",
      "first_name": "John", 
      "last_name": "Doe",
      "phone": "+1234567890",
      "current_address": "123 Main St"
    }
  ]
}
```

### **Step 2: Case Details Form**
```json
POST /api/cases
{
  "beneficiary_id": "ben-456",
  "assigned_to_id": "user-789",
  "title": "Housing Application Support",
  "description": "Help with housing authority application process",
  "case_type": "HOUSING",
  "priority": "MEDIUM", 
  "due_date": "2024-02-15"
}
```

### **Components to Build:**
- **👤 Beneficiary Selector**
  - Searchable dropdown/modal
  - Shows beneficiary card with key info
- **📝 Case Form Component**
  - Title (required)
  - Description (required)
  - Case type dropdown (Housing, Legal, Medical, Support, Other)
  - Priority selector (Low, Medium, High, Urgent)
  - Assigned to dropdown (organization users)
  - Due date picker
- **✅ Submit & Navigate** → Redirect to case details

---

## 3. **📄 Individual Case Details Screen**

### **Screen Layout:**
- **Header Section:** Case info & actions
- **Tabs:** Overview, Notes, Documents, Activity Log
- **Right Sidebar:** Beneficiary info, quick actions

### **Main Case Data:**
```json
GET /api/cases/case-123
{
  "id": "case-123",
  "case_number": "CASE-2024-001", 
  "title": "Housing Support Case",
  "description": "Help with housing application...",
  "status": "IN_PROGRESS",
  "priority": "HIGH",
  "case_type": "HOUSING",
  "beneficiary": {
    "id": "ben-456",
    "first_name": "John",
    "last_name": "Doe", 
    "phone": "+1234567890",
    "email": "john@example.com",
    "current_address": "123 Main St, City"
  },
  "assigned_to": {
    "id": "user-789",
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane@org.com"
  },
  "due_date": "2024-02-15",
  "notes_count": 8,
  "documents_count": 5,
  "created_at": "2024-01-01T09:00:00Z",
  "updated_at": "2024-01-10T14:30:00Z"
}
```

### **Components to Build:**
- **📌 Case Header Component**
  - Case number + title
  - Status badge (editable)
  - Priority indicator (editable)
  - Actions menu (Edit, Close, Reassign, Delete)
- **👤 Beneficiary Info Card**
  - Photo, name, contact info
  - Link to full beneficiary profile
- **📊 Case Stats Cards**
  - Notes count, documents count, days open
- **🔄 Status Update Component**
  - Quick status change dropdown
  - Auto-saves on change

---

## 4. **📝 Notes Tab/Section**

### **Recent Notes Data:**
```json
GET /api/cases/case-123/notes?limit=10
{
  "notes": [
    {
      "id": "note-123",
      "title": "Follow-up call completed",
      "content": "Spoke with beneficiary about housing application status. Documents submitted successfully to housing authority. Expecting response within 2 weeks.",
      "tags": ["follow-up", "phone-call", "housing"],
      "note_type": "CALL",
      "is_important": true,
      "created_by": {
        "name": "Jane Smith"
      },
      "created_at": "2024-01-10T14:30:00Z"
    },
    {
      "id": "note-124", 
      "title": "Document review",
      "content": "Reviewed submitted documents. All forms complete.",
      "tags": ["documents", "review"],
      "note_type": "UPDATE",
      "is_important": false,
      "created_by": {
        "name": "John Admin"
      },
      "created_at": "2024-01-08T11:20:00Z"
    }
  ]
}
```

### **Add Note Flow:**
```json
POST /api/cases/case-123/notes
{
  "title": "Scheduled appointment",
  "content": "Appointment scheduled for housing interview on 2024-01-20 at 2 PM",
  "tags": ["appointment", "housing", "interview"],
  "note_type": "APPOINTMENT",
  "is_important": true
}
```

### **Components to Build:**
- **📝 Add Note Form**
  - Title (required)
  - Content textarea (required)
  - Tags input (multi-select/chips)
  - Note type dropdown (Call, Meeting, Update, Appointment, Other)
  - Important flag checkbox
- **📋 Notes Timeline Component**
  - Chronological list
  - Expandable content
  - Tag badges
  - Important flag indicator
  - Edit/delete actions (for note creator)
- **🏷️ Tags Filter**
  - Filter notes by tags
  - Show all available tags

---

## 5. **📎 Documents Tab/Section**

### **Documents Data:**
```json
GET /api/cases/case-123/documents
{
  "documents": [
    {
      "id": "doc-123",
      "document_name": "Housing Application Form",
      "file_name": "housing_app_john_doe.pdf", 
      "document_type": "FORM",
      "file_size": 245760,
      "mime_type": "application/pdf",
      "description": "Completed housing application with all signatures",
      "tags": ["housing", "application", "signed"],
      "uploaded_by": {
        "name": "Jane Smith"
      },
      "created_at": "2024-01-08T16:20:00Z",
      "download_url": "/api/documents/doc-123/download"
    },
    {
      "id": "doc-124",
      "document_name": "ID Copy",
      "file_name": "john_doe_id.jpg",
      "document_type": "IDENTIFICATION", 
      "file_size": 102400,
      "mime_type": "image/jpeg",
      "description": "Copy of beneficiary's ID card",
      "tags": ["identification", "photo"],
      "uploaded_by": {
        "name": "Jane Smith"
      },
      "created_at": "2024-01-05T11:15:00Z",
      "download_url": "/api/documents/doc-124/download"
    }
  ]
}
```

### **Upload Document Flow:**
```json
POST /api/cases/case-123/documents
Content-Type: multipart/form-data

{
  "document_name": "Medical Report",
  "document_type": "MEDICAL", 
  "description": "Latest medical assessment from clinic",
  "tags": ["medical", "assessment", "clinic"],
  "file": [binary file data]
}
```

### **Components to Build:**
- **📤 Upload Document Component**
  - Drag & drop area
  - File picker button
  - Document name input (required)
  - Document type dropdown:
    - 📋 FORM (Applications, intake forms)
    - 📄 REPORT (Assessments, evaluations)
    - 📸 EVIDENCE (Photos, supporting materials)
    - 📧 CORRESPONDENCE (Emails, letters)
    - 🆔 IDENTIFICATION (ID copies, documents)
    - 💼 LEGAL (Contracts, legal documents)
    - 🏥 MEDICAL (Medical records, reports)
    - 📊 OTHER (Miscellaneous)
  - Description textarea
  - Tags input
  - Upload progress bar
- **📁 Documents Grid/List Component**
  - File type icons
  - Document name + description
  - File size display
  - Upload date & user
  - Tags badges
  - Actions: Download, Preview, Edit, Delete
- **👁️ Document Preview Component**
  - PDF viewer for PDFs
  - Image viewer for photos
  - Download button
- **🔍 Documents Filter**
  - Filter by document type
  - Filter by tags
  - Search by document name

---

## 6. **⚙️ Edit Case Modal/Screen**

### **Edit Case Data:**
```json
PUT /api/cases/case-123
{
  "title": "Updated Housing Support Case",
  "description": "Updated description...",
  "case_type": "HOUSING",
  "priority": "HIGH", 
  "status": "IN_PROGRESS",
  "assigned_to_id": "user-789",
  "due_date": "2024-03-01"
}
```

### **Components to Build:**
- **📝 Edit Case Form** (similar to create, but pre-filled)
- **🔄 Status Dropdown** with visual indicators
- **👥 Reassign User Dropdown**
- **📅 Due Date Picker**

---

## 7. **📊 Reports/Analytics Screen**

### **Analytics Data:**
```json
GET /api/reports/dashboard
{
  "overview": {
    "total_cases": 145,
    "open_cases": 23,
    "in_progress": 18,
    "overdue_cases": 5,
    "closed_this_month": 12,
    "avg_resolution_days": 14
  },
  "by_priority": {
    "HIGH": 8,
    "MEDIUM": 15, 
    "LOW": 22,
    "URGENT": 2
  },
  "by_type": {
    "HOUSING": 18,
    "LEGAL": 12,
    "MEDICAL": 8,
    "SUPPORT": 15
  },
  "by_user": [
    {
      "user_name": "Jane Smith",
      "active_cases": 8,
      "closed_cases": 12
    }
  ]
}
```

### **Components to Build:**
- **📈 Dashboard Charts**
  - Cases by status (pie chart)
  - Cases by priority (bar chart)
  - Cases by type (donut chart)
  - Cases over time (line chart)
- **📋 Performance Tables**
  - Cases by user
  - Overdue cases
  - Recently closed cases

---

## **🗂️ Navigation Structure:**

```
Main App Layout
├── 🏠 Dashboard (Cases overview)
├── ➕ Create Case (Modal/Page)
├── 📄 Case Details
│   ├── Overview Tab
│   ├── 📝 Notes Tab  
│   ├── 📎 Documents Tab
│   └── 📊 Activity Log Tab
├── 📊 Reports
└── ⚙️ Settings
```

## **🎨 UI/UX Design Considerations:**

### **Color Coding:**
- **🔴 HIGH/URGENT:** Red
- **🟡 MEDIUM:** Yellow/Orange  
- **🟢 LOW:** Green
- **⚫ Status:** Open (blue), In Progress (orange), Closed (green)

### **Icons to Use:**
- 📋 Cases, 📝 Notes, 📎 Documents
- 👤 Users, 🏠 Housing, ⚖️ Legal, 🏥 Medical
- ⬆️ Priority, 📅 Dates, 🏷️ Tags

### **Mobile Responsiveness:**
- Stack components vertically on mobile
- Collapsible sidebar
- Touch-friendly buttons and inputs

## **API Endpoints Summary:**

### **Cases:**
- `GET /api/cases` - List cases with filters
- `POST /api/cases` - Create new case
- `GET /api/cases/{id}` - Get case details
- `PUT /api/cases/{id}` - Update case
- `DELETE /api/cases/{id}` - Delete case

### **Notes:**
- `GET /api/cases/{case_id}/notes` - Get case notes
- `POST /api/cases/{case_id}/notes` - Add note
- `PUT /api/notes/{note_id}` - Update note
- `DELETE /api/notes/{note_id}` - Delete note

### **Documents:**
- `GET /api/cases/{case_id}/documents` - Get case documents
- `POST /api/cases/{case_id}/documents` - Upload document
- `GET /api/documents/{doc_id}/download` - Download document
- `PUT /api/documents/{doc_id}` - Update document metadata
- `DELETE /api/documents/{doc_id}` - Delete document

### **Reports:**
- `GET /api/reports/dashboard` - Get dashboard analytics
- `GET /api/reports/cases` - Get case reports

This complete guide provides everything needed to design a comprehensive case management frontend! 🚀 