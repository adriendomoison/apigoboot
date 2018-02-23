package repo

import (
	"database/sql"
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/go-boot-api/database/dbconn"
	"github.com/adriendomoison/go-boot-api/oauth/repo/model"
)

// Make sure the interface is implemented correctly
var _ osin.Storage = (*Storage)(nil)

// Storage implements interface "github.com/RangelReale/osin".Storage and interface "github.com/ory/osin-storage".Storage
type Storage struct {
	db *sql.DB
}

// New returns a new postgres storage instance.
func New(db *sql.DB) *Storage {
	dbconn.DB.AutoMigrate(&model.Client{})
	dbconn.DB.AutoMigrate(&model.Authorize{})
	dbconn.DB.AutoMigrate(&model.Access{})
	dbconn.DB.AutoMigrate(&model.Refresh{})
	return &Storage{db}
}

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (s *Storage) Clone() osin.Storage {
	return s
}

// Close the resources the Storage potentially holds (using Clone for example)
func (s *Storage) Close() {
}