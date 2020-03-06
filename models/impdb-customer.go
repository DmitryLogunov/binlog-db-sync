package models

type WsdbCustomer struct {
	Id           int    `gorm:"column:id"`
	Uid          string `gorm:"column:uid"`
	Pod          string `gorm:"column:pod"`
	Name         string `gorm:"column:name"`
	Status       int    `gorm:"column:status"`
	Editor       int    `gorm:"column:editor"`
	LastModified string `gorm:"column:last_modified"`
	LastSynced   string `gorm:"column:last_synced"`
	Major        int    `gorm:"column:major"`
	Email        string `gorm:"column:email"`
	Mobile       string `gorm:"column:mobile"`
	Retry        int    `gorm:"column:retry"`
	IsTest       int    `gorm:"column:is_test"`
	MailboxId    int    `gorm:"column:mailbox_id"`
}
