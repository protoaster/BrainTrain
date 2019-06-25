package main

import (
	"BrainTrain/app/controller"
	"net/http"
)

func main() {

	//Eingeloggte Seiten
	http.HandleFunc("/in_startseite", controller.In_startseite)
	http.HandleFunc("/in_karteikaesten", controller.In_karteikaesten)
	//	http.HandleFunc("/in_karteikaesten/kategorie", controller.In_karteikaesten)
	http.HandleFunc("/in_karteikarten_erstellen", controller.In_karteikarten_erstellen)
	http.HandleFunc("/in_karteikasten_anschauen", controller.In_karteikasten_anschauen)
	http.HandleFunc("/in_karteikasten_bearbeiten", controller.In_karteikasten_bearbeiten)
	http.HandleFunc("/in_lernen_antwort", controller.In_lernen_antwort)
	http.HandleFunc("/in_lernen_frage", controller.In_lernen_frage)
	http.HandleFunc("/in_meine_karteikaesten", controller.In_meine_karteikaesten)
	http.HandleFunc("/in_profil_popup", controller.In_profil_popup)
	http.HandleFunc("/in_profil", controller.In_profil)

	//Ausgeloggte Seiten
	http.HandleFunc("/", controller.Out_startseite)
	http.HandleFunc("/out_karteikaesten", controller.Out_karteikaesten)
	http.HandleFunc("/out_karteikasten_anschauen", controller.Out_karteikasten_anschauen)
	http.HandleFunc("/out_registrieren", controller.Out_registrieren)

	//bereitstellung der statischen inhalte
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/favicons/", http.StripPrefix("/favicons/", http.FileServer(http.Dir("./static/favicons"))))
	http.Handle("/font/", http.StripPrefix("/font/", http.FileServer(http.Dir("./static/font"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("./static/icons"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/logo/", http.StripPrefix("/logo/", http.FileServer(http.Dir("./static/logo"))))

	server := http.Server{
		Addr: ":8081",
	}

	server.ListenAndServe()

}
