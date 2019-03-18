package main

type nginxMeta struct {
	Dest string `json:"dest"`
}

type nignxLevel struct {
	Name     string       `json:"name"`
	Children []nignxLevel `json:"children"`
}
