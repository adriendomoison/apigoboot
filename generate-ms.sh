DOCKERFILE="FROM golang:1.10

MAINTAINER Adrien Domoison \"adomoison@gmail.com\"

WORKDIR /go/src/github.com/adriendomoison/apigoboot/$1-micro-service

CMD go build component/serve-micro-service.go && ./serve-micro-service

EXPOSE 4200"

##

SERVEMSFILE="// Package main
package main

import (
    \"github.com/adriendomoison/apigoboot/$1-micro-service/component/$1\"
    \"github.com/adriendomoison/apigoboot/$1-micro-service/component/$1/repo\"
    \"github.com/adriendomoison/apigoboot/$1-micro-service/component/$1/rest\"
    \"github.com/adriendomoison/apigoboot/$1-micro-service/component/$1/service\"
    \"github.com/adriendomoison/apigoboot/$1-micro-service/config\"
    \"github.com/adriendomoison/apigoboot/$1-micro-service/database/dbconn\"
    \"github.com/gin-contrib/cors\"
    \"github.com/gin-gonic/gin\"
    \"log\"
)

// startAPI start the API and keep it alive
func main() {
    // Init DB and plan to close it at the end of the programme
    dbconn.Connect()
    defer dbconn.DB.Close()

    // Set GIN in production mode if run in production
    if !config.GDevEnv {
        gin.SetMode(gin.ReleaseMode)
    }

    // Init router
    router := gin.Default()
    router.Use(cors.New(getCORSConfig()))

    // $1 component
    $1Component := $1.New(rest.New(service.New(repo.New())))
    $1Component.AttachPublicAPI(router.Group(\"/api/v1\"))
    $1Component.AttachPrivateAPI(router.Group(\"/api/private-v1\"))

    // Start router
    go log.Println(\"Service $1 started: Navigate to \" + config.GAppUrl)
    router.Run(\":\" + config.GPort)
}

// getCORSConfig Generate CORS config for router
func getCORSConfig() cors.Config {
    CORSConfig := cors.DefaultConfig()
    CORSConfig.AllowCredentials = true
    CORSConfig.AllowAllOrigins = true
    CORSConfig.AllowHeaders = []string{\"*\"}
    CORSConfig.AllowMethods = []string{\"GET\", \"PUT\", \"POST\", \"DELETE\", \"OPTIONS\"}
    return CORSConfig
}"

##

FILE="// Package $1 is a component to manage the $1
package $1

import (
    \"github.com/gin-gonic/gin\"
)

// RestInterface is the model for the rest package of $1
type RestInterface interface {
    Post(c *gin.Context)
    Get(c *gin.Context)
    Put(c *gin.Context)
    Delete(c *gin.Context)
}

// Component implement interface component
type Component struct {
    rest RestInterface
}

// New return a new micro service instance
func New(rest RestInterface) *Component {
    return &Component{rest}
}

// AttachPublicAPI add the $1 micro-service public api with its dependencies
func (component *Component) AttachPublicAPI(group *gin.RouterGroup) {
    group.POST(\"/$1s\", component.rest.Post)
    group.GET(\"/$1s/:publicId\", component.rest.Get)
    group.PUT(\"/$1s/:publicId\", component.rest.Put)
    group.DELETE(\"/$1s/:publicId\", component.rest.Delete)
}

// AttachPrivateAPI add the $1 micro-service $1 api with its dependencies
func (component *Component) AttachPrivateAPI(group *gin.RouterGroup) {

}"

##

REST="// Package rest implement the callback required by the $1 package
package rest

import (
    \"github.com/adriendomoison/apigoboot/api-tool/errorhandling/apihelper\"
    \"github.com/adriendomoison/apigoboot/api-tool/errorhandling/servicehelper\"
    \"github.com/adriendomoison/apigoboot/$1-micro-service/component/$1\"
    \"github.com/gin-gonic/gin\"
    \"net/http\"
)

// ServiceInterface is the model for the service package of $1
type ServiceInterface interface {
    GetResourceOwnerId(email string) uint
    Add(creation RequestDTO) (ResponseDTO, *servicehelper.Error)
    Retrieve(string) (ResponseDTO, *servicehelper.Error)
    Edit(RequestDTO) (ResponseDTO, *servicehelper.Error)
    Remove(string) *servicehelper.Error
    IsThatTheUserId(string, uint) (bool, *servicehelper.Error)
}

// RequestDTO is the object to map JSON request body
type RequestDTO struct {
    PublicId          string \`json:\"$1Id\" binding:\"required,min=16\"\`
}

// ResponseDTO is the object to map JSON response body
type ResponseDTO struct {
    PublicId          string \`json:\"$1Id\"\`
}

// Make sure the interface is implemented correctly
var _ $1.RestInterface = (*rest)(nil)

// Implement interface
type rest struct {
    service ServiceInterface
}

// New return a new rest instance
func New(service ServiceInterface) *rest {
    return &rest{service}
}

// Post allows to access the service to create a $1
func (r *rest) Post(c *gin.Context) {
    var reqDTO RequestDTO
    if err := c.BindJSON(&reqDTO); err != nil {
        c.JSON(apihelper.BuildRequestError(err))
    } else {
        if resDTO, err := r.service.Add(reqDTO); err != nil {
            c.JSON(apihelper.BuildResponseError(err))
        } else {
             c.JSON(http.StatusCreated, resDTO)
        }
    }
}

// Get allows to access the service to retrieve a $1 when sending its $1 public id
func (r *rest) Get(c *gin.Context) {
    if resDTO, err := r.service.Retrieve(c.Param(\"$1Id\")); err != nil {
        c.JSON(apihelper.BuildResponseError(err))
    } else {
        c.JSON(http.StatusOK, resDTO)
    }
}

// Put allows to access the service to update the properties of a $1
func (r *rest) Put(c *gin.Context) {
    var reqDTO RequestDTO
    if err := c.BindJSON(&reqDTO); err != nil {
        c.JSON(apihelper.BuildRequestError(err))
    } else {
        if resDTO, err := r.service.Edit(reqDTO); err != nil {
            c.JSON(apihelper.BuildResponseError(err))
        } else {
            c.JSON(http.StatusOK, resDTO)
        }
    }
}

// Delete allows to access the service to remove a $1 from the records
func (r *rest) Delete(c *gin.Context) {
    if err := r.service.Remove(c.Param(\"$1Id\")); err != nil {
        c.JSON(apihelper.BuildResponseError(err))
    } else {
        c.JSON(http.StatusOK, gin.H{\"message\": \"$1 has been deleted successfully\"})
    }
}"

##

SERVICE="// Package service implement the services required by the rest package
package service

import (
	\"errors\"
	\"github.com/adriendomoison/apigoboot/api-tool/gentool\"
	\"github.com/adriendomoison/apigoboot/api-tool/errorhandling/servicehelper\"
	\"github.com/adriendomoison/apigoboot/$1-micro-service/component/$1/rest\"
	\"github.com/jinzhu/copier\"
	\"github.com/jinzhu/gorm\"
)

// RepoInterface is the model for the repo package of $1
type RepoInterface interface {
	Create($1 Entity) bool
	FindByID(id uint) ($1 Entity, err error)
	FindByPublicId(publicId string) ($1 Entity, err error)
	Update($1 Entity) error
	Delete($1 Entity) error
}

// Entity is the model of a $1 in the database
type Entity struct {
	gorm.Model
	PublicId          string \`gorm:\"UNIQUE;NOT NULL\"\`
}

// TableName allow to gives a specific name to the $1 table
func (Entity) TableName() string {
	return \"$1\"
}

// Make sure the interface is implemented correctly
var _ rest.ServiceInterface = (*service)(nil)

// service implement interface
type service struct {
	repo RepoInterface
}

// New return a new service instance
func New(repo RepoInterface) *service {
	return &service{repo}
}

// GetResourceOwnerId ask database to retrieve a user ID from its $1 public id
func (s *service) GetResourceOwnerId(publicId string) (userId uint) {
	entity, err := s.repo.FindByPublicId(publicId)
	if err != nil {
		return 0
	}
	return entity.UserID
}

// createDTOFromEntity copy all data from an entity to a Response DTO
func createDTOFromEntity(entity Entity) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	copier.Copy(&resDTO, &entity)
	return resDTO, error
}

// createEntityFromDTO copy all data from a Request DTO to an entity and initialize entity with some initialization values
func createEntityFromDTO(reqDTO rest.RequestDTO, init bool) (entity Entity, error *servicehelper.Error) {
	copier.Copy(&entity, &reqDTO)
	if init {
		entity.PublicId = gentool.GenerateRandomString(8)
	}
	return
}

// Add set up and create a $1
func (s *service) Add(reqDTO rest.RequestDTO) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := createEntityFromDTO(reqDTO, true)
	if err != nil {
		return rest.ResponseDTO{}, err
	} else if s.repo.Create(entity) {
		return createDTOFromEntity(entity)
	}
	return rest.ResponseDTO{}, &servicehelper.Error{Detail: errors.New(\"could not be created\"), Code: servicehelper.AlreadyExist}
}

// Retrieve ask database to retrieve a $1 from its public_id
func (s *service) Retrieve(publicId string) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(publicId)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{Detail: errors.New(\"no result found\"), Code: servicehelper.NotFound}
	}
	return createDTOFromEntity(entity)
}

// Edit edit user $1 and ask database to save changes
func (s *service) Edit(reqDTO rest.RequestDTO) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(reqDTO.PublicId)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{Detail: errors.New(\"no result found\"), Code: servicehelper.NotFound}
	}

    // TODO add fields to copy
	//entity.Example = reqDTO.Example

	if err := s.repo.Update(entity); err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New(\"could not update $1\"),
			Message: \"We could not update the $1, please contact us or try again later\",
			Code:    servicehelper.UnexpectedError}
	}
	return createDTOFromEntity(entity)
}

// Remove find a $1 in the database and delete it
func (s *service) Remove(publicId string) (error *servicehelper.Error) {
	if entity, err := s.repo.FindByPublicId(publicId); err != nil {
		return &servicehelper.Error{
			Detail:  errors.New(\"could not find $1\"),
			Message: \"We could not find any $1, please check the provided public_id\",
			Param:   \"public_id\",
			Code:    servicehelper.BadRequest,
		}
	} else if err := s.repo.Delete(entity); err != nil {
		return &servicehelper.Error{
			Detail:  errors.New(\"failed to delete $1\"),
			Message: \"We could not delete the $1, please try again later\",
			Code:    servicehelper.UnexpectedError,
		}
	}
	return
}

// IsThatTheUserId check if userIdToCheck is the same than the resource
func (s *service) IsThatTheUserId(publicId string, userIdToCheck uint) (same bool, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(publicId)
	if err != nil {
		return false, &servicehelper.Error{
			Detail:  errors.New(\"could not find $1\"),
			Message: \"We could not find any $1, please check the provided $1 public id\",
			Param:   \"public_id\",
			Code:    servicehelper.BadRequest,
		}
	}
	return entity.ID == userIdToCheck, nil
}"

##

REPO="// Package repo implement the function that contact the db required by the service package
package repo

import (
	\"github.com/adriendomoison/apigoboot/$1-micro-service/component/$1/service\"
	\"github.com/adriendomoison/apigoboot/$1-micro-service/database/dbconn\"
)

// Make sure the interface is implemented correctly
var _ service.RepoInterface = (*repo)(nil)

// Implement interface
type repo struct {
	repo service.RepoInterface
}

// New return a new repo instance
func New() *repo {
	dbconn.DB.AutoMigrate(&service.Entity{})
	return &repo{}
}

// Create create $1 in Database
func (crud *repo) Create($1 service.Entity) bool {
	if dbconn.DB.NewRecord($1) {
		dbconn.DB.Create(&$1)
	}
	return !dbconn.DB.NewRecord($1)
}

// FindByID find $1 in Database by ID
func (crud *repo) FindByID(id uint) ($1 service.Entity, err error) {
	if err = dbconn.DB.Where(\"id = ?\", id).First(&$1).Error; err != nil {
		return service.Entity{}, err
	}
	return $1, nil
}

// FindByPublicId find $1 in Database by public_id
func (crud *repo) FindByPublicId(publicId string) ($1 service.Entity, err error) {
	if err = dbconn.DB.Where(\"public_id = ?\", publicId).First(&$1).Error; err != nil {
		return service.Entity{}, err
	}
	return $1, nil
}

// Update edit $1 in Database
func (crud *repo) Update($1 service.Entity) error {
	return dbconn.DB.Save(&$1).Error
}

// Delete remove $1 from Database
func (crud *repo) Delete($1 service.Entity) error {
	return dbconn.DB.Delete(&$1).Error
}"

##

CONFIG="// Package config generate the environment of the API
package config

import (
	\"log\"
	\"os\"
)

// GAppName define the app name
var GAppName = \"apigoboot\"

var devPort = \"4200\"
var devAppUrl = \"http://api.go.boot\"
var prodAppUrl = \"https://apigoboot.herokuapp.com\"

// GDevEnv define if environment is in dev mode
var GDevEnv bool

// GUnitTestingEnv define if environment is in testing mode
var GUnitTestingEnv bool

// GPort is the application current port
var GPort string

// GAppUrl is the application url
var GAppUrl string

// init initialize the default environment
func init() {
	GPort = os.Getenv(\"PORT\")
	if GPort == \"\" {
		GDevEnv = true
		GPort = devPort
		GAppUrl = devAppUrl + \":\" + GPort
		log.Println(\"Dev Environnement detected\")
	} else {
		GDevEnv = false
		GAppUrl = prodAppUrl
		log.Println(\"Heroku Environement detected\")
	}
}

// SetToTestingEnv set the test environment, this need to be called before testing to prevent the development database to be used
func SetToTestingEnv() {
	GDevEnv = false
	GUnitTestingEnv = true
}"

##

DBCONN="// Package dbconn allows connection to db
package dbconn

import (
	\"github.com/adriendomoison/apigoboot/$1-micro-service/config\"
	\"github.com/jinzhu/gorm\"
	// Add postgres for gorm
	_ \"github.com/jinzhu/gorm/dialects/postgres\"
	\"log\"
	\"os\"
	\"syscall\"
	\"time\"
)

// DB object for repo
var DB *gorm.DB

// Connect connect to database depending of the env
func Connect() (err error) {
	if config.GUnitTestingEnv {
		err = connectToDB(config.GAppName+\"_test\", config.GAppName+\"_test\", config.GAppName+\"_test\", \"localhost\")
	} else if _, ok := syscall.Getenv(\"DYNO\"); ok {
		err = connectToDB(os.Getenv(\"DB_USER\"), os.Getenv(\"DB_NAME\"), os.Getenv(\"DB_PASSWORD\"), os.Getenv(\"DB_HOST\"))
	} else if config.GDevEnv {
		err = connectToDB(config.GAppName, config.GAppName, config.GAppName, \"db\")
	}
	if err != nil {
		log.Panic(\"Database status: [Failed to connect]\", err)
	}
	return
}

// connectToDB do the connection request to the database depending on provided parameters
func connectToDB(username string, dbName string, password string, host string) (err error) {
	log.Println(\"CONNECTING TO [\" + dbName + \"] DB...\")
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(\"postgres\", \"host=\"+host+\" user=\"+username+\" dbname=\"+dbName+\" sslmode=disable password=\"+password)
		if err != nil {
			log.Println(\"Still trying...\")
		} else {
			DB.SingularTable(true)
			log.Println(\"Database status: [Connected]\")
			break
		}
		time.Sleep(5 * time.Second)
	}
	return
}"

mkdir $1-micro-service
cd $1-micro-service
mkdir -p component config database/dbconn
echo "$DOCKERFILE" >> Dockerfile
echo "$CONFIG" >> config/config.go
echo "$DBCONN" >> database/dbconn/dbconn.go
cd component
echo "$SERVEMSFILE" >> serve-micro-service.go
mkdir $1
cd $1
echo "$FILE" >> $1.go
mkdir repo rest service
echo "$REST" >> rest/rest.go
echo "$SERVICE" >> service/service.go
echo "$REPO" >> repo/repo.go
cd ../..
git add .