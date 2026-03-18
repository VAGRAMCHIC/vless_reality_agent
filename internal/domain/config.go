package domain

type Config struct {
	Log struct {
		Loglevel string `json:"loglevel"`
	} `json:"log"`

	Inbounds []Inbound `json:"inbounds"`

	Outbounds []any `json:"outbounds"`
}

type Inbound struct {
	Tag      string `json:"tag"`
	Listen   string `json:"listen"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`

	Settings struct {
		Clients []Client `json:"clients"`
	} `json:"settings"`
}

type Client struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
