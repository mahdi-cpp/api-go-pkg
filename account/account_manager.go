package account

import (
	"github.com/mahdi-cpp/api-go-pkg/network"
	"log"
)

//type UserPHCollection struct {
//	Collections []*asset.PHCollection[User] `json:"collections"`
//}

type Manager struct {
	networkUser     *network.Control[User]
	networkUserList *network.Control[[]User]
	//networkUsers *network.Control[UserPHCollection]
}

type requestBody struct {
	ID int `json:"id"`
}

func NewAccountManager() *Manager {
	manager := &Manager{
		networkUser:     network.NewNetworkManager[User]("http://localhost:8080/api/v1/user/get_user"),
		networkUserList: network.NewNetworkManager[[]User]("http://localhost:8080/api/v1/user/list"),
	}

	return manager
}

func (m *Manager) GetUser(id int) (*User, error) {

	user, err := m.networkUser.Read("", requestBody{ID: id})
	if err != nil {
		log.Fatalf("Error: %v", err)
		return nil, err
	}
	return user, nil
}

func (m *Manager) GetAll() (*[]User, error) {

	users, err := m.networkUserList.Read("", nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return nil, err
	}
	return users, nil
}
