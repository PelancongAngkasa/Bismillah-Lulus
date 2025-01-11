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
		DeleteFilesAfterSubmit bool `xml:"deleteFilesAfterSubmit,attr"`
		PartInfo               struct {
			Containment string `xml:"containment,attr"`
			MimeType    string `xml:"mimeType,attr"`
			Location    string `xml:"location,attr"`
		} `xml:"PartInfo"`
	} `xml:"PayloadInfo"`
}

// Fungsi untuk menyimpan file yang diterima ke folder tertentu
func saveFile(file multipart.File, header *multipart.FileHeader, destDir string) (string, error) {
	destPath := filepath.Join(destDir, header.Filename)
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

// Fungsi untuk mengambil address dinamis dari database berdasarkan nama party
func getAddressFromDB(partyName string) (string, error) {
	query := "SELECT endpoint_url FROM party WHERE name = ?"
	var address string
	err := db.QueryRow(query, partyName).Scan(&address)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no address found for party: %s", partyName)
		}
		return "", fmt.Errorf("database query failed: %v", err)
	}
	return address, nil
}

func updatePModeTemplate(partyName string) error {
	address, err := getAddressFromDB(partyName)
	if err != nil {
		return err
	}

	templateFile := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-A\repository\pmodes\ex-pm-push-init.xml`
	pmodeContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read P-Mode template: %v", err)
	}

	updatedContent := strings.ReplaceAll(string(pmodeContent), "${dynamicAddress}", address)

	err = os.WriteFile(templateFile, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to overwrite P-Mode template: %v", err)
	}

	return nil
}

// Fungsi untuk membuat file MMD dan menyertakan address dinamis
func writeMMDFile(message AS4Message, attachmentFileName, mimeType, dynamicAddress string) error {
	mmd := MessageMetaData{
		XMLNS:          "http://holodeck-b2b.org/schemas/2014/06/mmd",
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://holodeck-b2b.org/schemas/2014/06/mmd ../repository/xsd/messagemetadata.xsd",
	}
	mmd.CollaborationInfo.AgreementRef.PMode = "ex-pm-push-init"
	mmd.CollaborationInfo.ConversationId = "org:holodeckb2b:test:conversation"

	mmd.PayloadInfo.DeleteFilesAfterSubmit = false
	mmd.PayloadInfo.PartInfo.Containment = "attachment"
	mmd.PayloadInfo.PartInfo.MimeType = mimeType
	mmd.PayloadInfo.PartInfo.Location = filepath.Join("payloads", attachmentFileName)

	if strings.TrimSpace(dynamicAddress) == "" {
		return fmt.Errorf("dynamic address cannot be empty")
	}

	outputDir := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-A\data\msg_out`
	timestamp := time.Now().Format("20060102150405")
	fileName := filepath.Join(outputDir, fmt.Sprintf("%s_%s.mmd", message.MessageID, timestamp))

	mmdXML, err := xml.MarshalIndent(mmd, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal MMD to XML: %v", err)
	}

	err = os.WriteFile(fileName, mmdXML, 0644)
	if err != nil {
		return fmt.Errorf("failed to write MMD file: %v", err)
	}

	return nil
}

// Handler untuk menerima pesan AS4 dan memproses lampiran
func MessageHandler(w http.ResponseWriter, r *http.Request) {
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

	message := AS4Message{
		FromParty: r.FormValue("fromParty"),
		ToParty:   r.FormValue("toParty"),
		Service:   r.FormValue("service"),
		Action:    r.FormValue("action"),
		MessageID: r.FormValue("messageId"),
		Payload:   r.FormValue("payload"),
	}

	err := updatePModeTemplate(message.ToParty)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update P-Mode template: %v", err), http.StatusInternalServerError)
		return
	}

	dynamicAddress, err := getAddressFromDB(message.ToParty)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get address for party: %v", err), http.StatusInternalServerError)
		return
	}

	attachmentFileName := ""
	mimeType := ""

	if attachmentFile, fileHeader, err := r.FormFile("attachment"); err == nil {
		defer attachmentFile.Close()
		attachmentFileName, err = saveFile(attachmentFile, fileHeader, `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-A\data\msg_out\payloads`)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
			return
		}
		mimeType = fileHeader.Header.Get("Content-Type")
	} else {
		attachmentFileName = fmt.Sprintf("default_%s.xml", time.Now().Format("20060102150405"))
		mimeType = "application/xml"
		if err := os.WriteFile(filepath.Join(`C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-A\data\msg_out\payloads`, attachmentFileName), []byte(message.Payload), 0644); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write default XML: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if err := writeMMDFile(message, attachmentFileName, mimeType, dynamicAddress); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create MMD file: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Attachment and metadata processed successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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
