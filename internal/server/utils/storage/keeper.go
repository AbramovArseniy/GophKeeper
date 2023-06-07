package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

const (
	LoginPassword InfoType = iota + 1
	Card
	Text
)

type InfoType int

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
	DecodeBinary(bin []byte) error
}

type InfoMeta struct {
	Name  string
	Type  InfoType
	Login string
}

type InfoLoginPass struct {
	Login    string
	Password string
}

func (p *InfoLoginPass) MakeBinary() ([]byte, error) {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(p)

	return buff.Bytes(), err
}

func (p *InfoLoginPass) DecodeBinary(bin []byte) error {
	var buff bytes.Buffer
	_, err := buff.Write(bin)
	if err != nil {
		return fmt.Errorf("error while writing binary into buffer: %w", err)
	}
	dec := gob.NewDecoder(&buff)
	err = dec.Decode(p)
	if err != nil {
		return fmt.Errorf("error while decoding binary: %w", err)
	}
	return err
}

type InfoCard struct {
	CardNumber string
	Holder     string
	Date       string
	CVCcode    string
}

func (c *InfoCard) MakeBinary() ([]byte, error) {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(c)

	return buff.Bytes(), err
}

func (c *InfoCard) DecodeBinary(bin []byte) error {
	var buff bytes.Buffer
	_, err := buff.Write(bin)
	if err != nil {
		return fmt.Errorf("error while writing binary into buffer: %w", err)
	}
	dec := gob.NewDecoder(&buff)
	err = dec.Decode(c)
	if err != nil {
		return fmt.Errorf("error while decoding binary: %w", err)
	}
	return err
}

type InfoText struct {
	Text string
}

func (t *InfoText) MakeBinary() ([]byte, error) {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(t)

	return buff.Bytes(), err
}

func (t *InfoText) DecodeBinary(bin []byte) error {
	var buff bytes.Buffer
	_, err := buff.Write(bin)
	if err != nil {
		return fmt.Errorf("error while writing binary into buffer: %w", err)
	}
	dec := gob.NewDecoder(&buff)
	err = dec.Decode(t)
	if err != nil {
		return fmt.Errorf("error while decoding binary: %w", err)
	}
	return err
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
