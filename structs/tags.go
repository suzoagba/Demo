package structs

import (
	"database/sql"
	"log"
)

var Tags []Tag

func GetTags(db *sql.DB) {
	tags, err := getTagsFromDatabase(db)
	if err != nil {
		log.Fatal(err)
	}
	Tags = tags
}

func getTagsFromDatabase(db *sql.DB) ([]Tag, error) {
	query := "SELECT id, name FROM tags"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
