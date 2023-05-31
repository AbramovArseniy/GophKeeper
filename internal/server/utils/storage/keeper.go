package storage

import (
	"bytes"
	"encoding/gob"
)

const (
	LoginPassword InfoType = "login-password"
	Card          InfoType = "card"
	Text          InfoType = "text"
	Binary        InfoType = "binary"
)

type InfoType string

func (i InfoType) String() string {
	switch i {
	case Card:
		return "Bank Card"
	case LoginPassword:
		return "Login/Password"
	case Text:
		return "Text"
	}
	return ""
}

type Info interface {
	MakeBinary() ([]byte, error)
}

type InfoMeta struct {
	Name  string   `json:"name"`
	Type  InfoType `json:"type"`
	Login string   `json:"user_login"`
}

type InfoLoginPass struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (p *InfoLoginPass) MakeBinary() ([]byte, error) {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(p)

	return buff.Bytes(), err
}

type InfoCard struct {
	CardNumber string `json:"card_number"`
	Holder     string `json:"holder"`
	Date       string `json:"exp_date"`
	CVCcode    string `json:"cvc"`
}

func (c *InfoCard) MakeBinary() ([]byte, error) {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(c)

	return buff.Bytes(), err
}

type InfoText struct {
	Text string `json:"text"`
}

func (t *InfoText) MakeBinary() ([]byte, error) {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(t)

	return buff.Bytes(), err
}

func NewInfo(infoType InfoType) Info {
	switch infoType {
	case LoginPassword:
		return &InfoLoginPass{}
	case Card:
		return &InfoCard{}
	case Text:
		return &InfoText{}
	}

	return nil
}
