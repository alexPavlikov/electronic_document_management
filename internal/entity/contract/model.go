package contract

type Contract struct {
	Id        int
	Name      string
	Client    interface{}
	DataStart string
	DataEnd   string
	Amount    int
	File      string
	Status    bool
}
