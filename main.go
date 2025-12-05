package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	gnet "github.com/shirou/gopsutil/v3/net"
)

// -------------------------------------------------------------------
// Embebemos los archivos HTML (carpeta templates)
// -------------------------------------------------------------------

//go:embed templates/*
var templatesFS embed.FS

// -------------------------------------------------------------------
// Modelo de conexión que usará la API
// -------------------------------------------------------------------

type Connection struct {
	RemoteIP        string `json:"remote_ip"`
	Country         string `json:"country"`
	Range           string `json:"range"`
	ASN             string `json:"asn"`
	Protocol        string `json:"protocol"`
	Since           string `json:"since"`            // ISO RFC3339
	DurationSecs    int64  `json:"duration_secs"`    // duración en segundos
	DisplaySince    string `json:"display_since"`    // hora legible
	DisplayDuration string `json:"display_duration"` // duración legible
}

// -------------------------------------------------------------------
// Gestión de tiempo de conexión (primera vez vista)
// -------------------------------------------------------------------

var (
	connFirstSeen = make(map[string]time.Time)
	connMu        sync.Mutex
)

// clave única por conexión (local+remoto+tipo)
func connKey(c gnet.ConnectionStat) string {
	return fmt.Sprintf(
		"%s:%d->%s:%d:%d",
		c.Laddr.IP, c.Laddr.Port,
		c.Raddr.IP, c.Raddr.Port,
		c.Type,
	)
}

// -------------------------------------------------------------------
// Lógica para obtener conexiones ESTABLISHED reales
// -------------------------------------------------------------------

// getEstablishedConnections devuelve conexiones TCP ESTABLISHED reales
func getEstablishedConnections() []Connection {
	// Obtenemos todas las conexiones TCP del sistema
	conns, err := gnet.Connections("tcp")
	if err != nil {
		log.Println("error obteniendo conexiones:", err)
		return nil
	}

	now := time.Now()
	result := make([]Connection, 0, len(conns))

	connMu.Lock()
	defer connMu.Unlock()

	for _, c := range conns {
		// Filtrar solo ESTABLISHED y que tenga IP remota
		if c.Status != "ESTABLISHED" || c.Raddr.IP == "" {
			continue
		}

		key := connKey(c)
		first, ok := connFirstSeen[key]
		if !ok {
			first = now
			connFirstSeen[key] = first
		}

		dur := now.Sub(first)
		secs := int64(dur.Seconds())

		result = append(result, Connection{
			RemoteIP:        c.Raddr.IP,
			Country:         "", // luego lo llenaremos con GeoIP
			Range:           "",
			ASN:             "",
			Protocol:        "TCP",
			Since:           first.Format(time.RFC3339),
			DurationSecs:    secs,
			DisplaySince:    first.Format("15:04:05"),
			DisplayDuration: formatDurationHuman(secs),
		})
	}

	return result
}

// formatDurationHuman convierte segundos en texto tipo "1m 20s"
func formatDurationHuman(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	m := seconds / 60
	s := seconds % 60
	if m < 60 {
		if s == 0 {
			return fmt.Sprintf("%dm", m)
		}
		return fmt.Sprintf("%dm %ds", m, s)
	}
	h := m / 60
	m = m % 60
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}

// -------------------------------------------------------------------
// Handlers HTTP (dashboard + API)
// -------------------------------------------------------------------

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		http.Error(w, "Error cargando plantilla", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	data := map[string]any{
		"Title": "WiMon · Windows Lemon Monitor",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error renderizando plantilla", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}
}

// API: devolver lista de conexiones (ahora reales ESTABLISHED)
func connectionsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	conns := getEstablishedConnections()

	if err := json.NewEncoder(w).Encode(conns); err != nil {
		http.Error(w, "Error codificando JSON", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}
}

// -------------------------------------------------------------------
// main: arranca el servidor HTTP
// -------------------------------------------------------------------

func main() {
	mux := http.NewServeMux()

	// Dashboard
	mux.HandleFunc("/", dashboardHandler)

	// API de conexiones
	mux.HandleFunc("/api/connections", connectionsAPIHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("WiMon escuchando en http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
