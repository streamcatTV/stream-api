package data

import (
	"fmt"
	"streamcat-api/models"
)

// GetStreams Gets all streams.
func GetStreams() *[]models.StreamPrivate {
	const query = `SELECT * FROM stream ORDER BY id ASC`

	db, _ := ConnectDB()
	stream := []models.StreamPrivate{}
	db.Select(&stream, query)
	fmt.Println(&stream)

	return &stream
}

// GetFeaturedStreams Gets all streams.
// TODO: Featured streams should be live only.
func GetFeaturedStreams(limit int) (*[]models.StreamPrivate, error) {
	const query = `
		SELECT * FROM stream
		WHERE private = false
		ORDER BY id DESC
		LIMIT $1
	`

	db, _ := ConnectDB()
	stream := []models.StreamPrivate{}
	err := db.Select(&stream, query, limit)
	if err != nil {
		fmt.Println(err)
		return &stream, err
	}
	fmt.Println(&stream)

	return &stream, nil
}

// GetStreamByID Gets stream by Id
func GetStreamByID(id int) (*models.StreamPrivate, error) {
	const query = `SELECT * FROM stream WHERE id = $1`

	db, _ := ConnectDB()
	stream := models.StreamPrivate{}
	err := db.Get(&stream, query, id)
	fmt.Println(&stream)
	if err != nil {
		fmt.Println(err)
		return &stream, err
	}

	return &stream, nil
}

// GetStreamByName Gets stream by Name
func GetStreamByName(name string) (*models.StreamPrivate, error) {
	const query = `SELECT * FROM stream WHERE stream_name = $1`

	db, _ := ConnectDB()
	stream := models.StreamPrivate{}
	err := db.Get(&stream, query, name)
	fmt.Println(&stream)
	if err != nil {
		fmt.Println(err)
		return &stream, err
	}

	return &stream, nil
}

// CreateStream creates a stream.
// TODO: Check for proper claims from JWT.
func CreateStream(stream *models.Stream) *models.Stream {
	const query = `INSERT INTO stream
		(title, type, description, private, stream_name, stream_key)
		VALUES (:title, :type, :description, :private, :stream_name, :stream_key)`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	result, err := tx.NamedExec(query, &stream)
	if err != nil {
		fmt.Println("Error", err)
	}
	tx.Commit()

	fmt.Println("transaction done")

	lastID, _ := result.LastInsertId()
	stream.ID = lastID

	return stream
}

// UpdateStreamByID Update stream by id
func UpdateStreamByID(id int, stream models.StreamPrivate) *models.StreamPrivate {
	const query = `UPDATE stream
		SET stream_name = :stream_name,
		type = :type,
		description = :description
		stream_name = :stream_name
		stream_key = :stream_key
		private = :private
		WHERE id = :id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &stream)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

	return &stream
}

// UpdateStreamByName Update stream by name
func UpdateStreamByName(name string, stream models.StreamPrivate) *models.StreamPrivate {
	const query = `UPDATE stream
		SET stream_name = :stream_name,
		type = :type,
		description = :description
		stream_name = :stream_name
		stream_key = :stream_key
		private = :private
		WHERE stream_name = :name`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &stream)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

	return &stream
}

// DeleteStreamByID Deletes stream by id
func DeleteStreamByID(id int) error {
	const query = "DELETE FROM stream WHERE id = $1"

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.Exec(query, id)
	tx.Commit()
	return err
}

// DeleteStreamByName Deletes stream by name
func DeleteStreamByName(name string) error {
	const query = "DELETE FROM stream WHERE stream_name = $1"

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.Exec(query, name)
	tx.Commit()
	return err
}

// StreamExistsByName Checks if stream exists.
func StreamExistsByName(streamName string) bool {
	const query = `SELECT EXISTS (SELECT id FROM stream WHERE stream_name = $1)`

	var exists bool
	db, _ := ConnectDB()
	err := db.QueryRow(query, streamName).Scan(&exists)
	if err != nil {
		fmt.Println(err)
	}
	return exists
}

// ValidateStreamKey Validates stream by checking stream name against stream key.
func ValidateStreamKey(streamName, streamKey string) bool {
	const query = `SELECT EXISTS (SELECT id FROM stream WHERE stream_name = $1 AND stream_key = $2)`

	var valid bool
	db, _ := ConnectDB()
	err := db.QueryRow(query, streamName, streamKey).Scan(&valid)
	if err != nil {
		fmt.Println(err)
	}
	return valid
}
