package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"syscall"
)

var (
	httpAddr        = ":8080"
	RDSUserName     = GetEnv("RDS_USERNAME", "admin")
	RDSPass         = GetEnv("RDS_PASS", "Sisi4ever")
	RDSAddr         = GetEnv("RDS_ADDR", "db-pf-1.cb2q0ky6cy86.us-east-2.rds.amazonaws.com")
	RDSDatabaseName = GetEnv("RDS_DATABASE_NAME", "alertas")
)

func main() {

	fmt.Println("🧩 RDS Configuration:")
	fmt.Println("  Username:", RDSUserName)
	fmt.Println("  Password:", RDSPass)
	fmt.Println("  Address :", RDSAddr)
	fmt.Println("  Database:", RDSDatabaseName)

	mux := http.NewServeMux()

	RDSClient, err := NewRDSStorage(
		RDSUserName,     // usuario
		RDSPass,         // contraseña
		RDSAddr,         // endpoint RDS
		"3306",          // puerto
		RDSDatabaseName, // nombre BD
	)
	if err != nil {
		log.Fatal("Error connecting to RDS: ", err)
	}

	store := NewStore(RDSClient.db)
	handler := NewHandler(store)
	handler.registerRoutes(mux)

	// ✅ Apply CORS middleware
	corsMux := enableCORS(mux)

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		InsecureSkipVerify:       true, // ⚠️ Acepta cualquier certificado (solo para pruebas)
		PreferServerCipherSuites: true,
	}

	server := &http.Server{
		Addr:      httpAddr,
		Handler:   corsMux,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting HTTPS server at %s", httpAddr)
	if err := server.ListenAndServeTLS("/app/certs/fullchain.pem",
		"/app/certs/privkey.pem"); err != nil {
		log.Fatal("Failed to start https server:", err)
	}
}

func GetEnv(key string, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adjust the origin as needed
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
