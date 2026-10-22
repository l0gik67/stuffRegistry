package httpRouting

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"stuffRegistry/dto"
	"stuffRegistry/postgres"

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
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id < 1 {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}
	emp, err := empHandler.empService.DeleteEmployee(id)
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
	err := json.NewDecoder(r.Body).Decode(&empDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	if empDto.FullName == "" {
		return errors.New("full name is required")
	}
	if empDto.Position == "" {
		return errors.New("position is required")
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
