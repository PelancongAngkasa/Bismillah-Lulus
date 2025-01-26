package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Struktur pesan yang diterima melalui API
type AS4Message struct {
	FromParty string `json:"fromParty"`
	ToParty   string `json:"toParty"`
	Service   string `json:"service"`
	Action    string `json:"action"`
	MessageID string `json:"messageId"`
	Payload   string `json:"payload"`
}

// Struktur XML untuk metadata pesan (MMD)
type MessageMetaData struct {
	XMLName           xml.Name `xml:"MessageMetaData"`
	XMLNS             string   `xml:"xmlns,attr"`
	SchemaLocation    string   `xml:"xsi:schemaLocation,attr"`
	XSI               string   `xml:"xmlns:xsi,attr"`
	CollaborationInfo struct {
		AgreementRef struct {
			PMode string `xml:"pmode,attr"`
		} `xml:"AgreementRef"`
		ConversationId string `xml:"ConversationId"`
	} `xml:"CollaborationInfo"`
	PayloadInfo struct {
		DeleteFilesAfterSubmit bool       `xml:"deleteFilesAfterSubmit,attr"`
		PartInfo               []struct { // Ubah menjadi array
			Containment string `xml:"containment,attr"`
			MimeType    string `xml:"mimeType,attr"`
			Location    string `xml:"location,attr"`
		} `xml:"PartInfo"`
	} `xml:"PayloadInfo"`
}

