package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type DatabaseBackup struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`
}

type Backup struct {
	ID           int64            `json:"id"`
	DayOfWeek    *int             `json:"day_of_week"`
	Time         *string          `json:"time"`
	Provider     string           `json:"provider"`
	ProviderName string           `json:"provider_name"`
	Status       string           `json:"status"`
	Databases    []DatabaseBackup `json:"databases"`
	Backups      []struct {
		ID            int64   `json:"id"`
		BackupID      int64   `json:"backup_id"`
		Status        string  `json:"status"`
		RestoreStatus *string `json:"restore_status"`
		ArchivePath   string  `json:"archive_path"`
		Duration      int     `json:"duration"`
		Date          string  `json:"date"`
	} `json:"backups"`
	LastBackupTime *string `json:"last_backup_time"`
}

func (c *Client) ListBackupConfigurations(ctx context.Context, serverID int) ([]Backup, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/backups", serverID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(b *Backup, id int64) { b.ID = id })
}

type BackupCredentials struct {
	Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type BackupFrequency struct {
	Type   string `json:"type"`
	Time   string `json:"time,omitempty"`
	Day    *int   `json:"day,omitempty"`
	Custom string `json:"custom,omitempty"`
}

type CreateBackupConfigurationRequest struct {
	Provider    string            `json:"provider"`
	Credentials BackupCredentials `json:"credentials"`
	Frequency   BackupFrequency   `json:"frequency"`
	Directory   string            `json:"directory"`
	Email       string            `json:"email"`
	Retention   int               `json:"retention"`
	Databases   []int             `json:"databases"`
}

func (c *Client) CreateBackupConfiguration(ctx context.Context, serverID int, req CreateBackupConfigurationRequest) (*Backup, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/backups", serverID))
	var backup Backup
	id, err := c.PostJsonApi(ctx, path, req, &backup)
	if err != nil {
		return nil, err
	}
	backup.ID = int64(id)
	return &backup, nil
}

func (c *Client) UpdateBackupConfiguration(ctx context.Context, serverID, backupID int, req CreateBackupConfigurationRequest) (*Backup, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/backups/%d", serverID, backupID))
	var backup Backup
	id, err := c.PutJsonApi(ctx, path, req, &backup)
	if err != nil {
		return nil, err
	}
	backup.ID = int64(id)
	return &backup, nil
}

func (c *Client) GetBackupConfiguration(ctx context.Context, serverID, backupID int) (*Backup, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/backups/%d", serverID, backupID))
	var backup Backup
	id, err := c.GetJsonApi(ctx, path, &backup)
	if err != nil {
		return nil, err
	}
	backup.ID = int64(id)
	return &backup, nil
}

func (c *Client) RunBackupConfiguration(ctx context.Context, serverID, backupID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/backups/%d", serverID, backupID))
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) DeleteBackupConfiguration(ctx context.Context, serverID, backupID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/backups/%d", serverID, backupID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type RestoreBackupRequest struct {
	Database int `json:"database,omitempty"`
}

func (c *Client) RestoreBackup(ctx context.Context, serverID, backupID, backupItemID int, database int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/backups/%d/backups/%d", serverID, backupID, backupItemID))
	req := RestoreBackupRequest{Database: database}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) DeleteBackup(ctx context.Context, serverID, backupID, backupItemID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/backups/%d/backups/%d", serverID, backupID, backupItemID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
