package models

import (
	"database/sql"
	"errors"
	
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
	UpdatedAt time.Time
}
type SnippetModelInterface interface {
	Insert(title string, content string, expires int, authorName string) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
	UpdateSnippet(id int, title string, content string, expires int, authorName string) (int, error)
	DeleteSnippet(id int) (int, error) 
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int, authorName string) (int, error) {

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
func (m *SnippetModel) Get(id int) (Snippet, error) {

	stmt := `SELECT id, title, content, created, expires, updated_at FROM snippets
WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	// Initialize a new zeroed Snippet struct.
	var s Snippet
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.UpdatedAt)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	// If everything went OK, then return the filled Snippet struct.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {

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
func (m *SnippetModel) DeleteSnippet(id int) (int, error) {

	stmt := `DELETE FROM snippets WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {	
		return 0, err
	}

	return 0, nil
}

func (m *SnippetModel) UpdateSnippet(id int, title string, content string, expires int, authorName string) (int, error) {

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