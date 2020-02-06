// Package conf contains functions to configure logging, the database and router
package conf

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/epip-io/go-birthday-calendar/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// BirthdayService reponses with a API version and name
func BirthdayService() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{
			"Name":        "Birthday",
			"Description": "A birthday reminder and wishes service",
			"Version":     "v1alpha1",
		})
	})
}

// BirthdayMessage reponse with a reminder or wishes happy birthday
func BirthdayMessage(db *gorm.DB, logger *log.Entry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]string
		var c int

		vars := mux.Vars(r)
		person := models.Person{
			Name: vars["Name"],
		}

		logger.WithFields(log.Fields{
			"Name": person.Name,
		}).Debug("getting birthday")

		db.Where(&person).Find(&person).Count(&c)

		if c == 0 {
			respondWithError(w, http.StatusNotFound, "Name not found")
			return
		}

		by, bm, bd := time.Time(person.BirthDate).Date()
		ty, tm, td := time.Now().Date()

		switch {
		case bd < td && bm == tm, bm < tm:
			by = ty + 1
		default:
			by = ty

		}

		db := time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
		dt := time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)

		dd := int(math.Floor(db.Sub(dt).Hours() / 24))

		switch dd {
		case 0:
			resp = map[string]string{
				"message": fmt.Sprintf("Hi, %s! Happy birthday!", person.Name),
			}
		case 5:
			resp = map[string]string{
				"message": fmt.Sprintf("Hi, %s! Your birtday is in %d days!", person.Name, dd),
			}
		default:
			resp = map[string]string{
				"message": fmt.Sprintf("Hi, %s!", person.Name),
			}
		}

		respondWithJSON(w, http.StatusOK, resp)
	})
}

// PersonsBirthday adds / updates a person's birthday
func PersonsBirthday(db *gorm.DB, logger *log.Entry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		person := models.Person{
			Name: vars["Name"],
		}

		// Finding existing person
		db.Where(&person).First(&person)

		// Adding birthday
		if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create / Update person
		logger.WithFields(log.Fields{
			"Name":      person.Name,
			"BirthDate": person.BirthDate,
		}).Debug("saving person")
		db.Save(&person)

		w.WriteHeader(http.StatusNoContent)
	})
}

// Health response with the health of the service
func Health(db *gorm.DB, logger *log.Entry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !db.HasTable(&models.Person{}) {
			respondWithError(w, http.StatusServiceUnavailable, "persons table doesn't exist in database yet")
			return
		}

		var count int
		var people []models.Person

		db.Find(&people).Count(&count)

		logger.WithFields(log.Fields{
			"People": people,
			"Count":  count,
		}).Debugf("health check")

		respondWithJSON(w, http.StatusOK, map[string]string{
			"count": fmt.Sprintf("%d", count),
		})
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
