package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type AS4Message struct {
	FromParty string `json:"fromParty"`
	ToParty   string `json:"toParty"`
	Service   string `json:"service"`
	Action    string `json:"action"`
	MessageID string `json:"messageId"`
	Payload   string `json:"payload"`
	Subject   string `json:"subject"`
}

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

type Partner struct {
	PartyID     string `json:"partyid"`
	Name        string `json:"name"`
	EndpointURL string `json:"endpoint_url"`
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

func replacePlaceholders(template, address, partyID string) (string, error) {
	if template == "" {
		return "", fmt.Errorf("template is empty")
	}
	replaced := strings.ReplaceAll(template, "${dynamic_responder_party_id}", strings.TrimSpace(partyID))
	replaced = strings.ReplaceAll(replaced, "${dynamic_address}", strings.TrimSpace(address))
	return replaced, nil
}

func updatePModeTemplate(partyName string) error {
	dynamicAddress, partyID, err := getAddressFromDB(partyName)
	if err != nil {
		return fmt.Errorf("failed to get address and partyID for party %s: %v", partyName, err)
	}
	log.Printf("Dynamic Address: %s, Party ID: %s", dynamicAddress, partyID)

	templateFile := "/opt/holodeckb2b/examples/pmodes/pm-push.xml"

	pmodeContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read P-Mode template: %v", err)
	}

	updatedContent, err := replacePlaceholders(string(pmodeContent), dynamicAddress, partyID)
	if err != nil {
		return fmt.Errorf("failed to replace placeholders in template: %v", err)
	}

	activePMode := "/opt/holodeckb2b/repository/pmodes/current-pmode.xml"

	if err := os.WriteFile(activePMode, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to overwrite active P-Mode file: %v", err)
	}

	log.Printf("P-Mode successfully updated for party: %s", partyName)
	return nil
}

// Tambahkan di bawah func SavePModeFileHandler
func UpdatePModeHandler(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		ToParty string `json:"toParty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.ToParty == "" {
		http.Error(w, "toParty is required", http.StatusBadRequest)
		return
	}
	if err := updatePModeTemplate(req.ToParty); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update P-Mode template: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"PMode updated successfully"}`))
}

func writePayloadAsSOAP(message AS4Message, outputDir, soapPayloadFileName string) error {
	payloadPath := filepath.Join(outputDir, soapPayloadFileName)

	type SoapEnvelope struct {
		XMLName xml.Name `xml:"SOAP:Envelope"`
		SOAP    string   `xml:"xmlns:SOAP,attr"`
		EB      string   `xml:"xmlns:eb,attr"`
		MyNS    string   `xml:"xmlns:myns,attr"`
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
						Subject string `xml:"eb:Subject"`
					} `xml:"eb:CollaborationInfo"`
				} `xml:"eb:UserMessage"`
			} `xml:"eb:Messaging"`
		} `xml:"SOAP:Header"`
		Body struct {
			MessageContent struct {
				Document string `xml:",innerxml"`
			} `xml:"myns:MessageContent"`
		} `xml:"SOAP:Body"`
	}

	soapEnv := SoapEnvelope{
		SOAP: "http://www.w3.org/2003/05/soap-envelope",
		EB:   "http://docs.oasis-open.org/ebxml-msg/ebms/v3.0/ns/core/200704",
		MyNS: "http://example.org/myns",
	}

	// Isi Header
	soapEnv.Header.Messaging.UserMessage.MessageInfo.MessageId = message.MessageID
	soapEnv.Header.Messaging.UserMessage.MessageInfo.Timestamp = time.Now().UTC().Format(time.RFC3339)
	soapEnv.Header.Messaging.UserMessage.PartyInfo.From.PartyId = message.FromParty
	soapEnv.Header.Messaging.UserMessage.PartyInfo.To.PartyId = message.ToParty
	soapEnv.Header.Messaging.UserMessage.CollaborationInfo.Service = message.Service
	soapEnv.Header.Messaging.UserMessage.CollaborationInfo.Action = message.Action
	soapEnv.Header.Messaging.UserMessage.CollaborationInfo.Subject = message.Subject

	// Isi Body dengan XML langsung (bisa kosong)
	soapEnv.Body.MessageContent.Document = message.Payload

	soapXML, err := xml.MarshalIndent(soapEnv, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SOAP payload: %v", err)
	}

	var buffer bytes.Buffer
	buffer.WriteString(xml.Header)
	buffer.Write(soapXML)

	return os.WriteFile(payloadPath, buffer.Bytes(), 0644)
}

