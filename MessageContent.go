package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"path/filepath"
)

type AS4Message struct {
	XMLName     xml.Name `xml:"AS4Message"`
	ID          string   `xml:"ID"`
	Payload     string   `xml:"Payload"`
	Sender      string   `xml:"Sender"`
	Recipient   string   `xml:"Recipient"`
	ContentType string   `xml:"ContentType"`
}

func handleAS4Message(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read and parse the incoming XML body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var message AS4Message
	if err := xml.Unmarshal(body, &message); err != nil {
		http.Error(w, "Invalid XML payload", http.StatusBadRequest)
		return
	}

	// Log the received message
	log.Printf("Received AS4 Message: %+v\n", message)

	// Save the message to a file in the data/msg_in folder
	msgFilePath := filepath.Join("data", "msg_in", fmt.Sprintf("%s.xml", message.ID))
	if err := ioutil.WriteFile(msgFilePath, body, 0644); err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "AS4 Message received and saved successfully")
}

func displayMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read messages and images from the data/msg_in folder
	msgDir := filepath.Join("data", "msg_in")
	files, err := ioutil.ReadDir(msgDir)
	if err != nil {
		http.Error(w, "Unable to read message directory", http.StatusInternalServerError)
		return
	}

	// Build the HTML response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<html><body>")
	fmt.Fprintln(w, "<h1>Received Messages</h1>")
	fmt.Fprintln(w, "<ul>")

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".xml" {
			msgPath := filepath.Join(msgDir, file.Name())
			content, err := ioutil.ReadFile(msgPath)
			if err != nil {
				log.Printf("Error reading message file: %v", err)
				continue
			}

			fmt.Fprintf(w, "<li><pre>%s</pre></li>", content)
		} else if filepath.Ext(file.Name()) == ".jpg" || filepath.Ext(file.Name()) == ".png" {
			imgPath := filepath.Join("data", "msg_in", file.Name())
			fmt.Fprintf(w, "<li><img src='/%s' alt='%s' style='max-width:300px;'/></li>", imgPath, file.Name())
		}
	}

	fmt.Fprintln(w, "</ul>")
	fmt.Fprintln(w, "</body></html>")
}

/*func main() {
	http.HandleFunc("/api/as4/receive", handleAS4Message)
	http.HandleFunc("/", displayMessages)

	// Serve static files (images) from the data/msg_in folder
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))

	log.Println("Starting AS4 API server on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}*/
