// tessera project main.go

package main
	
import (
	"os"
	"log"
	"io"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
)

//Templates erstellen und auf Fehler überprüfen 
var t = template.Must(template.ParseFiles("temp/index.html"))
var t1 = template.Must(template.ParseFiles("temp/upload.html"))
var t3 = template.Must(template.ParseFiles("temp/main.html"))
var terr = template.Must(template.ParseFiles("temp/error.html"))

func main() {
	//Verbindug mit der DB 
	
	// "lokale" function:
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	//Nur in der FH 
	//dbSession, err := mgo.Dial("mongodb://borsti.inf.fh-flensburg.de:27017")
	//Zuhause 
	dbSession, err := mgo.Dial("localhost")
	check(err)
	defer dbSession.Close()
	
	db := dbSession.DB("JTwiessel_MatrNr570705_MongoDB")
	
	//Collection stellt die Tabelle der DB da
	coll := db.C("JTessera")
	
	//DB Inhalt?
	type user struct{
		Id bson.ObjectId  `bson:"_id"`
		Username string
		Password string 
	}
	
	var doc3 = user {bson.NewObjectId(), "Foxi", "44"}
	
	err = coll.Insert(doc3)	// kann eine beliebig lange Liste enthalten 
	check(err)
	
	http.HandleFunc("/regist", response)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/main", htmain)
	http.Handle("/", http.FileServer(http.Dir("./")))
	
	
	http.ListenAndServe(":4242", nil)
	
	// Datenbank löschen:
	err = db.DropDatabase()
	check(err)
}
func htmain(w http.ResponseWriter, r *http.Request){
	switch r.Method{
		case "GET":
			t3.ExecuteTemplate(w, "main.html", nil)
		default: 
			terr.ExecuteTemplate(w, "error.html", nil)		
	}
}
func response(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	
	switch r.Method{
		case "GET": 
			t.ExecuteTemplate(w ,"index.html", nil)	
		case "POST":
			t.ExecuteTemplate(w, "index.html", r.FormValue("name"))
	}
}

func upload(w http.ResponseWriter, r *http.Request){
	
	//Aus der Vorlesungsfolie 
	
	switch r.Method {
	case "GET": // erster Aufruf ist GET -> http://localhost:4242/upload
		t1.ExecuteTemplate(w, "upload.html", nil) // upload-Form senden

	case "POST": // Daten der multipart-form empfangen, files speichern
		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Jeden "part" in den Ordner ./upDownL kopieren:
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			// falls part.FileName() leer, überspringen:
			if part.FileName() == "" {
				continue
			}

			dst, err := os.Create("./upDownL/" + part.FileName())
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Erfolgsmeldung anzeigen:
		t1.ExecuteTemplate(w, "upload.html", "Upload erfolgreich.")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}