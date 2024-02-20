package main

type StopCity struct {
	City      City
	State     State
	Latitude  float32
	Longitude float32
}

type City struct {
	CityId int
	Name   string
}

type State struct {
	Name    string
	Country string
}
