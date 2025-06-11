package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/chand-magar/SolidBaseGoStructure/internal/dto"
	"github.com/chand-magar/SolidBaseGoStructure/internal/interfaces"
	"github.com/chand-magar/SolidBaseGoStructure/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const __ROLE_TBL__ = "master.roles"
const __PROFILE_TBL__ = "master.users"
const __CREDENTIAL_TBL__ = "master.user_credentials"

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(data dto.RequestDTO) (uuid.UUID, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return uuid.Nil, tx.Error
	}

	// Rollback on panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("panic during Create transaction: %v", r)
		}
	}()

	insertFields := make(map[string]interface{})

	if data.UserFullName != "" {
		insertFields["user_full_name"] = data.UserFullName
	}
	if data.RoleId != uuid.Nil {
		var roleNo int
		query := fmt.Sprintf(`SELECT role_no FROM %s WHERE role_id = ?`, __ROLE_TBL__)
		if err := tx.Raw(query, data.RoleId).Scan(&roleNo).Error; err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		insertFields["role_no"] = roleNo
	}
	if data.EmailId != "" {
		insertFields["email_id"] = data.EmailId
	}
	if data.Gender != "" {
		insertFields["gender"] = data.Gender
	}
	if data.Dob != nil {
		insertFields["dob"] = data.Dob
	}
	if data.MobileNo != "" {
		insertFields["mobile_no"] = data.MobileNo
	}
	if len(data.Address) != 0 {
		insertFields["address"] = data.Address
	}

	currentTime := time.Now().UTC()
	profileID := uuid.New()

	insertFields["profile_id"] = profileID
	insertFields["created_at"] = currentTime

	userCols, userVals, userArgs := buildSQLParts(insertFields)

	rawQuery := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) RETURNING profile_no`,
		__PROFILE_TBL__, userCols, userVals,
	)

	var profileNo int
	if err := tx.Raw(rawQuery, userArgs...).Scan(&profileNo).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("failed to insert profile: %w", err)
	}

	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("password hashing failed: %w", err)
	}

	credentialID := uuid.New()
	credFields := map[string]interface{}{
		"credential_id": credentialID,
		"profile_no":    profileNo,
		"username":      data.Username,
		"password":      hashedPassword,
		"created_at":    currentTime,
	}

	credCols, credVals, credArgs := buildSQLParts(credFields)
	credQuery := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, __CREDENTIAL_TBL__, credCols, credVals)

	if err := tx.Exec(credQuery, credArgs...).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("failed to insert credentials: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return uuid.Nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return profileID, nil
}

func (r *userRepo) FindOne(id uuid.UUID) (*dto.ResponseDTO, error) {

	var user dto.ResponseDTO

	query := fmt.Sprintf(`
		SELECT profile.profile_no, 
			profile.profile_id, 
			profile.user_full_name,
			profile.email_id,
			profile.gender,
			profile.dob,
			profile.mobile_no,
			profile.address,
			profile.status,
			profile.created_at,
			profile.created_by,
			profile.updated_at,
			profile.updated_by,
			role.role_id,
			role.role_name
		FROM %s AS profile

	INNER JOIN %s AS role 
		ON role.role_no = profile.role_no
	
	WHERE profile.profile_id = ?
		LIMIT 1`, __PROFILE_TBL__, __ROLE_TBL__)

	err := r.db.Raw(query, id).Scan(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetAll(params dto.PaginationParams) ([]dto.ResponseDTO, int64, error) {

	var users []dto.ResponseDTO
	var total int64

	where := "WHERE 1=1"
	args := []interface{}{}

	if params.Search != "" {
		where += " AND (profile.user_full_name ILIKE ? OR profile.email_id ILIKE ?)"
		search := "%" + params.Search + "%"
		args = append(args, search, search)
	}

	if params.Status != "" {
		where += " AND profile.status = ?"
		args = append(args, params.Status)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s AS profile %s", __PROFILE_TBL__, where)
	if err := r.db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Order != "ASC" && params.Order != "DESC" {
		params.Order = "ASC"
	}
	if params.SortBy == "" {
		params.SortBy = "user_full_name"
	}

	offset := (params.Page - 1) * params.Size

	query := fmt.Sprintf(`
		SELECT profile.profile_no, 
			profile.profile_id,
			profile.user_full_name,
			profile.email_id,
			profile.gender,
			profile.dob,
			profile.mobile_no,
			profile.address,
			profile.status,
			profile.created_at,
			profile.updated_at,
			role.role_id,
			role.role_name
		FROM %s AS profile
		INNER JOIN %s AS role ON role.role_no = profile.role_no
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?`,
		__PROFILE_TBL__, __ROLE_TBL__, where, params.SortBy, params.Order)

	args = append(args, params.Size, offset)

	if err := r.db.Raw(query, args...).Scan(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepo) Update(id uuid.UUID, data dto.RequestDTO) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	updateFields := map[string]interface{}{}

	if data.UserFullName != "" {
		updateFields["user_full_name"] = data.UserFullName
	}
	if data.RoleId != uuid.Nil {
		var roleNo int
		query := fmt.Sprintf(`SELECT role_no FROM %s WHERE role_id = ?`, __ROLE_TBL__)
		if err := tx.Raw(query, data.RoleId).Scan(&roleNo).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to fetch role_no: %v", err)
		}
		updateFields["role_no"] = roleNo
	}
	if data.EmailId != "" {
		updateFields["email_id"] = data.EmailId
	}
	if data.Gender != "" {
		updateFields["gender"] = data.Gender
	}
	if data.Dob != nil {
		updateFields["dob"] = data.Dob
	}
	if data.MobileNo != "" {
		updateFields["mobile_no"] = data.MobileNo
	}
	if len(data.Address) != 0 {
		updateFields["address"] = data.Address
	}
	if data.Status != "" {
		updateFields["status"] = data.Status
	}

	if len(updateFields) == 0 {
		tx.Rollback()
		return fmt.Errorf("no fields provided for update")
	}

	updateFields["updated_at"] = time.Now().UTC()
	updateFields["updated_by"] = data.UpdatedBy

	query := fmt.Sprintf("UPDATE %s SET ", __PROFILE_TBL__)
	values := []interface{}{}
	i := 0
	for field, value := range updateFields {
		if i > 0 {
			query += ", "
		}
		query += field + " = ?"
		values = append(values, value)
		i++
	}
	query += " WHERE profile_id = ?"
	values = append(values, id)

	if err := tx.Exec(query, values...).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func buildSQLParts(fields map[string]interface{}) (columns string, values string, args []interface{}) {
	i := 1
	for col, val := range fields {
		if i > 1 {
			columns += ", "
			values += ", "
		}
		columns += col
		values += fmt.Sprintf("$%d", i)
		args = append(args, val)
		i++
	}
	return
}
