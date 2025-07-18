package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Eb3Messaging struct {
	XMLName     xml.Name `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ Messaging"`
	UserMessage struct {
		MessageInfo struct {
			Timestamp string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ Timestamp"`
			MessageId string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ MessageId"`
		} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ MessageInfo"`
		UserMessage struct {
			From struct {
				PartyId string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ PartyId"`
				Role    string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ Role"`
			} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ From"`
			To struct {
				PartyId string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ PartyId"`
				Role    string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ Role"`
			} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ To"`
		} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ UserMessage"`
		CollaborationInfo struct {
			AgreementRef   string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ AgreementRef"`
			Service        string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ Service"`
			Action         string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ Action"`
			ConversationId string `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ ConversationId"`
		} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ CollaborationInfo"`
		PayloadInfo struct {
			PartInfo []struct {
				Href           string `xml:"href,attr"`
				PartProperties struct {
					Property struct {
						Name  string `xml:"name,attr"`
						Value string `xml:",chardata"`
					} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ Property"`
				} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ PartProperties"`
			} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ PartInfo"`
		} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ PayloadInfo"`
	} `xml:"http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704/ UserMessage"`
}

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

// Ubah struct SecurityInfo agar hanya mengandung KeystoreAlias dan DName
// Parse DName dari certificate PEM hasil keytool

type SecurityInfo struct {
	KeystoreAlias string `json:"keystoreAlias"`
	DName         string `json:"dname"`
}

// Fungsi untuk mengambil DName dari keytool -list -v
func getDNameFromKeystore(alias string) (string, error) {
	partnerKeysPath := "/opt/holodeckb2b/repository/certs/partnerkeys.jks"
	cmd := fmt.Sprintf("keytool -list -v -keystore %s -storepass nosecrets -alias %s", partnerKeysPath, alias)
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("failed to get DName: %v", err)
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Owner:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Owner:")), nil
		}
		if strings.HasPrefix(line, "Subject:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Subject:")), nil
		}
	}
	return "", fmt.Errorf("DName not found in keytool output")
}

// Fungsi untuk mengekstrak informasi keamanan berdasarkan party_id
func extractSecurityInfoFromDB(partyID string) SecurityInfo {
	var securityInfo SecurityInfo
	securityInfo.KeystoreAlias = ""
	securityInfo.DName = ""

	if partyID == "" {
		return securityInfo
	}

	alias, err := getKeystoreAliasFromDB(partyID)
	if err != nil || alias == "" {
		return securityInfo
	}
	securityInfo.KeystoreAlias = alias

	dname, err := getDNameFromKeystore(alias)
	if err == nil && dname != "" {
		securityInfo.DName = dname
	}
	return securityInfo
}

type MessageDetail struct {
	ID           string       `json:"id"`
	Content      string       `json:"content"`
	Sender       string       `json:"sender"`
	Receiver     string       `json:"receiver"`
	Date         string       `json:"date"`
	Subject      string       `json:"subject"`
	FileName     string       `json:"fileName"`
	Attachments  []Attachment `json:"attachments"`
	SecurityInfo SecurityInfo `json:"securityInfo"`
}

