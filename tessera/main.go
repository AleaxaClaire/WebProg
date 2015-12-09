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
var t = template.Must(template.ParseFiles("docu/tessera.html"))
var t1 = template.Must(template.ParseFiles("docu/upload.html"))
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
	
	err = coll.Insert(doc3)	// beliebig lange Liste
	check(err)
	
	//http.Handle("/static", http.FileServer(http.Dir("/docu/")))
		
	//http.Handle("/", http.FileServer(http.Dir("./docu/tessera.css")))
	http.HandleFunc("/regist", response)
	http.HandleFunc("/upload", upload)
	
	http.ListenAndServe(":4242", nil)
	
	// Datenbank löschen:
	err = db.DropDatabase()
	check(err)
}

func response(w http.ResponseWriter, r *http.Request){
	switch r.Method{
		case "GET": 
			t.ExecuteTemplate(w ,"tessera.html", nil)	
		case "POST":
			t.ExecuteTemplate(w, "tessera.html",  "Erfolg")
	}
}

func upload(w http.ResponseWriter, r *http.Request){
	
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