// saveFile saves an uploaded file to the specified directory
func saveFile(file multipart.File, header *multipart.FileHeader, destDir string) (string, error) {
	destPath := filepath.Join(destDir, header.Filename)

	// Detect MIME type
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	file.Seek(0, io.SeekStart) // Reset file reader position
	mimeType := http.DetectContentType(buffer)

	// Validate MIME type
	if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "application/pdf" {
		return "", fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// Save file
	outFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return header.Filename, nil
}

func getAddressFromDB(partyName string) (string, string, error) {
	query := "SELECT endpoint_url, party_id FROM party WHERE name = ?"
	var address, partyID string

	err := db.QueryRow(query, partyName).Scan(&address, &partyID)
	if err != nil {
		return "", "", err
	}

	return address, partyID, nil
}

// Fungsi untuk menggantikan placeholder dalam file template
func replacePlaceholders(template, address, partyID string) (string, error) {
	if template == "" {
		return "", fmt.Errorf("template is empty")
	}
	replaced := strings.ReplaceAll(template, "${dynamic_responder_party_id}", partyID)
	replaced = strings.ReplaceAll(replaced, "${dynamic_address}", address)
	return replaced, nil
}

// Fungsi untuk memperbarui P-Mode dengan nilai dynamicAddress dan partyID
func updatePModeTemplate(partyName string) error {

	// Ambil dynamicAddress dan partyID dari database
	dynamicAddress, partyID, err := getAddressFromDB(partyName)
	if err != nil {
		return fmt.Errorf("failed to get address and partyID for party %s: %v", partyName, err)
	}
	log.Printf("Dynamic Address: %s, Party ID: %s", dynamicAddress, partyID)

	// Path file template P-Mode (hardcoded untuk mode push)
	templateFile := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-A\examples\pmodes\pm-push.xml`

	// Baca template P-Mode
	pmodeContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read P-Mode template: %v", err)
	}

	// Gantikan placeholder dengan dynamicAddress dan partyID
	updatedContent, err := replacePlaceholders(string(pmodeContent), dynamicAddress, partyID)
	if err != nil {
		return fmt.Errorf("failed to replace placeholders in template: %v", err)
	}

	// Path untuk file P-Mode aktif
	activePMode := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-A\repository\pmodes\current-pmode.xml`

	// Overwrite P-Mode aktif dengan konten yang diperbarui
	if err := os.WriteFile(activePMode, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to overwrite active P-Mode file: %v", err)
	}

	log.Printf("P-Mode successfully updated for party: %s", partyName)
	return nil
}

func writePayloadAsSOAP(message AS4Message, outputDir string) error {
	payloadFileName := fmt.Sprintf("%s_payload.xml", message.MessageID)
	payloadPath := filepath.Join(outputDir, payloadFileName)

	// Struktur SOAP Envelope
	soapEnvelope := struct {
		XMLName xml.Name `xml:"SOAP:Envelope"`
		SOAP    string   `xml:"xmlns:SOAP,attr"`
		EB      string   `xml:"xmlns:eb,attr"`
		Header  struct {
			Messaging struct {
				UserMessage struct {
					MessageInfo struct {
						MessageId string `xml:"eb:MessageId"`
						Timestamp string `xml:"eb:Timestamp"`
					} `xml:"eb:MessageInfo"`
					PartyInfo struct {
						From struct {
							PartyId string `xml:"eb:PartyId"`
						} `xml:"eb:From"`
						To struct {
							PartyId string `xml:"eb:PartyId"`
						} `xml:"eb:To"`
					} `xml:"eb:PartyInfo"`
					CollaborationInfo struct {
						Service string `xml:"eb:Service"`
						Action  string `xml:"eb:Action"`
					} `xml:"eb:CollaborationInfo"`
					PayloadInfo struct {
						PartInfo struct {
							Href         string `xml:"eb:href,attr"`
							PartProperty struct {
								Name  string `xml:"name,attr"`
								Value string `xml:",chardata"`
							} `xml:"eb:PartProperties>eb:Property"`
						} `xml:"eb:PartInfo"`
					} `xml:"eb:PayloadInfo"`
				} `xml:"eb:UserMessage"`
			} `xml:"eb:Messaging"`
		} `xml:"SOAP:Header"`
		Body struct {
			MessageContent struct {
				Document string `xml:"myns:Document"`
			} `xml:"myns:MessageContent"`
		} `xml:"SOAP:Body"`
	}{
		SOAP: "http://www.w3.org/2003/05/soap-envelope",
		EB:   "http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704",
	}

	// Isi Header SOAP
	soapEnvelope.Header.Messaging.UserMessage.MessageInfo.MessageId = message.MessageID
	soapEnvelope.Header.Messaging.UserMessage.MessageInfo.Timestamp = time.Now().UTC().Format(time.RFC3339)
	soapEnvelope.Header.Messaging.UserMessage.PartyInfo.From.PartyId = message.FromParty
	soapEnvelope.Header.Messaging.UserMessage.PartyInfo.To.PartyId = message.ToParty
	soapEnvelope.Header.Messaging.UserMessage.CollaborationInfo.Service = message.Service
	soapEnvelope.Header.Messaging.UserMessage.CollaborationInfo.Action = message.Action
	soapEnvelope.Header.Messaging.UserMessage.PayloadInfo.PartInfo.Href = "cid:payload1"
	soapEnvelope.Header.Messaging.UserMessage.PayloadInfo.PartInfo.PartProperty.Name = "MimeType"
	soapEnvelope.Header.Messaging.UserMessage.PayloadInfo.PartInfo.PartProperty.Value = "application/xml"

	// Isi Body SOAP
	soapEnvelope.Body.MessageContent.Document = message.Payload

	// Serialize ke file
	soapXML, err := xml.MarshalIndent(soapEnvelope, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SOAP payload: %v", err)
	}
	return os.WriteFile(payloadPath, soapXML, 0644)
}

// Fungsi untuk membuat file MMD
func writeMMDFile(message AS4Message, attachmentFileName string, mimeType string, outputDir string) error {
	// Direktori untuk payload
	payloadDir := filepath.Join(outputDir, "payloads")
	if _, err := os.Stat(payloadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(payloadDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create payload directory: %v", err)
		}
	}

	// Buat file SOAP envelope
	err := writePayloadAsSOAP(message, payloadDir)
	if err != nil {
		return fmt.Errorf("failed to write SOAP payload: %v", err)
	}

	soapPayloadFileName := fmt.Sprintf("%s_payload.xml", message.MessageID)
	soapPayloadPath := filepath.Join("payloads", soapPayloadFileName)

	// Metadata MMD
	mmd := MessageMetaData{
		XMLNS:          "http://holodeck-b2b.org/schemas/2014/06/mmd",
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://holodeck-b2b.org/schemas/2014/06/mmd ../repository/xsd/messagemetadata.xsd",
	}
	mmd.CollaborationInfo.AgreementRef.PMode = "current-pmode-push"
	mmd.CollaborationInfo.ConversationId = "org:holodeckb2b:test:conversation"
	mmd.PayloadInfo.DeleteFilesAfterSubmit = false

	// Tambahkan PartInfo untuk SOAP payload
	mmd.PayloadInfo.PartInfo = append(mmd.PayloadInfo.PartInfo, struct {
		Containment string `xml:"containment,attr"`
		MimeType    string `xml:"mimeType,attr"`
		Location    string `xml:"location,attr"`
	}{
		Containment: "inline",
		MimeType:    "application/xml",
		Location:    soapPayloadPath,
	})

	// Jika ada attachment, tambahkan ke PartInfo
	if attachmentFileName != "" {
		mmd.PayloadInfo.PartInfo = append(mmd.PayloadInfo.PartInfo, struct {
			Containment string `xml:"containment,attr"`
			MimeType    string `xml:"mimeType,attr"`
			Location    string `xml:"location,attr"`
		}{
			Containment: "attachment",
			MimeType:    mimeType,
			Location:    filepath.Join("payloads", attachmentFileName),
		})
	}

	// Serialize metadata MMD ke file
	mmdFileName := fmt.Sprintf("%s.mmd", message.MessageID)
	mmdFilePath := filepath.Join(outputDir, mmdFileName)

	mmdXML, err := xml.MarshalIndent(mmd, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal MMD to XML: %v", err)
	}

	if err := os.WriteFile(mmdFilePath, mmdXML, 0644); err != nil {
		return fmt.Errorf("failed to write MMD file: %v", err)
	}

	return nil
}

// Handler untuk menerima pesan AS4 dan memproses lampiran
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	// Preflight CORS handling
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// Ambil data dari request
	message := AS4Message{
		FromParty: r.FormValue("fromParty"),
		ToParty:   r.FormValue("toParty"),
		Service:   r.FormValue("service"),
		Action:    r.FormValue("action"),
		MessageID: r.FormValue("messageId"),
		Payload:   r.FormValue("payload"),
	}

	// Update P-Mode dengan dynamicAddress dan partyID
	if err := updatePModeTemplate(message.ToParty); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update P-Mode template: %v", err), http.StatusInternalServerError)
		return
	}

	payloadDir := "C:/Users/Yusuf/Documents/Kuliah/RPLK/Tugas Akhir/holodeckb2b-7.0.0-A/data/msg_out/payloads"
	if _, err := os.Stat(payloadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(payloadDir, os.ModePerm); err != nil {
			http.Error(w, "Failed to create payload directory", http.StatusInternalServerError)
			return
		}
	}

	attachmentFileName := ""
	mimeType := ""

	// Tangani file attachment jika ada
	if attachmentFile, fileHeader, err := r.FormFile("attachment"); err == nil {
		defer attachmentFile.Close()
		attachmentFileName, err = saveFile(attachmentFile, fileHeader, payloadDir)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
			return
		}
		mimeType = fileHeader.Header.Get("Content-Type")
		log.Printf("Attachment saved: %s (MIME: %s)", attachmentFileName, mimeType)
	} else {
		log.Println("No attachment provided; using only payload.")
	}

	// Tulis file SOAP payload
	if err := writePayloadAsSOAP(message, payloadDir); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write SOAP payload: %v", err), http.StatusInternalServerError)
		return
	}

	// Tulis file MMD (metadata)
	outputDir := "C:/Users/Yusuf/Documents/Kuliah/RPLK/Tugas Akhir/holodeckb2b-7.0.0-A/data/msg_out"
	if err := writeMMDFile(message, attachmentFileName, mimeType, outputDir); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create MMD file: %v", err), http.StatusInternalServerError)
		return
	}

	// Respon sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message processed successfully"})
}

// Fungsi utama
func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/proyekta")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	http.HandleFunc("/api/as4/send", MessageHandler)
	log.Println("Starting server on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
