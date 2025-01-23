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
	MessageContent string `xml:"messageContent"`
	PayloadInfo    struct {
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

func savePayloadToXMLFile(payload, filePath string) error {
	if strings.TrimSpace(payload) == "" {
		return fmt.Errorf("payload is empty")
	}

	err := os.WriteFile(filePath, []byte(payload), 0644)
	if err != nil {
		return fmt.Errorf("failed to save payload to XML: %v", err)
	}

	return nil
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

// Fungsi untuk membuat file MMD
func writeMMDFile(message AS4Message, attachmentFileName, mimeType string) error {
	mmd := MessageMetaData{
		XMLNS:          "http://holodeck-b2b.org/schemas/2014/06/mmd",
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://holodeck-b2b.org/schemas/2014/06/mmd ../repository/xsd/messagemetadata.xsd",
	}
	mmd.CollaborationInfo.AgreementRef.PMode = "current-pmode-push"
	mmd.CollaborationInfo.ConversationId = "org:holodeckb2b:test:conversation"
	mmd.MessageContent = message.Payload
	// Konfigurasi PayloadInfo untuk attachment
	mmd.PayloadInfo.DeleteFilesAfterSubmit = false
	mmd.PayloadInfo.PartInfo = struct {
		Containment string `xml:"containment,attr"`
		MimeType    string `xml:"mimeType,attr"`
		Location    string `xml:"location,attr"`
	}{
		Containment: "attachment",
		MimeType:    mimeType,
		Location:    filepath.Join("payloads", attachmentFileName),
	}

	// Path untuk direktori output
	outputDir := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-A\data\msg_out`
	timestamp := time.Now().Format("20060102150405")
	fileName := filepath.Join(outputDir, fmt.Sprintf("%s_%s.mmd", message.MessageID, timestamp))

	// Serialize ke XML
	mmdXML, err := xml.MarshalIndent(mmd, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal MMD to XML: %v", err)
	}

	// Tulis file MMD
	err = os.WriteFile(fileName, mmdXML, 0644)
	if err != nil {
		return fmt.Errorf("failed to write MMD file: %v", err)
	}

	log.Printf("MMD file created: %s", fileName)
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

	message := AS4Message{
		FromParty: r.FormValue("fromParty"),
		ToParty:   r.FormValue("toParty"),
		Service:   r.FormValue("service"),
		Action:    r.FormValue("action"),
		MessageID: r.FormValue("messageId"),
		Payload:   r.FormValue("payload"),
	}

	// Update P-Mode menggunakan nilai dynamicAddress dan partyID
	if err := updatePModeTemplate(message.ToParty); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update P-Mode template: %v", err), http.StatusInternalServerError)
		return
	}

	// Direktori payload
	payloadDir := `C:\Users\Yusuf\Documents\Kuliah\RPLK\Tugas Akhir\holodeckb2b-7.0.0-A\data\msg_out\payloads`
	attachmentFileName := ""
	mimeType := ""

	// Cek dan simpan attachment jika ada
	if attachmentFile, fileHeader, err := r.FormFile("attachment"); err == nil {
		defer attachmentFile.Close()
		attachmentFileName, err = saveFile(attachmentFile, fileHeader, payloadDir)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
			return
		}
		mimeType = fileHeader.Header.Get("Content-Type")
	} else {
		// Jika tidak ada attachment, gunakan payload sebagai default file
		attachmentFileName = fmt.Sprintf("default_%s.xml", time.Now().Format("20060102150405"))
		mimeType = "application/xml"
		if err := savePayloadToXMLFile(message.Payload, filepath.Join(payloadDir, attachmentFileName)); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write default XML: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// **Selalu simpan payload sebagai file XML terpisah**
	payloadFileName := fmt.Sprintf("%s_payload.xml", message.MessageID)
	if err := savePayloadToXMLFile(message.Payload, filepath.Join(payloadDir, payloadFileName)); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save payload XML: %v", err), http.StatusInternalServerError)
		return
	}

	// **Tulis file MMD dengan kedua file: attachment dan payload**
	if err := writeMMDFile(message, attachmentFileName, mimeType); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write MMD: %v", err), http.StatusInternalServerError)
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
