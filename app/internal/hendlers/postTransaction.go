package hendlers

import (
	"bytes"
	"cachManagerApp/app/db/models"
	"encoding/json"
	"log"
	"net/http"
)

func PostTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := new(models.Transactions)
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(transaction)

	resultDB := database.DB.Db.Create(&transaction)
	if resultDB.Error != nil {
		http.Error(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
