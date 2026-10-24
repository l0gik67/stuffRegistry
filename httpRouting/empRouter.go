package httpRouting

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"stuffRegistry/dto"
	"stuffRegistry/postgres"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()
	empHandler := NewEmployeeHandler()

	router.HandleFunc("/employees", empHandler.getEmps).Methods(http.MethodGet)
	router.HandleFunc("/employees", empHandler.newEmp).Methods(http.MethodPost)
	router.HandleFunc("/employees/{id}", empHandler.deleteEmp).Methods(http.MethodDelete)

	return router
}

func (empHandler *EmployeeHandler) deleteEmp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
	if err != nil || id < 1 {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}
	emp, err := empHandler.empService.DeleteEmployee(int(id))
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "employee not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func (empHandler *EmployeeHandler) getEmps(w http.ResponseWriter, r *http.Request) {
	emps, err := empHandler.empService.GetEmployees()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	if json.NewEncoder(w).Encode(emps) != nil {
		return
	}
}

func (empHandler EmployeeHandler) newEmp(w http.ResponseWriter, r *http.Request) {
	empDto := dto.EmployeeDto{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&empDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "request body must contain a single JSON value", http.StatusBadRequest)
		return
	}
	if err := validateEmpDto(empDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	emp, err := empHandler.empService.NewEmp(empDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func validateEmpDto(empDto dto.EmployeeDto) error {
	if strings.TrimSpace(empDto.FullName) == "" {
		return errors.New("full name is required")
	}
	if utf8.RuneCountInString(empDto.FullName) > 255 {
		return errors.New("full name must not exceed 255 characters")
	}
	if strings.TrimSpace(empDto.Position) == "" {
		return errors.New("position is required")
	}
	if utf8.RuneCountInString(empDto.Position) > 255 {
		return errors.New("position must not exceed 255 characters")
	}
	return nil
}

type EmployeeHandler struct {
	empService *postgres.EmployeeService
}

func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{
		empService: postgres.NewEmployeeService(),
	}
}
