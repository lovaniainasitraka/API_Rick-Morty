package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

// Structure pour représenter un personnage Rick & Morty
type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Gender  string `json:"gender"`
	URL     string `json:"url"`
	Origin  struct {
		Name string `json:"name"`
	} `json:"origin"`
	Image string `json:"image"`
}

// Structure pour la réponse de l'API
type APIResponse struct {
	Results []Character `json:"results"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// URL de l'API Rick & Morty (20 premiers personnages)
	url := "https://rickandmortyapi.com/api/character?page=1"

	// Effectuer une requête GET
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erreur lors de la requête : %v", err)
	}
	defer resp.Body.Close()

	// Lire le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du corps : %v", err)
	}

	// Décoder le JSON
	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		log.Fatalf("Erreur lors du décodage JSON : %v", err)
	}

	// Charger le template HTML
	tmpl := template.Must(template.ParseFiles("index.html"))
	err = tmpl.Execute(w, apiResp.Results[:20]) // limiter à 20 personnages
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution du template : %v", err)
	}
}

func main() {
	// Servir les fichiers statiques (CSS, vidéo, fonts)
	http.Handle("/style.css", http.FileServer(http.Dir(".")))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/fonts/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/", handler)

	fmt.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
