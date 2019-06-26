package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type tmp_b_home struct {
	Nutzer     string
	Lernkarten string
	Karteien   string
}

type tmp_in struct {
	Nutzer        string
	Lernkarten    string
	Karteien      string
	MeineKarteien string
	Nutzername    string
}

type in_menu struct {
	Karteien      string
	MeineKarteien string
	Nutzername    string
}

type out_menu struct {
	Karteien string
}

type in_erstellen struct {
	Karteien              string
	MeineKarteien         string
	Nutzername            string
	Naturwissenschaften   []Karteikasten
	Sprachen              []Karteikasten
	Gesellschaft          []Karteikasten
	Wirtschaft            []Karteikasten
	Geisteswissenschaften []Karteikasten
	Sonstige              []Karteikasten
}

type tmp_nL_Karteikasten struct {
	Karteien              string
	MeineKarteien         string
	Nutzername            string
	Abfrage               string
	Naturwissenschaften   []Karteikasten
	Sprachen              []Karteikasten
	Gesellschaft          []Karteikasten
	Wirtschaft            []Karteikasten
	Geisteswissenschaften []Karteikasten
	Sonstige              []Karteikasten
}

type tmp_L_MeineKarteikaesten struct {
	Karteien                  string
	MeineKarteien             string
	Nutzername                string
	GespeicherteKarteikaesten []Karteikasten
	MeineKarteikaesten        []Karteikasten
}

type tmp_L_modkarteikasten1 struct {
	MeineKarteien         string
	Karteien              string
	Ersteller             string
	Nutzername            string
	AlleKarten            []Karte
	AlleFortschritte      []int
	AktuelleKarte         Karte
	AktuellerKarteikasten Karteikasten
	Fach                  [5]int
}

func Out_startseite(w http.ResponseWriter, r *http.Request) {
	p := tmp_b_home{Nutzer: strconv.Itoa(GetNutzeranz()), Lernkarten: strconv.Itoa(GetKartenAnz()), Karteien: strconv.Itoa(GetKarteikastenAnz())}
	t, _ := template.ParseFiles("./templates/out_menu.html", "./templates/out_startseite.html")
	t.ExecuteTemplate(w, "layout", p)
}

func Out_karteikaesten(w http.ResponseWriter, r *http.Request) {

	data := tmp_nL_Karteikasten{
		Karteien:              strconv.Itoa(GetKarteikastenAnz()),
		Naturwissenschaften:   []Karteikasten{},
		Sprachen:              []Karteikasten{},
		Gesellschaft:          []Karteikasten{},
		Wirtschaft:            []Karteikasten{},
		Geisteswissenschaften: []Karteikasten{},
		Sonstige:              []Karteikasten{},
	}

	kk := []Karteikasten{}
	kk = GetAlleKarteikaesten()

	for _, element := range kk {
		if element.Kategorie == "Naturwissenschaften" {
			data.Naturwissenschaften = append(data.Naturwissenschaften, element)
		} else if element.Kategorie == "Sprachen" {
			data.Sprachen = append(data.Sprachen, element)
		} else if element.Kategorie == "Gesellschaft" {
			data.Gesellschaft = append(data.Gesellschaft, element)
		} else if element.Kategorie == "Wirtschaft" {
			data.Wirtschaft = append(data.Wirtschaft, element)
		} else if element.Kategorie == "Geisteswissenschaften" {
			data.Geisteswissenschaften = append(data.Geisteswissenschaften, element)
		} else {
			data.Sonstige = append(data.Sonstige, element)
		}
	}

	t, _ := template.ParseFiles("./templates/out_menu.html", "./templates/out_karteikaesten.html")
	t.ExecuteTemplate(w, "layout", data)
}

func Out_karteikasten_anschauen(w http.ResponseWriter, r *http.Request) {
	data := tmp_L_modkarteikasten1{
		Karteien:              strconv.Itoa(GetKarteikastenAnz()),
		Ersteller:             "",
		AktuellerKarteikasten: Karteikasten{},
		AlleKarten:            []Karte{},
		AlleFortschritte:      []int{},
		AktuelleKarte:         Karte{},
	}

	temp_kk := GetKarteikastenByid(1)
	temp_kk.FortschrittP = int(GetKarteikastenFortschritt(GetKarteikastenByid(1), GetNutzerById(1)))

	data.Ersteller = temp_kk.Ersteller

	//gewählte Karte

	Num := r.FormValue("Num")
	if Num == "" {
		Num = "1"
	}

	data.AktuellerKarteikasten = temp_kk

	for i, element := range temp_kk.Karten {
		data.AlleKarten = append(data.AlleKarten, element)
		data.AlleKarten[i].Num = i + 1
	}

	for _, element := range GetKarteikastenWiederholungArr(temp_kk, GetNutzerById(1)) {
		data.AlleFortschritte = append(data.AlleFortschritte, element)
	}

	akt, _ := strconv.Atoi(Num)
	akt = akt - 1

	//fmt.Println("#########################################################################################")
	//fmt.Println(akt)
	//fmt.Println("#########################################################################################")
	data.AktuelleKarte = data.AlleKarten[akt]
	data.AktuelleKarte.NutzerFach = strconv.Itoa(data.AlleFortschritte[akt])

	fmt.Println(data.AktuelleKarte.NutzerFach)

	t, _ := template.ParseFiles("./templates/out_menu.html", "./templates/out_karteikasten_anschauen.html")
	t.ExecuteTemplate(w, "layout", data)
}