func writeMMDFile(message AS4Message, soapFileName string, attachmentFileNames []string, mimeTypes []string, outputDir string) error {
	mmd := MessageMetaData{
		XMLNS:          "http://holodeck-b2b.org/schemas/2014/06/mmd",
		XSI:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://holodeck-b2b.org/schemas/2014/06/mmd ../repository/xsd/messagemetadata.xsd",
	}

	mmd.MessageInfo.MessageId = message.MessageID
	mmd.MessageInfo.Timestamp = time.Now().UTC().Format(time.RFC3339)

	mmd.CollaborationInfo.AgreementRef.PMode = "current-pmode-push"
	mmd.CollaborationInfo.Service = message.Service
	mmd.CollaborationInfo.Action = message.Action
	mmd.CollaborationInfo.ConversationId = "org:holodeckb2b:test:conversation"

	mmd.PayloadInfo.DeleteFilesAfterSubmit = false

	type PartInfo struct {
		URI         string `xml:"uri,attr"`
		Containment string `xml:"containment,attr"`
		MimeType    string `xml:"mimeType,attr"`
		Location    string `xml:"location,attr"`
	}

	// Tambahkan SOAP payload jika ada
	if soapFileName != "" {
		mmd.PayloadInfo.PartInfo = append(mmd.PayloadInfo.PartInfo, PartInfo{
			URI:         "soapPart",
			Containment: "attachment",
			MimeType:    "application/xml",
			Location:    path.Join("payloads", soapFileName),
		})
	}

	// Tambahkan attachment dengan URI unik
	for i, fileName := range attachmentFileNames {
		mimeType := mimeTypes[i]
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		mmd.PayloadInfo.PartInfo = append(mmd.PayloadInfo.PartInfo, PartInfo{
			URI:         fmt.Sprintf("part%d", i+1),
			Containment: "attachment",
			MimeType:    mimeType,
			Location:    path.Join("payloads", fileName),
		})
	}

	mmdFileName := fmt.Sprintf("%s.mmd", message.MessageID)
	mmdFilePath := filepath.Join(outputDir, mmdFileName)

	mmdXML, err := xml.MarshalIndent(mmd, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal MMD to XML: %v", err)
	}

	var buffer bytes.Buffer
	buffer.WriteString(xml.Header)
	buffer.Write(mmdXML)

	return os.WriteFile(mmdFilePath, buffer.Bytes(), 0644)
}

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

	if err := r.ParseMultipartForm(20 << 20); err != nil {
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
		Subject:   r.FormValue("subject"),
	}

	// Generate MessageID jika kosong
	if message.MessageID == "" {
		message.MessageID = fmt.Sprintf("msg_%d", time.Now().UnixNano())
	}

	payloadDir := "/opt/holodeckb2b/data/msg_out/payloads"
	if err := os.MkdirAll(payloadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create payload directory", http.StatusInternalServerError)
		return
	}

	attachmentFileNames := []string{}
	mimeTypes := []string{}
	totalSize := int64(0)

	files := r.MultipartForm.File["attachments"]
	if len(files) > 5 {
		http.Error(w, "Too many attachments (maximum is 5)", http.StatusBadRequest)
		return
	}

	// Proses attachments
	for _, fileHeader := range files {
		if fileHeader.Size+totalSize > 20<<20 {
			http.Error(w, "Total attachment size exceeds 20 MB", http.StatusBadRequest)
			return
		}
		totalSize += fileHeader.Size

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to open attachment: %v", err), http.StatusInternalServerError)
			return
		}

		destPath := filepath.Join(payloadDir, fileHeader.Filename)
		destFile, err := os.Create(destPath)
		if err != nil {
			file.Close()
			http.Error(w, fmt.Sprintf("Failed to create file: %v", err), http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(destFile, file)
		file.Close()
		destFile.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
			return
		}

		attachmentFileNames = append(attachmentFileNames, fileHeader.Filename)
		mimeType := fileHeader.Header.Get("Content-Type")
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		mimeTypes = append(mimeTypes, mimeType)
		log.Printf("Saved attachment: %s (Size: %d bytes, MIME: %s)", fileHeader.Filename, fileHeader.Size, mimeType)
	}

	var soapPayloadFileName string
	// Buat SOAP payload hanya jika ada payload atau tidak ada attachment
	if message.Payload != "" || len(attachmentFileNames) == 0 {
		soapPayloadFileName = fmt.Sprintf("%s_payload.xml", message.MessageID)
		if err := writePayloadAsSOAP(message, payloadDir, soapPayloadFileName); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write SOAP payload: %v", err), http.StatusInternalServerError)
			return
		}
	}

	outputDir := "/opt/holodeckb2b/data/msg_out"
	if err := writeMMDFile(message, soapPayloadFileName, attachmentFileNames, mimeTypes, outputDir); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create MMD file: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Message processed successfully",
		"messageId": message.MessageID,
	})
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		":", "_", "/", "_", "\\", "_",
		"*", "_", "?", "_", "\"", "_",
		"<", "_", ">", "_", "|", "_",
	)
	return replacer.Replace(name)
}

