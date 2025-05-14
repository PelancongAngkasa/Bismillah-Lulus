package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PartInfo struct {
	Href           string `xml:"href,attr"`
	PartProperties struct {
		Properties []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:",chardata"`
		} `xml:"Property"`
	} `xml:"PartProperties"`
}

// Tambahkan field xmlns agar bisa parsing tag dengan prefix (misalnya myns)
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
					PartInfos []PartInfo `xml:"PartInfo"`
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
	ID          string   `json:"id"`
	Content     string   `json:"content"`
	Sender      string   `json:"sender"`
	Receiver    string   `json:"receiver"`
	Date        string   `json:"date"`
	Subject     string   `json:"subject"`
	FileNames   []string `json:"fileNames"`
	Attachments []string `json:"attachments"`
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

			currentMsgID := envelope.Header.Messaging.UserMessage.MessageInfo.MessageId
			if currentMsgID == messageID {
				var fileNames []string
				var attachments []string

				for _, part := range envelope.Header.Messaging.UserMessage.PayloadInfo.PartInfos {
					var originalFileName string
					for _, prop := range part.PartProperties.Properties {
						if prop.Name == "OriginalFileName" {
							originalFileName = prop.Value
							break
						}
					}

					if originalFileName != "" {
						fileNames = append(fileNames, originalFileName)

						attachmentPath := filepath.Join(msgDir, originalFileName)
						if _, err := os.Stat(attachmentPath); err == nil {
							attachments = append(attachments, "/attachments/"+originalFileName)
						} else {
							log.Printf("Attachment file %s not found", originalFileName)
						}
					}
				}

				messageDetail := MessageDetail{
					ID:          currentMsgID,
					Content:     envelope.Body.MessageContent.Document,
					Sender:      envelope.Header.Messaging.UserMessage.PartyInfo.From.PartyId,
					Receiver:    envelope.Header.Messaging.UserMessage.PartyInfo.To.PartyId,
					Date:        envelope.Header.Messaging.UserMessage.MessageInfo.Timestamp,
					Subject:     envelope.Header.Messaging.UserMessage.CollaborationInfo.Subject,
					FileNames:   fileNames,
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

func downloadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get file name from query parameter
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Set the path to your attachments directory
	attachmentDir := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-B\data\msg_in`
	filePath := filepath.Join(attachmentDir, fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the content disposition header to force download
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	// Set content type based on file extension
	if strings.HasSuffix(fileName, ".pdf") {
		w.Header().Set("Content-Type", "application/pdf")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Copy the file to response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error sending file", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/api/mail", viewMessage)
	http.HandleFunc("/download", downloadFile)                                                                                                                                     // Tambahkan route baru untuk download
	http.Handle("/attachments/", http.StripPrefix("/attachments/", http.FileServer(http.Dir(`C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-B\data\msg_in`)))) // Serve attachments

	log.Println("Starting ViewMessage API server on port 9092...")
	if err := http.ListenAndServe(":9092", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