func Out_registrieren(w http.ResponseWriter, r *http.Request) {
	data := out_menu{Karteien: strconv.Itoa(GetKarteikastenAnz())}

	t, _ := template.ParseFiles("./templates/out_menu.html", "./templates/out_registrieren.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_startseite(w http.ResponseWriter, r *http.Request) {
	data := tmp_in{
		Nutzer:        strconv.Itoa(GetNutzeranz()),
		Lernkarten:    strconv.Itoa(GetKartenAnz()),
		Karteien:      strconv.Itoa(GetKarteikastenAnz()),
		MeineKarteien: "",
		Nutzername:    "",
	}
	tempMK := GetNutzerById(1)
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername

	fmt.Println(data.MeineKarteien)

	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_startseite.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_karteikaesten(w http.ResponseWriter, r *http.Request) {

	data := tmp_nL_Karteikasten{
		MeineKarteien:         "",
		Nutzername:            "",
		Karteien:              strconv.Itoa(GetKarteikastenAnz()),
		Naturwissenschaften:   []Karteikasten{},
		Sprachen:              []Karteikasten{},
		Gesellschaft:          []Karteikasten{},
		Wirtschaft:            []Karteikasten{},
		Geisteswissenschaften: []Karteikasten{},
		Sonstige:              []Karteikasten{},
		Abfrage:               "Alle",
	}
	test := ""
	kk := []Karteikasten{}
	kk = GetAlleKarteikaesten()

	if r.Method == "GET" {

		r.ParseForm()
		test = r.FormValue("kategorie")
		fmt.Println(test)
	}

	data.Abfrage = test

	for _, element := range kk {
		if element.Kategorie == "Naturwissenschaften" && (test == "Naturwissenschaften" || test == "Alle") {
			data.Naturwissenschaften = append(data.Naturwissenschaften, element)

		} else if element.Kategorie == "Sprachen" && (test == "Sprachen" || test == "Alle") {
			data.Sprachen = append(data.Sprachen, element)

		} else if element.Kategorie == "Gesellschaft" && (test == "Gesellschaft" || test == "Alle") {
			data.Gesellschaft = append(data.Gesellschaft, element)

		} else if element.Kategorie == "Wirtschaft" && (test == "Wirtschaft" || test == "Alle") {
			data.Wirtschaft = append(data.Wirtschaft, element)

		} else if element.Kategorie == "Geisteswissenschaften" && (test == "Geisteswissenschaften" || test == "Alle") {
			data.Geisteswissenschaften = append(data.Geisteswissenschaften, element)

		} else if element.Kategorie == "Sonstige" && (test == "Sonstige" || test == "Alle") {
			data.Sonstige = append(data.Sonstige, element)

		}
	}

	tempMK := GetNutzerById(1)
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername

	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_karteikaesten.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_karteikarten_erstellen(w http.ResponseWriter, r *http.Request) {
	data := in_erstellen{
		MeineKarteien: "",
		Nutzername:    "",
		Karteien:      strconv.Itoa(GetKarteikastenAnz()),
	}
	tempMK := GetNutzerById(1)

	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername

	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_karteikarten_erstellen.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_lernen_antwort(w http.ResponseWriter, r *http.Request) {
	data := tmp_L_modkarteikasten1{
		MeineKarteien:         "",
		Nutzername:            "",
		Karteien:              strconv.Itoa(GetKarteikastenAnz()),
		AktuellerKarteikasten: Karteikasten{},
		AlleKarten:            []Karte{},
		AlleFortschritte:      []int{},
		AktuelleKarte:         Karte{},
		Fach:                  [5]int{},
	}
	temp_kk := GetKarteikastenByid(1)
	temp_kk.FortschrittP = int(GetKarteikastenFortschritt(GetKarteikastenByid(1), GetNutzerById(1)))

	data.AktuellerKarteikasten = temp_kk

	for i, element := range temp_kk.Karten {
		data.AlleKarten = append(data.AlleKarten, element)
		data.AlleKarten[i].Num = i + 1
	}

	for _, element := range GetKarteikastenWiederholungArr(temp_kk, GetNutzerById(1)) {
		data.AlleFortschritte = append(data.AlleFortschritte, element)
	}

	for n := 0; n <= 4; n++ {
		data.Fach[n] = GetKarteikartenAnzByFach(temp_kk, n, GetNutzerById(1))
	}

	fmt.Println(data.Fach)

	data.AktuelleKarte = data.AlleKarten[0]
	data.AktuelleKarte.NutzerFach = strconv.Itoa(data.AlleFortschritte[0])

	tempMK := GetNutzerById(1)
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername

	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_lernen_antwort.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_lernen_frage(w http.ResponseWriter, r *http.Request) {
	data := tmp_L_modkarteikasten1{
		MeineKarteien:         "",
		Nutzername:            "",
		Karteien:              strconv.Itoa(GetKarteikastenAnz()),
		AktuellerKarteikasten: Karteikasten{},
		AlleKarten:            []Karte{},
		AlleFortschritte:      []int{},
		AktuelleKarte:         Karte{},
		Fach:                  [5]int{},
	}
	temp_kk := GetKarteikastenByid(1)
	temp_kk.FortschrittP = int(GetKarteikastenFortschritt(GetKarteikastenByid(1), GetNutzerById(1)))

	data.AktuellerKarteikasten = temp_kk

	for i, element := range temp_kk.Karten {
		data.AlleKarten = append(data.AlleKarten, element)
		data.AlleKarten[i].Num = i + 1
	}

	for _, element := range GetKarteikastenWiederholungArr(temp_kk, GetNutzerById(1)) {
		data.AlleFortschritte = append(data.AlleFortschritte, element)
	}

	for n := 0; n <= 4; n++ {
		data.Fach[n] = GetKarteikartenAnzByFach(temp_kk, n, GetNutzerById(1))
	}

	fmt.Println(data.Fach)

	data.AktuelleKarte = data.AlleKarten[0]
	data.AktuelleKarte.NutzerFach = strconv.Itoa(data.AlleFortschritte[0])

	tempMK := GetNutzerById(1)
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername

	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_lernen_frage.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_karteikasten_anschauen(w http.ResponseWriter, r *http.Request) {
	data := tmp_L_modkarteikasten1{
		Nutzername:            "",
		MeineKarteien:         "",
		Karteien:              strconv.Itoa(GetKarteikastenAnz()),
		AktuellerKarteikasten: Karteikasten{},
		AlleKarten:            []Karte{},
		AlleFortschritte:      []int{},
		AktuelleKarte:         Karte{},
	}

	temp_kk := GetKarteikastenByid(1)
	temp_kk.FortschrittP = int(GetKarteikastenFortschritt(GetKarteikastenByid(1), GetNutzerById(1)))

	data.Ersteller = temp_kk.Ersteller

	//gewählte Karte

	Num := r.FormValue("Num")
	if Num == "" {
		Num = "1"
	}

	data.AktuellerKarteikasten = temp_kk

	for i, element := range temp_kk.Karten {
		data.AlleKarten = append(data.AlleKarten, element)
		data.AlleKarten[i].Num = i + 1
	}

	for _, element := range GetKarteikastenWiederholungArr(temp_kk, GetNutzerById(1)) {
		data.AlleFortschritte = append(data.AlleFortschritte, element)
	}

	akt, _ := strconv.Atoi(Num)
	akt = akt - 1

	//fmt.Println("#########################################################################################")
	//fmt.Println(akt)
	//fmt.Println("#########################################################################################")
	data.AktuelleKarte = data.AlleKarten[akt]
	data.AktuelleKarte.NutzerFach = strconv.Itoa(data.AlleFortschritte[akt])

	fmt.Println(data.AktuelleKarte.NutzerFach)

	tempMK := GetNutzerById(1)
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername

	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_karteikasten_anschauen.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_karteikasten_bearbeiten(w http.ResponseWriter, r *http.Request) {
	data := tmp_L_modkarteikasten1{
		Nutzername:            "",
		MeineKarteien:         "",
		Karteien:              strconv.Itoa(GetKarteikastenAnz()),
		AktuellerKarteikasten: Karteikasten{},
		AlleKarten:            []Karte{},
	}

	temp_kk := GetKarteikastenByid(1)
	temp_kk.FortschrittP = int(GetKarteikastenFortschritt(GetKarteikastenByid(1), GetNutzerById(1)))

	data.AktuellerKarteikasten = temp_kk

	for i, element := range temp_kk.Karten {
		data.AlleKarten = append(data.AlleKarten, element)
		data.AlleKarten[i].Num = i + 1
	}

	tempMK := GetNutzerById(1)
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername

	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_karteikasten_bearbeiten.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_meine_karteikaesten(w http.ResponseWriter, r *http.Request) {
	data := tmp_L_MeineKarteikaesten{
		Nutzername:                "",
		MeineKarteien:             "",
		Karteien:                  "",
		GespeicherteKarteikaesten: []Karteikasten{},
		MeineKarteikaesten:        []Karteikasten{},
	}
	nutzer := GetNutzerById(1) //muss noch dynamisch gehlot werden

	//#######################ADD###################

	titel := ""
	beschreibung := ""
	kategorie := ""
	radio := ""
	if r.Method == "POST" {

		r.ParseForm()
		titel = r.FormValue("titel")
		fmt.Println(titel)
		beschreibung = r.FormValue("beschreibung")
		fmt.Println(beschreibung)
		kategorie = r.FormValue("kategorie")
		fmt.Println(kategorie)
		radio = r.FormValue("answer")
		fmt.Println(radio)

		OberKategorie := ""

		if kategorie == "Biologie" || kategorie == "Chemie" || kategorie == "Elektrotechnik" || kategorie == "Informatik" || kategorie == "Mathematik" || kategorie == "Medizin" || kategorie == "Naturkunde" || kategorie == "Physik" {
			OberKategorie = "Naturwissenschaften"
		}
		if kategorie == "Chinesisch" || kategorie == "Deutsch" || kategorie == "Englisch" || kategorie == "Französisch" || kategorie == "Griechisch" || kategorie == "Italienisch" || kategorie == "Latein" || kategorie == "Russisch" {
			OberKategorie = "Sprachen"
		}
		if kategorie == "Ethik" || kategorie == "Geschichte" || kategorie == "Literatur" || kategorie == "Musik" || kategorie == "Politik" || kategorie == "Recht" || kategorie == "Soziales" || kategorie == "Sport" || kategorie == "Verkehrskunde" {
			OberKategorie = "Gesellschaft"
		}
		if kategorie == "BWL" || kategorie == "Finanzen" || kategorie == "Landwirtschaft" || kategorie == "Marketing" || kategorie == "VWL" {
			OberKategorie = "Wirtschaft"
		}
		if kategorie == "Kriminologie" || kategorie == "Philosophie" || kategorie == "Psychologie" || kategorie == "Pädagogik" || kategorie == "Theologie" {
			OberKategorie = "Geisteswissenschaften"
		}
		if kategorie == "Sonstige" {
			OberKategorie = "Sonstige"
		}

		kk := Karteikasten{}
		kk.ID = GetKarteikastenAnz() + 1
		kk.TYP = "Karteikasten"
		kk.NutzerID = nutzer.ID
		kk.Ersteller = nutzer.Nutzername
		kk.Sichtbarkeit = radio
		kk.Kategorie = OberKategorie
		kk.Unterkategorie = kategorie
		kk.Titel = titel
		kk.Anzahl = 0
		kk.Beschreibung = beschreibung

		Add(kk)
		//FUNKTIONIERT noch nicht
		AddErstellte(nutzer, GetKarteikastenAnz())
	}

	//#############################################

	for _, element := range nutzer.ErstellteKarteien {
		temp_kk := GetKarteikastenByid(element)
		temp_kk.FortschrittP = int(GetKarteikastenFortschritt(temp_kk, GetNutzerById(1)))
		data.MeineKarteikaesten = append(data.MeineKarteikaesten, temp_kk)
	}

	for _, element := range nutzer.GelernteKarteien {
		temp_kk := GetKarteikastenByid(element)
		temp_kk.FortschrittP = int(GetKarteikastenFortschritt(temp_kk, GetNutzerById(1)))
		data.GespeicherteKarteikaesten = append(data.GespeicherteKarteikaesten, temp_kk)
	}

	data.Karteien = strconv.Itoa(GetKarteikastenAnz())
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(nutzer))
	data.Nutzername = nutzer.Nutzername

	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_meine_karteikaesten.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_profil_popup(w http.ResponseWriter, r *http.Request) {
	data := in_menu{MeineKarteien: "", Nutzername: "", Karteien: strconv.Itoa(GetKarteikastenAnz())}

	tempMK := GetNutzerById(1)
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername
	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_profil_popup.html")
	t.ExecuteTemplate(w, "layout", data)
}

func In_profil(w http.ResponseWriter, r *http.Request) {
	data := in_menu{MeineKarteien: "", Nutzername: "", Karteien: strconv.Itoa(GetKarteikastenAnz())}

	tempMK := GetNutzerById(1)
	data.MeineKarteien = strconv.Itoa(GetMeineKarteikaestenAnz(tempMK))
	data.Nutzername = tempMK.Nutzername
	t, _ := template.ParseFiles("./templates/in_menu.html", "./templates/in_profil.html")
	t.ExecuteTemplate(w, "layout", data)
}
