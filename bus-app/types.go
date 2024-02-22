package main

type StopCity struct {
	State     State
	City      City
	Latitude  float32
	Longitude float32
}

type City struct {
	Name   string
	CityId int
}

type State struct {
	Name    string
	Country string
}
