package client

type Client struct {
	Id         int
	Name       string
	INN        string
	KPP        string
	OGRN       string
	Owner      string
	Phone      string
	Email      string
	Address    string
	CreateDate string
	Status     bool
}

type ClientObject struct {
	Id        int
	Equipment int
	Client    int
	Object    int
	Contract  int
}
