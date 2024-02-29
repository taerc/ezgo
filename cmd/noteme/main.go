package main

import "flag"

var AccessToken string
var AccessSecret string
var NoteText string
var NoteLink string
var ShowVersion bool

func init() {
	flag.BoolVar(&ShowVersion, "version", false, "print program build version")
	flag.StringVar(&AccessToken, "token", "", "path of configure file.")
	flag.StringVar(&AccessSecret, "secret", "", "path of configure file.")
	flag.StringVar(&NoteText, "text", "", "path of configure file.")
	flag.StringVar(&NoteLink, "link", "", "path of configure file.")
	flag.Parse()
}

func main() {

}
