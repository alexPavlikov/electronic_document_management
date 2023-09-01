package requests

type Request struct {
	Id           string
	Title        string
	Description  string
	Priority     string
	StartDate    string
	EndDate      string
	Files        string
	Client       interface{}
	Worker       interface{}
	ClientObject interface{}
	Equipment    interface{}
	Contract     interface{}
	Status       interface{}
}
