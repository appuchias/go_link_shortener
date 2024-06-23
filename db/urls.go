package db

import (
	"database/sql"
	"log"
)

type URL struct {
	IDLink   int
	Src      string
	Dst      string
	IsCustom bool
}

func GetURLRedirect(src string) (string, error) {
	var dst string
	err := db.QueryRow("SELECT dst FROM Link WHERE src = ?", src).Scan(&dst)
	if err == sql.ErrNoRows {
		// log.Println("No rows found")
		return "", nil
	} else if err != nil {
		log.Println("Error getting URL redirect from DB:", err)
		return "", err
	}
	return dst, nil
}

func GetUserURLs(username string) ([]URL, error) {
	rows, err := db.Query("SELECT id_link, src, dst, is_custom FROM Link WHERE owner = (SELECT id_user FROM User WHERE username = ?)", username)
	if err != nil {
		log.Println("Error getting URLs from DB:", err)
		return nil, err
	}
	defer rows.Close()

	urls := []URL{}
	for rows.Next() {
		var src, dst string
		var id_link int
		var is_custom bool
		if err := rows.Scan(&id_link, &src, &dst, &is_custom); err != nil {
			log.Println("Error scanning URLs from DB:", err)
			return nil, err
		}
		urls = append(urls, URL{IDLink: id_link, Src: src, Dst: dst, IsCustom: is_custom})
	}

	return urls, nil
}

func GetURLDetails(id_link int) (URL, error) {
	var src, dst string
	var isCustom bool
	err := db.QueryRow("SELECT src, dst, is_custom FROM Link WHERE id_link = ?", id_link).Scan(&src, &dst, &isCustom)
	if err != nil {
		log.Println("Error getting URL details from DB:", err)
		return URL{}, err
	}
	return URL{IDLink: id_link, Src: src, Dst: dst, IsCustom: isCustom}, nil
}

func AddURL(owner int, src string, dst string, isCustom bool) error {
	_, err := db.Exec("INSERT INTO Link (owner, src, dst, is_custom) VALUES (?, ?, ?, ?)", owner, src, dst, isCustom)
	if err != nil {
		log.Println("Error adding URL to DB:", err)
		return err
	}
	return nil
}
