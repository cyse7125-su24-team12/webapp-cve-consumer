package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func computeHash(jsonData []byte) string {
    hash := sha256.Sum256(jsonData)
    return hex.EncodeToString(hash[:])
}

func checkExistingRecord(db *sql.DB, cveID, dataHash string) (bool, int, error) {
    query := `SELECT version FROM cve.cve_details WHERE cve_id = $1 AND cve_data_hash = $2 ORDER BY version DESC LIMIT 1`
    var version int
    err := db.QueryRow(query, cveID, dataHash).Scan(&version)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, 0, nil
        }
        return false, 0, err
    }
    return true, version, nil
}

// insertJSONData inserts a JSON object into the PostgreSQL database
func insertJSONData(db *sql.DB, cveID string, jsonData []byte, version int, dataHash string) error {
    query := `INSERT INTO cve.cve_details (cve_id, cve_data, version, cve_data_hash) VALUES ($1, $2, $3, $4)`
    _, err := db.Exec(query, cveID, jsonData, version, dataHash)
    return err
}

func SaveToDb(data []byte, db *sql.DB) error {
    var jsonObj map[string]interface{}
    if err := json.Unmarshal(data, &jsonObj); err != nil {
        return err
    }

    cveID := getCveId(jsonObj)
    dataHash := computeHash(data)
	
	exists, version, err := checkExistingRecord(db, cveID, dataHash)
    if err != nil {
        return err
    }
    if exists {
        fmt.Println("No operation needed, duplicate data.")
        return nil
    }

	if version != 0 {
        version++
    } else {
        version = 1
    }

    if err := insertJSONData(db, cveID, data, version, dataHash); err != nil {
        return err
    }

    return nil
}

func getCveId(jsonObj map[string]interface{}) string {
    if metaData, ok := jsonObj["cveMetadata"].(map[string]interface{}); ok {
        if cveID, ok := metaData["cveId"].(string); ok {
            return cveID
        }
    }
    return ""
}