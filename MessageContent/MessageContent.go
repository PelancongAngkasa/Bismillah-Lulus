package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Tambahkan field Body ke struct SOAPEnvelope
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
	Body struct {
		MessageContent struct {
			Document string `xml:"Document"`
		} `xml:"MessageContent"`
	} `xml:"Body"`
}

type MessageDetail struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	Sender     string `json:"sender"`
	Receiver   string `json:"receiver"`
	Date       string `json:"date"`
	Subject    string `json:"subject"`
	FileName   string `json:"fileName"`
	Attachment string `json:"attachment,omitempty"`
}

func viewMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	messageID := r.URL.Query().Get("id")
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	msgDir := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-B\data\msg_in`
	files, err := os.ReadDir(msgDir)
	if err != nil {
		http.Error(w, "Unable to read message directory", http.StatusInternalServerError)
		log.Printf("Error reading directory: %v", err)
		return
	}

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

			if envelope.Header.Messaging.UserMessage.MessageInfo.MessageId == messageID {
				messageDetail := MessageDetail{
					ID:       envelope.Header.Messaging.UserMessage.MessageInfo.MessageId,
					Content:  envelope.Body.MessageContent.Document,
					Sender:   envelope.Header.Messaging.UserMessage.PartyInfo.From.PartyId,
					Receiver: envelope.Header.Messaging.UserMessage.PartyInfo.To.PartyId,
					Date:     envelope.Header.Messaging.UserMessage.MessageInfo.Timestamp,
					Subject:  envelope.Header.Messaging.UserMessage.CollaborationInfo.Subject,
				}

				attachmentPath := filepath.Join(msgDir, strings.TrimSuffix(file.Name(), ".xml")+".pdf")
				if _, err := os.Stat(attachmentPath); err == nil {
					messageDetail.Attachment = "/attachments/" + filepath.Base(attachmentPath)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(messageDetail)
				return
			}
		}
	}

	http.Error(w, "Message not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/api/mail", viewMessage)
	http.Handle("/attachments/", http.StripPrefix("/attachments/", http.FileServer(http.Dir(`C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-B\data\attachments`))))

	log.Println("Starting ViewMessage API server on port 9092...")
	if err := http.ListenAndServe(":9092", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