// Fungsi untuk mengambil keystore alias dari database berdasarkan party_id
func getKeystoreAliasFromDB(partyID string) (string, error) {
	query := "SELECT keystore_alias FROM party WHERE party_id = ?"
	var alias string
	if err := db.QueryRow(query, partyID).Scan(&alias); err != nil {
		return "", err
	}
	return alias, nil
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

	msgDir := "/opt/holodeckb2b/data/msg_in"
	files, err := os.ReadDir(msgDir)
	if err != nil {
		http.Error(w, "Unable to read message directory", http.StatusInternalServerError)
		log.Printf("Error reading directory: %v", err)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".xml") {
			filePath := filepath.Join(msgDir, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", file.Name(), err)
				continue
			}

			// 1. Coba parse sebagai MessageMetaData (MMD)
			var mmd MessageMetaData
			if err := xml.Unmarshal(content, &mmd); err == nil && mmd.MessageInfo.MessageId == messageID {
				var attachments []Attachment
				var soapFile string

				for _, part := range mmd.PayloadInfo.PartInfo {
					if part.MimeType == "application/xml" && strings.Contains(strings.ToLower(part.Location), "soappart") {
						soapFile = part.Location
						continue
					}
					attachments = append(attachments, Attachment{
						FileName: part.Location,
						MimeType: part.MimeType,
						Url:      "/attachments/" + part.Location,
					})
				}

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

				// Ambil informasi keamanan dari database berdasarkan sender (party_id)
				securityInfo := extractSecurityInfoFromDB(sender)

				messageDetail := MessageDetail{
					ID:           mmd.MessageInfo.MessageId,
					Content:      contentMsg,
					Sender:       sender,
					Receiver:     receiver,
					Date:         date,
					Subject:      subject,
					Attachments:  attachments,
					SecurityInfo: securityInfo,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(messageDetail)
				return
			}

			// 2. Coba parse sebagai Eb3Messaging
			var eb3 Eb3Messaging
			if err := xml.Unmarshal(content, &eb3); err == nil && eb3.UserMessage.MessageInfo.MessageId == messageID {
				var attachments []Attachment
				var contentMsg, sender, receiver, date, subject string

				for _, part := range eb3.UserMessage.PayloadInfo.PartInfo {
					fileName := part.PartProperties.Property.Value
					if strings.HasSuffix(strings.ToLower(fileName), ".xml") && strings.Contains(strings.ToLower(fileName), "soap") {
						soapPath := filepath.Join(msgDir, fileName)
						soapContent, err := os.ReadFile(soapPath)
						if err == nil {
							var envelope SOAPEnvelope
							if err := xml.Unmarshal(soapContent, &envelope); err == nil {
								contentMsg = envelope.Body.MessageContent
								sender = envelope.Header.Messaging.UserMessage.PartyInfo.From.PartyId
								receiver = envelope.Header.Messaging.UserMessage.PartyInfo.To.PartyId
								date = envelope.Header.Messaging.UserMessage.MessageInfo.Timestamp
								subject = envelope.Header.Messaging.UserMessage.CollaborationInfo.Subject
								continue // jangan tambahkan ke attachments
							}
						}
					} else {
						attachments = append(attachments, Attachment{
							FileName: fileName,
							MimeType: "",
							Url:      "/attachments/" + fileName,
						})
					}
				}

				// Ambil informasi keamanan dari database berdasarkan sender (party_id)
				securityInfo := extractSecurityInfoFromDB(sender)

				messageDetail := MessageDetail{
					ID:           eb3.UserMessage.MessageInfo.MessageId,
					Content:      contentMsg,
					Sender:       sender,
					Receiver:     receiver,
					Date:         date,
					Subject:      subject,
					Attachments:  attachments,
					SecurityInfo: securityInfo,
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

	msgDir := "/opt/holodeckb2b/data/msg_in"
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

// Fungsi untuk mendapatkan daftar sertifikat yang tersedia
func listCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type CertificateInfo struct {
		FileName string `json:"fileName"`
		Type     string `json:"type"`
		Size     int64  `json:"size"`
	}

	var certificates []CertificateInfo

	// Cari file sertifikat di folder certs
	certsDir := "/opt/holodeckb2b/repository/certs"
	files, err := os.ReadDir(certsDir)
	if err != nil {
		log.Printf("Error reading certs directory: %v", err)
	} else {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".cert") ||
				strings.HasSuffix(file.Name(), ".pem") ||
				strings.HasSuffix(file.Name(), ".crt") ||
				strings.HasSuffix(file.Name(), ".jks") {
				info, err := file.Info()
				if err != nil {
					continue
				}
				certificates = append(certificates, CertificateInfo{
					FileName: file.Name(),
					Type:     filepath.Ext(file.Name()),
					Size:     info.Size(),
				})
			}
		}
	}

	// Coba list alias dari partnerkeys.jks
	type KeystoreInfo struct {
		KeystoreFile string   `json:"keystoreFile"`
		Aliases      []string `json:"aliases"`
	}

	var keystoreInfo KeystoreInfo
	keystoreInfo.KeystoreFile = "partnerkeys.jks"

	// Gunakan keytool untuk list alias
	cmd := "keytool -list -keystore /opt/holodeckb2b/repository/certs/partnerkeys.jks -storepass nosecrets"
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "alias name:") {
				parts := strings.Split(line, "alias name:")
				if len(parts) > 1 {
					alias := strings.TrimSpace(parts[1])
					keystoreInfo.Aliases = append(keystoreInfo.Aliases, alias)
				}
			}
		}
	}

	response := map[string]interface{}{
		"certificateFiles": certificates,
		"keystoreInfo":     keystoreInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Inisialisasi koneksi database di main() saja, tidak perlu di init()

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(xampp-mariadb:3306)/proyekta")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	http.HandleFunc("/api/mail", viewMessage)
	http.HandleFunc("/api/certificates", listCertificates)
	http.HandleFunc("/download", downloadAttachment)
	http.Handle("/attachments/", http.StripPrefix("/attachments/", http.FileServer(http.Dir("/opt/holodeckb2b/data/msg_in"))))

	log.Println("Starting ViewMessage API server on port 8083...")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
