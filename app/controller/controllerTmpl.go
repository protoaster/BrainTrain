package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	couchdb "github.com/leesper/couchdb-golang"
	"strconv" //strconv.Itoa -> int to string
)

// Structs
type Nutzer struct {
	ID                int
	Nutzername        string
	EMail             string
	Passwort          string
	ErstellteKarteien []int
	GelernteKarteien  []int
}

type alleNutzer struct {
	_id    string
	_rev   string
	Nutzer []Nutzer
}

type Karte struct {
	Num        int
	Titel      string
	Frage      string
	Antwort    string
	NutzerFach string
}

type Fortschritt struct {
	ID           int
	Wiederholung []int
}

type Karteikasten struct {
	ID             int
	_rev           string
	TYP            string
	NutzerID       int
	Ersteller      string
	Sichtbarkeit   string
	Kategorie      string
	Unterkategorie string
	Titel          string
	Anzahl         int
	Beschreibung   string
	Karten         []Karte
	Fortschritt    []Fortschritt
	FortschrittP   int
}

// Gibt die DB zurück (wenn nicht vorhanden = nil)
func GetDB() (d *couchdb.Database) {
	a, b := couchdb.NewDatabase(couchdb.DefaultBaseURL + "/web")
	var err error

	if b != nil {
		fmt.Print(b)
	} else {
		err = a.Available()
		if err == nil {
			//fmt.Println("DB is available")
			return a
		}
	}

	return nil
}

// ############################### START Kartei Methoden ############################### //
func GetKartenAnz() (anz int) {
	kk := GetAlleKarteikaesten()

	anz = 0

	for _, element := range kk {
		anz += len(element.Karten)
	}

	return anz
}

func GetKarteikastenFortschritt(k Karteikasten, nutzer Nutzer) (fortschritt float64) {
	fortschritt = 0
	var zaehler = 0
	xgesamt := len(k.Karten)

	for n := 0; n <= 4; n++ {
		zaehler += n * GetKarteikartenAnzByFach(k, n, nutzer)
	}

	fortschritt = float64(zaehler) / float64(4*float64(xgesamt)) * 100

	return fortschritt
}

func GetKarteikastenWiederholungArr(k Karteikasten, nutzer Nutzer) (i []int) {

	for _, element := range k.Fortschritt {
		if element.ID == nutzer.ID {
			for _, wd := range element.Wiederholung {
				i = append(i, wd)
			}
			return i
		} else {
			return nil
		}
	}
	return nil

}

func GetKarteikartenAnzByFach(k Karteikasten, fach int, n Nutzer) (anz int) {
	var anzahl_fach = 0
	for index, _ := range k.Karten {
		if k.Fortschritt[n.ID-1].Wiederholung[index] == fach {
			anzahl_fach++
		}
	}

	return anzahl_fach
}

// ############################### Ende Kartei Methoden ################################ //

// ############################### START Karteikasten Methoden ############################### //
func GetKarteikastenAnz() (anz int) {
	return len(GetAlleKarteikaesten())
}

func GetAlleKarteikaesten() (kk []Karteikasten) {
	var db *couchdb.Database = GetDB()

	var inmap []map[string]interface{}

	inmap, err := db.QueryJSON(`
	{
	  "selector": {
		"TYP": "Karteikasten"
	  }
	}`)

	for _, element := range inmap {
		var in = mapToJSON(element)

		var temp_kk = Karteikasten{}
		if err == nil {
			json.Unmarshal([]byte(in), &temp_kk)

			kk = append(kk, temp_kk)
		} else {
			fmt.Println(err)
		}
	}

	//for _, element := range kk {
	//	TerminalOutKarteikasten(element)
	//}

	return kk
}

func GetKarteikastenByid(id int) (k Karteikasten) {

	kk := GetAlleKarteikaesten()

	for _, element := range kk {
		if element.ID == id {
			return element
		}
	}

	return k
}

func GetMeineKarteikaestenAnz(n Nutzer) (anz int) {

	var i int
	for range n.ErstellteKarteien {
		i += 1
	}

	for range n.GelernteKarteien {
		i += 1
	}
	fmt.Println(i)
	return i
}

func TerminalOutKarteikasten(k Karteikasten) {
	fmt.Println("############# KARTEIKASTEN ##############")
	fmt.Println("id : " + strconv.Itoa(k.ID))
	fmt.Println("NutzerID : " + strconv.Itoa(k.NutzerID))
	fmt.Println("Oeffentlich : " + k.Sichtbarkeit)
	fmt.Println("Kategorie : " + k.Kategorie)
	fmt.Println("Unterkategorie : " + k.Unterkategorie)
	fmt.Println("Titel : " + k.Titel)
	fmt.Println("Anzahl : " + strconv.Itoa(k.Anzahl))
	fmt.Println("Beschreibung : " + k.Beschreibung)
	fmt.Println("#########################################")
}

