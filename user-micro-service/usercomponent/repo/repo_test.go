package repo

import (
	"os"
	"testing"
	"github.com/adriendomoison/apigoboot/user-micro-service/config"
	"github.com/adriendomoison/apigoboot/user-micro-service/database/dbconn"
	"github.com/adriendomoison/apigoboot/user-micro-service/usercomponent/service"
)

var r *repo

func TestMain(m *testing.M) {
	config.SetToTestingEnv()
	dbconn.Connect()
	defer dbconn.DB.Close()

	r = New()

	code := m.Run()

	dbconn.DB.DropTable(&service.Entity{})

	os.Exit(code)
}

func TestRepository_Create(t *testing.T) {
	if !r.Create(service.Entity{
		Email:    "john@example.dev",
		Username: "John",
		Password: "QNDNwefwf44DfY@wDNwfEC#H4$$fNEC4H4WEw&@w4NFw$wHwf4WEwfFwSsf@As$Dsdfsdf$JsFHIWE",
	}) {
		t.Error("Could not create perfectly formed user")
	}

	if r.Create(service.Entity{
		Email:    "john@example.dev",
		Username: "John",
		Password: "QNDNwefwf44DfY@wDNwfEC#H4$$fNEC4H4WEw&@w4NFw$wHwf4WEwfFwSsf@As$Dsdfsdf$JsFHIWE",
	}) {
		t.Error("The same user was created twice")
	}
}

func TestRepository_FindByID(t *testing.T) {
	entity, err := r.FindByID(1)

	if err != nil {
		t.Error(err)
	}

	if entity.Email != "john@example.dev" || entity.Username != "John" || entity.Password != "QNDNwefwf44DfY@wDNwfEC#H4$$fNEC4H4WEw&@w4NFw$wHwf4WEwfFwSsf@As$Dsdfsdf$JsFHIWE" {
		t.Error("Could not create perfectly formde user")
	}
}
func TestRepository_FindByEmail(t *testing.T) {
	entity, err := r.FindByEmail("john@example.dev")

	if err != nil {
		t.Error(err)
	}

	if entity.Email != "john@example.dev" || entity.Username != "John" || entity.Password != "QNDNwefwf44DfY@wDNwfEC#H4$$fNEC4H4WEw&@w4NFw$wHwf4WEwfFwSsf@As$Dsdfsdf$JsFHIWE" {
		t.Error("Could not create perfectly formde user")
	}
}

var entityToDelete service.Entity

func TestRepository_Delete(t *testing.T) {
	entityToDelete, _ = r.FindByEmail("john@example.dev")

	if err := r.Delete(entityToDelete); err != nil {
		t.Error(err)
	}
}

func TestRepository_DeleteUserThatDoesNotExist(t *testing.T) {
	if err := r.Delete(entityToDelete); err != nil {
		t.Error(err)
	}
}
