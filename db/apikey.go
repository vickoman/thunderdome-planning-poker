package db

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

// GenerateApiKey generates a new API key for a User
func (d *Database) GenerateApiKey(UserID string, KeyName string) (*model.APIKey, error) {
	apiPrefix, prefixErr := randomString(8)
	if prefixErr != nil {
		err := errors.New("error generating api prefix")
		log.Println(err)
		log.Println(prefixErr)
		return nil, err
	}

	apiSecret, secretErr := randomString(32)
	if secretErr != nil {
		err := errors.New("error generating api secret")
		log.Println(err)
		log.Println(secretErr)
		return nil, err
	}

	APIKEY := &model.APIKey{
		Name:        KeyName,
		Key:         apiPrefix + "." + apiSecret,
		UserId:      UserID,
		Prefix:      apiPrefix,
		Active:      true,
		CreatedDate: time.Now(),
	}
	hashedKey := hashString(APIKEY.Key)
	keyID := apiPrefix + "." + hashedKey

	e := d.db.QueryRow(
		`SELECT createdDate FROM user_apikey_add($1, $2, $3);`,
		keyID,
		KeyName,
		UserID,
	).Scan(&APIKEY.CreatedDate)
	if e != nil {
		log.Println(e)
		return nil, errors.New("unable to create new api key")
	}

	return APIKEY, nil
}

// GetUserApiKeys gets a list of api keys for a user
func (d *Database) GetUserApiKeys(UserID string) ([]*model.APIKey, error) {
	var APIKeys = make([]*model.APIKey, 0)
	rows, err := d.db.Query(
		"SELECT id, name, user_id, active, created_date, updated_date FROM api_keys WHERE user_id = $1 ORDER BY created_date",
		UserID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak model.APIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.UserId,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.Id = key
				APIKeys = append(APIKeys, &ak)
			}
		}
	}

	return APIKeys, err
}

// UpdateUserApiKey updates a user api key (active column only)
func (d *Database) UpdateUserApiKey(UserID string, KeyID string, Active bool) ([]*model.APIKey, error) {
	if _, err := d.db.Exec(
		`CALL user_apikey_update($1, $2, $3);`, KeyID, UserID, Active); err != nil {
		log.Println(err)
		return nil, err
	}

	keys, keysErr := d.GetUserApiKeys(UserID)
	if keysErr != nil {
		log.Println(keysErr)
		return nil, keysErr
	}

	return keys, nil
}

// DeleteUserApiKey removes a users api key
func (d *Database) DeleteUserApiKey(UserID string, KeyID string) ([]*model.APIKey, error) {
	if _, err := d.db.Exec(
		`CALL user_apikey_delete($1, $2);`, KeyID, UserID); err != nil {
		log.Println(err)
		return nil, err
	}

	keys, keysErr := d.GetUserApiKeys(UserID)
	if keysErr != nil {
		log.Println(keysErr)
		return nil, keysErr
	}

	return keys, nil
}

// GetApiKeyUser checks to see if the API key exists and returns the User
func (d *Database) GetApiKeyUser(APK string) (*model.User, error) {
	User := &model.User{}

	splitKey := strings.Split(APK, ".")
	hashedKey := hashString(APK)
	keyID := splitKey[0] + "." + hashedKey

	e := d.db.QueryRow(`
		SELECT u.id, u.name, u.email, u.type, u.avatar, u.verified, u.notifications_enabled, COALESCE(u.country, ''), COALESCE(u.locale, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), u.created_date, u.updated_date, u.last_active 
		FROM api_keys ak
		LEFT JOIN users u ON u.id = ak.user_id
		WHERE ak.id = $1 AND ak.active = true
`,
		keyID,
	).Scan(
		&User.Id,
		&User.Name,
		&User.Email,
		&User.Type,
		&User.Avatar,
		&User.Verified,
		&User.NotificationsEnabled,
		&User.Country,
		&User.Locale,
		&User.Company,
		&User.JobTitle,
		&User.CreatedDate,
		&User.UpdatedDate,
		&User.LastActive)
	if e != nil {
		log.Println(e)
		return nil, errors.New("active API Key match not found")
	}

	User.GravatarHash = createGravatarHash(User.Email)

	return User, nil
}
