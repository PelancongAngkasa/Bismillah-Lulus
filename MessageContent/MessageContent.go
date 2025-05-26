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

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Header  struct {
		Messaging struct {
			UserMessage struct {
				MessageInfo struct {
					MessageId string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 MessageId"`
					Timestamp string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 Timestamp"`
				} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 MessageInfo"`
				PartyInfo struct {
					From struct {
						PartyId string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 PartyId"`
					} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 From"`
					To struct {
						PartyId string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 PartyId"`
					} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 To"`
				} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 PartyInfo"`
				CollaborationInfo struct {
					Service string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 Service"`
					Action  string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 Action"`
					Subject string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 Subject"`
				} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 CollaborationInfo"`
			} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 UserMessage"`
		} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704 Messaging"`
	} `xml:"http://www.w3.org/2003/05/soap-envelope Header"`
	Body struct {
		MessageContent string `xml:"http://example.org/myns MessageContent"`
	} `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
}

// Struct untuk membaca file .mmd.xml
type MessageMetaData struct {
	XMLName        xml.Name `xml:"MessageMetaData"`
	XMLNS          string   `xml:"xmlns,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	XSI            string   `xml:"xmlns:xsi,attr"`
	MessageInfo    struct {
		Timestamp string `xml:"Timestamp"`
		MessageId string `xml:"MessageId"`
	} `xml:"MessageInfo"`
	CollaborationInfo struct {
		AgreementRef struct {
			PMode string `xml:"pmode,attr"`
		} `xml:"AgreementRef"`
		Service        string `xml:"Service"`
		Action         string `xml:"Action"`
		ConversationId string `xml:"ConversationId"`
	} `xml:"CollaborationInfo"`
	PayloadInfo struct {
		DeleteFilesAfterSubmit bool `xml:"deleteFilesAfterSubmit,attr"`
		PartInfo               []struct {
			URI         string `xml:"uri,attr"`
			Containment string `xml:"containment,attr"`
			MimeType    string `xml:"mimeType,attr"`
			Location    string `xml:"location,attr"`
		} `xml:"PartInfo"`
	} `xml:"PayloadInfo"`
}

type Attachment struct {
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	Url      string `json:"url"`
}

type MessageDetail struct {
	ID          string       `json:"id"`
	Content     string       `json:"content"`
	Sender      string       `json:"sender"`
	Receiver    string       `json:"receiver"`
	Date        string       `json:"date"`
	Subject     string       `json:"subject"`
	FileName    string       `json:"fileName"`
	Attachments []Attachment `json:"attachments"`
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

	msgDir := "C:/Users/Yusuf/Documents/Kuliah/RPLK/Tugas Akhir/holodeckb2b-7.0.0-B/data/msg_in"
	files, err := os.ReadDir(msgDir)
	if err != nil {
		http.Error(w, "Unable to read message directory", http.StatusInternalServerError)
		log.Printf("Error reading directory: %v", err)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mmd.xml") {
			filePath := filepath.Join(msgDir, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", file.Name(), err)
				continue
			}

			var mmd MessageMetaData
			if err := xml.Unmarshal(content, &mmd); err != nil {
				log.Printf("Error unmarshalling XML %s: %v", file.Name(), err)
				continue
			}

			if mmd.MessageInfo.MessageId == messageID {
				// Ambil daftar attachment
				var attachments []Attachment
				var soapFile string
				for _, part := range mmd.PayloadInfo.PartInfo {
					// Simpan file SOAP untuk parsing, tapi jangan tampilkan sebagai attachment
					if part.MimeType == "application/xml" && strings.Contains(strings.ToLower(part.Location), "soappart") {
						soapFile = part.Location
						continue // skip dari list attachment
					}

					attachments = append(attachments, Attachment{
						FileName: part.Location,
						MimeType: part.MimeType,
						Url:      "/attachments/" + part.Location,
					})
				}

				// Default value jika SOAP tidak ditemukan
				var envelope SOAPEnvelope
				var contentMsg, sender, receiver, date, subject string

				if soapFile != "" {
					soapPath := filepath.Join(msgDir, soapFile)
					soapContent, err := os.ReadFile(soapPath)
					if err == nil {
						if err := xml.Unmarshal(soapContent, &envelope); err == nil {
							contentMsg = envelope.Body.MessageContent
							sender = envelope.Header.Messaging.UserMessage.PartyInfo.From.PartyId
							receiver = envelope.Header.Messaging.UserMessage.PartyInfo.To.PartyId
							date = envelope.Header.Messaging.UserMessage.MessageInfo.Timestamp
							subject = envelope.Header.Messaging.UserMessage.CollaborationInfo.Subject
						}
					}
				}

				messageDetail := MessageDetail{
					ID:          mmd.MessageInfo.MessageId,
					Content:     contentMsg,
					Sender:      sender,
					Receiver:    receiver,
					Date:        date,
					Subject:     subject,
					Attachments: attachments,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(messageDetail)
				return
			}
		}
	}

	http.Error(w, "Message not found", http.StatusNotFound)
}

func downloadAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Nama file diambil dari query string (?name=...)
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	msgDir := "C:/Users/Yusuf/Documents/Kuliah/RPLK/Tugas Akhir/holodeckb2b-7.0.0-B/data/msg_in"
	filePath := filepath.Join(msgDir, fileName)

	// Cek apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set agar browser langsung download
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	http.ServeFile(w, r, filePath)
}

func main() {
	http.HandleFunc("/api/mail", viewMessage)
	http.HandleFunc("/download", downloadAttachment)
	http.Handle("/attachments/", http.StripPrefix("/attachments/", http.FileServer(http.Dir("C:/Users/Yusuf/Documents/Kuliah/RPLK/Tugas Akhir/holodeckb2b-7.0.0-B/data/msg_in"))))

	log.Println("Starting ViewMessage API server on port 9093...")
	if err := http.ListenAndServe(":9093", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