// Fungsi untuk generate PMode dari template
func GeneratePModeFromTemplate(templatePath, outputPath, senderPartyId string) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}
	replaced := strings.ReplaceAll(string(content), "${sender}", senderPartyId)
	return os.WriteFile(outputPath, []byte(replaced), 0644)
}

// Handler utama untuk nambah partner
func AddPartnerHandler(w http.ResponseWriter, r *http.Request) {
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

	var data struct {
		PartyID     string `json:"partyid"`
		Name        string `json:"name"`
		EndpointURL string `json:"endpoint_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if data.PartyID == "" || data.Name == "" || data.EndpointURL == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received: partyid=%s, name=%s, endpoint_url=%s\n", data.PartyID, data.Name, data.EndpointURL)

	_, err := db.Exec("INSERT INTO party (party_id, name, endpoint_url) VALUES (?, ?, ?)", data.PartyID, data.Name, data.EndpointURL)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	templatePath := "/opt/holodeckb2b/examples/pmodes/ex-pm-push-resp.xml"
	safeName := sanitizeFilename(data.Name)
	outputPath := fmt.Sprintf("/opt/holodeckb2b/repository/pmodes/pmode-resp-%s.xml", safeName)

	err = GeneratePModeFromTemplate(templatePath, outputPath, data.PartyID)
	if err != nil {
		http.Error(w, "Failed to generate PMode file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Partner added and PMode response generated",
	})
}

func GetPartnersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT party_id, name, endpoint_url FROM party")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var partners []Partner
	for rows.Next() {
		var p Partner
		if err := rows.Scan(&p.PartyID, &p.Name, &p.EndpointURL); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		partners = append(partners, p)
	}
	json.NewEncoder(w).Encode(partners)
}

// Handler: Update partner
func UpdatePartnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	var p Partner
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if p.PartyID == "" || p.Name == "" || p.EndpointURL == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE party SET name=?, endpoint_url=? WHERE party_id=?", p.Name, p.EndpointURL, p.PartyID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Partner updated"})
}

// Handler: Delete partner
func DeletePartnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	var p Partner
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if p.PartyID == "" {
		http.Error(w, "PartyID is required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM party WHERE party_id=?", p.PartyID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Partner deleted"})
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}

	// Ganti path log sesuai aplikasi Anda
	logPath := "/opt/holodeckb2b/logs/holodeckb2b.log"
	data, err := ioutil.ReadFile(logPath)
	if err != nil {
		http.Error(w, "Gagal membaca log: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}

func ListPModeFilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dir := "/opt/holodeckb2b/repository/pmodes"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		http.Error(w, "Gagal membaca folder: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var xmlFiles []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".xml" {
			xmlFiles = append(xmlFiles, f.Name())
		}
	}
	json.NewEncoder(w).Encode(xmlFiles)
}

// Get content of a PMode XML file
func GetPModeFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	name := r.URL.Query().Get("name")
	if name == "" || filepath.Ext(name) != ".xml" {
		http.Error(w, "Nama file tidak valid", http.StatusBadRequest)
		return
	}
	path := filepath.Join("/opt/holodeckb2b/repository/pmodes", name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "Gagal membaca file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/xml")
	w.Write(data)
}

// Save (update) content of a PMode XML file
func SavePModeFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Name == "" || filepath.Ext(req.Name) != ".xml" {
		http.Error(w, "Nama file tidak valid", http.StatusBadRequest)
		return
	}
	path := filepath.Join("/opt/holodeckb2b/repository/pmodes", req.Name)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "Gagal membuka file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close() // File akan langsung ditutup setelah fungsi selesai
	if _, err := f.WriteString(req.Content); err != nil {
		http.Error(w, "Gagal menulis file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"message":"File berhasil disimpan"}`))
}

func PartnerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetPartnersHandler(w, r)
	case http.MethodPost:
		AddPartnerHandler(w, r)
	case http.MethodPut:
		UpdatePartnerHandler(w, r)
	case http.MethodDelete:
		DeletePartnerHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

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

	http.HandleFunc("/api/as4/send", MessageHandler)
	http.HandleFunc("/api/partner", PartnerHandler)
	http.HandleFunc("/api/log", LogHandler)
	http.HandleFunc("/api/pmode/update", UpdatePModeHandler)
	http.HandleFunc("/api/pmode/list", ListPModeFilesHandler)
	http.HandleFunc("/api/pmode/get", GetPModeFileHandler)
	http.HandleFunc("/api/pmode/save", SavePModeFileHandler)
	log.Println("Starting server on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
