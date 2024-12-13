package models

import (
	"database/sql"
	"errors"

	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Artikel struct {
	ID        int        `gorm:"column:id" json:"artikel_id"`
	Title     string     `gorm:"column:title" json:"title"`
	Kontent   string     `gorm:"column:kontent" json:"kontent"`
	Komentar  []Komentar `gorm:"Foreignkey:Artikel_ID;association_foreignkey:ID;" json:"komentar"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

type Komentar struct {
	ID         int        `gorm:"column:id" json:"komentar_id"`
	Artikel_ID string     `gorm:"column:artikel_id" json:"artikel_id"`
	Komentar   string     `gorm:"column:komentar" json:"komentar"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"-"`
}

type ArtikelModelInterface interface {
	InsertArtikel(title string, content string, expires int, authorName string) (int, error)
	GetArtikel(id int) (Artikel, error)
	LatestArtikel() ([]Artikel, error)
	UpdateArtikel(id int, title string, content string, expires int, authorName string) (int, error)
	DeleteArtikel(id int) (int, error)
}

// Define a ArtikelModel type which wraps a sql.DB connection pool.
type ArtikelModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) InsertArtikel(title string, content string, expires int, authorName string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires, author_name)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`

	result, err := m.DB.Exec(stmt, title, content, expires, authorName)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId() //LastInsertId() get the ID of our  newly inserted record in the snippets table.
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *ArtikelModel) GetArtikel(id int) (Artikel, error) {


	stmt := `SELECT id, title, kontent, created, expires, updated_at FROM snippets
WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	// Initialize a new zeroed Snippet struct.
	var s Artikel
	err := row.Scan(&s.ID, &s.Title, &s.Kontent, &s.CreatedAt, &s.UpdatedAt, &s.UpdatedAt)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return Artikel{}, ErrNoRecord
		} else {
			return Artikel{}, err
		}
	}
	// If everything went OK, then return the filled Snippet struct.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) LatestArtikel() ([]Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > UTC_TIMESTAMP() ORDER BY id ASC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the Snippet structs.
	var snippets []Snippet
	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a new zeroed Snippet struct.
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return snippets, nil
}

// this is for delete snippet row
func (m *SnippetModel) DeleteArtikel(id int) (int, error) {

	stmt := `DELETE FROM snippets WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (m *SnippetModel) UpdateArtikel(id int, title string, content string, expires int, authorName string) (int, error) {

	stmt := `UPDATE snippets SET title= ?, content= ?, expires= DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), 
	updated_at= UTC_TIMESTAMP(), author_name=?
	WHERE id = ?`

	_, err := m.DB.Exec(stmt, title, content, expires, authorName, id)
	if err != nil {
		return 0, err
	}
	//log.Println("UPDATE: Title: " + title + " | Content: " + content, expires, "Day", "id is", id)

	return id, nil
}
