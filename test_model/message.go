package test_model

import "time"

func (a *Message) SetID(id int)                    { a.ID = id }
func (a *Message) SetCreationDate(t time.Time)     { a.CreationDate = t }
func (a *Message) SetModificationDate(t time.Time) { a.ModificationDate = t }
func (a *Message) GetID() int                      { return a.ID }
func (a *Message) GetCreationDate() time.Time      { return a.CreationDate }
func (a *Message) GetModificationDate() time.Time  { return a.ModificationDate }

type Message struct {
	ID      int    `json:"id"`
	ChatID  int    `json:"chatId"`  // Identifier of the chat where the message belongs
	UserID  int    `json:"userId"`  // Identifier of the user who sent the message
	Content string `json:"content"` // The actual message content
	Type    string `json:"type"`

	Assets []MessageAsset `json:"assets"`

	//Date------------------------------------------
	ExpiryDate       time.Time `json:"expiryDate"`
	CreationDate     time.Time `json:"creationDate"`
	ModificationDate time.Time `json:"modificationDate"`
}

type MessageAsset struct {
	ID     int `json:"id"`
	ChatID int `json:"chatId"`
	UserID int `json:"userId"`
}
