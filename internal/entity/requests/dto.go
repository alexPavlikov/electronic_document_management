package requests

import (
	"github.com/alexPavlikov/electronic_document_management/internal/entity/client"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/contract"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/equipment"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/user"
)

var (
	ClientsDTO      []client.Client
	ContractsDTO    []contract.Contract
	WorkerDTO       []user.User
	ClientObjectDTO []client.ClientObject
	EquipmentDTO    []equipment.Equipment
	StatusDTO       []ReqStatus
)
