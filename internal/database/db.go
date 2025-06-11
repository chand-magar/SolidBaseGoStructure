package db

import (
	"fmt"
	"log"

	"github.com/chand-magar/SolidBaseGoStructure/internal/models"

	utils "github.com/chand-magar/SolidBaseGoStructure/internal/utils"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// Initialize the database connection
func InitDB(dataSourceName string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable the master schema
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS master").Error; err != nil {
		return nil, fmt.Errorf("failed to create master schema: %w", err)
	}

	// Run Migrations
	// MigrateAllTables(db)

	return db, nil
}

// Migrate all tables and set up initial data
func MigrateAllTables(db *gorm.DB) {

	// ADD ENUM
	addEnumQuery := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
			CREATE TYPE status_enum AS ENUM ('A', 'I', 'D');
		END IF;
	END $$;
	`
	if err := db.Exec(addEnumQuery).Error; err != nil {
		log.Fatalf("Failed to add status ENUM: %v", err)
	}

	// Auto-migrate all tables
	tables := []interface{}{
		&models.User{},
		&models.UsersCredentials{},
		&models.Section{},
		&models.Page{},
		&models.Role{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			log.Fatalf("Failed to migrate table %T: %v", table, err)
		}
	}

	// ADD FOREIGN KEY CONSTRAINTS
	foreignKeys := []string{
		`ALTER TABLE master.users ADD CONSTRAINT fk_role_no FOREIGN KEY (role_no) REFERENCES master.roles(role_no) ON UPDATE SET NULL;`,
		`ALTER TABLE master.user_credentials ADD CONSTRAINT fk_profile_no FOREIGN KEY (profile_no) REFERENCES master.users(profile_no) ON DELETE CASCADE;`,
		`ALTER TABLE master.pages ADD CONSTRAINT fk_section_no FOREIGN KEY (section_no) REFERENCES master.sections(section_no) ON DELETE CASCADE;`,
	}

	for _, query := range foreignKeys {
		if err := db.Exec(query).Error; err != nil {
			log.Fatalf("Failed to add foreign key constraint: %v", err)
		}
	}

	// START TRANSACTION
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("Transaction rolled back due to panic: %v", r)
		}
	}()

	// Insert Sections
	sections := []models.Section{
		{SectionId: uuid.New(), SectionName: "Webmasters", SectionPath: "webmasters", SectionIcon: "fa-laptop", SectionOrder: 1, Status: "A"},
		{SectionId: uuid.New(), SectionName: "Employees", SectionPath: "systems", SectionIcon: "fa-users", SectionOrder: 2, Status: "A"},
	}

	if err := tx.Create(&sections).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert sections: %v", err)
	}

	// Insert Pages
	pages := []models.Page{
		{PageId: uuid.New(), SectionNo: sections[0].SectionNo, PageName: "Sections", PagePath: "section", PageOrder: 1, Status: "A"},
		{PageId: uuid.New(), SectionNo: sections[0].SectionNo, PageName: "Pages", PagePath: "pages", PageOrder: 2, Status: "A"},
		{PageId: uuid.New(), SectionNo: sections[1].SectionNo, PageName: "Roles", PagePath: "roles", PageOrder: 1, Status: "A"},
		{PageId: uuid.New(), SectionNo: sections[1].SectionNo, PageName: "Employees", PagePath: "employees", PageOrder: 2, Status: "A"},
	}

	if err := tx.Create(&pages).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert pages: %v", err)
	}

	// Insert Roles
	roles := []models.Role{
		{RoleNo: 1, RoleName: "Super Admin", Status: "A"},
	}

	if err := tx.Create(&roles).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert roles: %v", err)
	}

	users := []models.User{
		{ProfileId: uuid.New(), RoleNo: 1, UserFullName: "Chand Kumar Magar", EmailId: "chand.magar@gmail.com", MobileNo: "9804590230", Status: "A"},
	}

	if err := tx.Create(&users).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert users: %v", err)
	}

	// Hash Passwords
	var Password string
	var err error

	if Password, err = utils.HashPassword("Khelaghar@786"); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to hash password for admin: %v", err)
	}

	// Insert User Credentials
	userCreds := []models.UsersCredentials{
		{CredentialId: uuid.New(), ProfileNo: users[0].ProfileNo, Username: "chand.magar", Password: Password},
	}

	if err := tx.Create(&userCreds).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert user credentials: %v", err)
	}

	// Commit Transaction if Everything is Successful
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Database migration completed successfully!")
}

// Getter function to return the database instance
func GetDB() *gorm.DB {
	return db
}
