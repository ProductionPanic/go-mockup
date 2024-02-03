package crawl

import (
	"database/sql"
	"mock/url"
)

type CrawlDB struct {
	sqliteDbFile string
	db           *sql.DB
}

const sqliteDbFile = "crawl.db"

func DB() *CrawlDB {
	db := &CrawlDB{
		sqliteDbFile: sqliteDbFile,
	}
	err := db.Open()
	if err != nil {
		return nil

	}
	err = db.CreateSchema()
	if err != nil {
		return nil
	}
	return db
}

func (c *CrawlDB) Open() error {
	db, err := sql.Open("sqlite3", c.sqliteDbFile)
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

func (c *CrawlDB) Close() error {
	return c.db.Close()
}

func (c *CrawlDB) CreateSchema() error {
	_, err := c.db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			depth INTEGER NOT NULL,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    host TEXT NOT NULL
		);

	    		CREATE TABLE IF NOT EXISTS links (
	    		    			id INTEGER PRIMARY KEY AUTOINCREMENT,
	    		    			from_url_id INTEGER NOT NULL,
	    		    			to_url_id INTEGER NOT NULL,
	    		    			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    		    			FOREIGN KEY(from_url_id) REFERENCES urls(id),
	    		    			FOREIGN KEY(to_url_id) REFERENCES urls(id)
	    		                                 	    		    			);

	`)
	return err
}

func (c *CrawlDB) InsertURL(urlstr string, depth int) error {
	urlObj := url.NewURL(urlstr)
	_, err := c.db.Exec(`INSERT INTO urls (url, depth, host) VALUES (?, ?, ?)`, urlstr, depth, urlObj.Host())
	return err
}

func (c *CrawlDB) InsertLink(fromURL string, toURL string) error {
	_, err := c.db.Exec(`INSERT INTO links (from_url_id, to_url_id) VALUES ((SELECT id FROM urls WHERE url = ?), (SELECT id FROM urls WHERE url = ?))`, fromURL, toURL)
	return err
}

func (c *CrawlDB) GetURLs() ([]string, error) {
	rows, err := c.db.Query(`SELECT url FROM urls`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	urls := []string{}
	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (c *CrawlDB) GetByHost(host string) ([]string, error) {
	rows, err := c.db.Query(`SELECT url FROM urls WHERE host = ?`, host)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	urls := []string{}
	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
