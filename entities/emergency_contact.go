package entities

type EmergencyContact struct {
	Relationship string
	FullName     string
	Emails       []string
	Phones       []string
}
