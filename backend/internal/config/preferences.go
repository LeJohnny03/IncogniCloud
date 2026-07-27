package config

// UserPreferences repräsentiert die individuellen Einstellungen eines Nutzers,
// die als JSONB in der Datenbank gespeichert werden.
type ServerPreferences struct {
	StorageDrives   []string `json:"storageDrives"` // Array für mehrere Festplatten!
	StorageFolder   string   `json:"storageFolder"`
	EnableBackups   bool     `json:"enableBackups"`
	BackupDrive     string   `json:"backupDrive"`
	BackupType      string   `json:"backupType"`
	CloudFolderName string   `json:"cloudFolderName"`
}
