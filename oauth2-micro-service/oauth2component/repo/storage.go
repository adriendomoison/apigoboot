package repo

import (
	"database/sql"
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/gobootapi/oauth2-micro-service/database/dbconn"
	"github.com/adriendomoison/gobootapi/oauth2-micro-service/oauth2component/service"
)

// Make sure the interface is implemented correctly
var _ osin.Storage = (*Storage)(nil)

// Storage implements interface "github.com/RangelReale/osin".Storage and interface "github.com/ory/osin-storage".Storage
type Storage struct {
	db *sql.DB
}

// NewStorage returns a new postgres storage instance.
func NewStorage(db *sql.DB) *Storage {
	dbconn.DB.AutoMigrate(&service.Client{})
	dbconn.DB.AutoMigrate(&service.Authorize{})
	dbconn.DB.AutoMigrate(&service.Access{})
	dbconn.DB.AutoMigrate(&service.Refresh{})
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