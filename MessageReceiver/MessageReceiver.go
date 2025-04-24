package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  struct {
		Messaging struct {
			UserMessage struct {
				MessageInfo struct {
					MessageId string `xml:"MessageId"`
					Timestamp string `xml:"Timestamp"`
				} `xml:"MessageInfo"`
				PartyInfo struct {
					From struct {
						PartyId string `xml:"PartyId"`
					} `xml:"From"`
					To struct {
						PartyId string `xml:"PartyId"`
					} `xml:"To"`
				} `xml:"PartyInfo"`
				CollaborationInfo struct {
					Service string `xml:"Service"`
					Action  string `xml:"Action"`
					Subject string `xml:"Subject"`
				} `xml:"CollaborationInfo"`
				PayloadInfo struct {
					PartInfo struct {
						Href           string `xml:"href,attr"`
						PartProperties struct {
							Property struct {
								Name  string `xml:"name,attr"`
								Value string `xml:",chardata"`
							} `xml:"Property"`
						} `xml:"PartProperties"`
					} `xml:"PartInfo"`
				} `xml:"PayloadInfo"`
			} `xml:"UserMessage"`
		} `xml:"Messaging"`
	} `xml:"Header"`
}

type Mail struct {
	ID       string `json:"id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Date     string `json:"date"`
	Subject  string `json:"subject"`
	FileName string `json:"fileName"`
}

func getMails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	msgDir := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-B\data\msg_in`
	files, err := ioutil.ReadDir(msgDir)
	if err != nil {
		http.Error(w, "Unable to read message directory", http.StatusInternalServerError)
		return
	}

	var mails []Mail

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".xml" {
			filePath := filepath.Join(msgDir, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", file.Name(), err)
				continue
			}

			var envelope SOAPEnvelope
			if err := xml.Unmarshal(content, &envelope); err != nil {
				log.Printf("Error unmarshalling XML %s: %v", file.Name(), err)
				continue
			}

			// Validasi data
			if envelope.Header.Messaging.UserMessage.MessageInfo.MessageId == "" ||
				envelope.Header.Messaging.UserMessage.PartyInfo.From.PartyId == "" ||
				envelope.Header.Messaging.UserMessage.PartyInfo.To.PartyId == "" ||
				envelope.Header.Messaging.UserMessage.MessageInfo.Timestamp == "" {
				log.Printf("Invalid structure in file %s", file.Name())
				continue
			}

			mails = append(mails, Mail{
				ID:       envelope.Header.Messaging.UserMessage.MessageInfo.MessageId,
				Sender:   envelope.Header.Messaging.UserMessage.PartyInfo.From.PartyId,
				Receiver: envelope.Header.Messaging.UserMessage.PartyInfo.To.PartyId,
				Date:     envelope.Header.Messaging.UserMessage.MessageInfo.Timestamp,
				Subject:  envelope.Header.Messaging.UserMessage.CollaborationInfo.Subject,
				FileName: file.Name(),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mails); err != nil {
		http.Error(w, "Failed to encode mails", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/api/mails", getMails)

	log.Println("Starting AS4 API server on port 9091...")
	if err := http.ListenAndServe(":9091", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
