package repo

import (
	"database/sql"
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/gobootapi/database/dbconn"
	"github.com/adriendomoison/gobootapi/oauth/repo/dbmodel"
)

// Make sure the interface is implemented correctly
var _ osin.Storage = (*Storage)(nil)

// Storage implements interface "github.com/RangelReale/osin".Storage and interface "github.com/ory/osin-storage".Storage
type Storage struct {
	db *sql.DB
}

// New returns a new postgres storage instance.
func New(db *sql.DB) *Storage {
	dbconn.DB.AutoMigrate(&dbmodel.Client{})
	dbconn.DB.AutoMigrate(&dbmodel.Authorize{})
	dbconn.DB.AutoMigrate(&dbmodel.Access{})
	dbconn.DB.AutoMigrate(&dbmodel.Refresh{})
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