package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/config"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/postgres"
	repository "github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/postgres/sqlc"
)

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LastNames string `json:"lastNames"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Status    string `json:"status"`
	BirthDate string `json:"birthDate"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	RoleID    string `json:"roleId"`
}

func main() {
	cfg := config.Load()
	db, err := postgres.New(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.New(db)
	ctx := context.Background()

	rolesFile, err := os.ReadFile("seeders/roles-data.json")
	if err != nil {
		log.Fatal(err)
	}

	var roles []Role
	if err := json.Unmarshal(rolesFile, &roles); err != nil {
		log.Fatal(err)
	}

	for _, r := range roles {
		_, err := db.ExecContext(ctx, "INSERT INTO role (id, name, description) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", r.ID, r.Name, r.Description)
		if err != nil {
			log.Printf("Failed to seed role %s: %v", r.Name, err)
		}
	}
	log.Println("Roles seeded")

	usersFile, err := os.ReadFile("seeders/users-data.json")
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	if err := json.Unmarshal(usersFile, &users); err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		birthDate, _ := time.Parse("1/2/2006", u.BirthDate)
		_, err := repo.CreateUser(ctx, repository.CreateUserParams{
			ID:        u.ID,
			Name:      u.Name,
			LastNames: u.LastNames,
			Email:     u.Email,
			Username:  u.Username,
			Password:  u.Password,
			Status:    u.Status,
			BirthDate: birthDate,
			Address:   u.Address,
			RoleID:    u.RoleID,
			Phone:     u.Phone,
		})
		if err != nil {
			log.Printf("Failed to seed user %s: %v", u.Username, err)
		}
	}
	log.Println("Users seeded")
}