// ############################### ENDE Karteikasten Methoden ############################### //

// ############################### START Nutzer Methoden ############################### //

//Wenn nicht vorhanden ID = -1
func GetNutzerById(id int) (n Nutzer) {

	var arr, err = getNutzerArr()

	if err == nil {
		for _, n := range arr {
			if n.ID == id {
				return n
			}
		}
	}

	n = Nutzer{}
	n.ID = -1
	return n
}

//-1 = db not da
//-2 = abfrage nicht möglich
func GetNutzeranz() (anz int) {
	var n, err = getNutzerArr()

	if err == nil {
		return len(n)
	} else {
		fmt.Println(err)
		return -2
	}

	return -1
}

func getNutzerArr() (n []Nutzer, err error) {
	var db *couchdb.Database = GetDB()

	if db == nil {
		return nil, errors.New("Datenbank Verbindung nicht möglich!")
	}

	//Nutzer Wählen
	var result map[string]interface{}

	//result, err = db.Get("nutzer", nil)
	result, err = db.Get("e3c55bfe2a805192b1ab0a0abf03a2d5", nil)

	in := mapToJSON(result)

	an := alleNutzer{}
	json.Unmarshal([]byte(in), &an)

	return an.Nutzer, nil

}

func TerminalOutNutzer(n Nutzer) {
	fmt.Println("ID 		: " + strconv.Itoa(n.ID))
	fmt.Println("Vorname 	: " + n.Nutzername)
	fmt.Println("Email 		: " + n.EMail)
	fmt.Println("Passwort 	: " + n.Passwort)
}

// ###########AddKarteikasten##
func Add(kk Karteikasten) error {
	// Convert Todo suct to map[string]interface as required by Save() method
	KarteiK := kk2Map(kk)

	// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
	delete(KarteiK, "_id")
	delete(KarteiK, "_rev")
	delete(KarteiK, "FortschrittP")

	var db *couchdb.Database = GetDB()
	// Add todo to DB
	_, _, err := db.Save(KarteiK, nil)

	if err != nil {
		fmt.Printf("[Add] error: %s", err)
	}

	return err
}

//FUNKTIONIERT NOCH NICHT

func AddErstellte(n Nutzer, neu int) error {
	var db *couchdb.Database = GetDB()

	an := GetNutzerDatei()
	var err error

	fmt.Println(neu)

	for i := 0; i <= len(an[0].Nutzer); i++ {

		if n.ID == i {

			fmt.Println(i)
			fmt.Println("komm ich hier hin")
			ek := an[0].Nutzer[i-1].ErstellteKarteien
			fmt.Println(ek)
			t := append(ek, neu)
			fmt.Println(t)
			an[0].Nutzer[i-1].ErstellteKarteien = t

			err = db.Set("e3c55bfe2a805192b1ab0a0abf03a2d5", nutzer2Map(an[0]))
			fmt.Println(err)
		}

	}

	return err

}

//FUNKTIONIERT NOCH NICHT
func GetNutzerDatei() (an []alleNutzer) {
	var db *couchdb.Database = GetDB()

	var inmap []map[string]interface{}

	inmap, err := db.QueryJSON(`
	{
		"selector": {
		"TYP": "nutzer"
		}
	}`)

	for _, element := range inmap {
		var in = mapToJSON(element)

		var temp_an = alleNutzer{}
		if err == nil {
			json.Unmarshal([]byte(in), &temp_an)

			an = append(an, temp_an)
		} else {
			fmt.Println(err)
		}
	}

	//for _, element := range kk {
	//	TerminalOutKarteikasten(element)
	//}

	return an
}

// ############################### ENDE Nutzer Methoden ############################### //

func mapToJSON(inMap map[string]interface{}) (s string) {
	var b []byte

	b, err := json.Marshal(inMap)
	jsonString := string(b)

	//Error Output
	if err != nil {
		fmt.Print("JSON Convertion Error: ")
		fmt.Println(err)
	}

	return jsonString
}

func kk2Map(kk Karteikasten) map[string]interface{} {
	var doc map[string]interface{}
	tJSON, _ := json.Marshal(kk)
	json.Unmarshal(tJSON, &doc)

	return doc
}

func nutzer2Map(an alleNutzer) map[string]interface{} {
	var doc map[string]interface{}
	tJSON, _ := json.Marshal(an)
	json.Unmarshal(tJSON, &doc)

	return doc
}
