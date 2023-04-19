package swag

type SwagTermOfService struct {
	URL string // no prop
}

type SwagContact struct {
	Email string
}

type SwagLicence struct {
	Name string
	URL  string
}

type SwagInfo struct {
	Title         string
	Version       string
	Description   string
	TermOfService *SwagTermOfService // no prop
	Contact       *SwagContact
	License       *SwagLicence
}
