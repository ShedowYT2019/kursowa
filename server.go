package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Структура даних пристрою згідно з завданням
type DeviceLog struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	IPAddress   string    `json:"ip_address"`
	RoutingType string    `json:"routing_type"`
	Timestamp   time.Time `json:"timestamp"`
}

var (
	logFile = "network_devices.log"
	mu      sync.Mutex
)

func handleRequests(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		// Читання та повернення лог-файлу в термінал
		data, err := ioutil.ReadFile(logFile)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Лог-файл порожній або не знайдений.")
			return
		}
		fmt.Fprintf(w, "--- Лог входження пристроїв (Варіант 14) ---\n%s", string(data))

	case http.MethodPost:
		// Обробка POST-запиту
		var device DeviceLog
		if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		device.Timestamp = time.Now()

		logEntry := fmt.Sprintf("[%s] Пристрій: %s | Тип: %s | IP: %s | Маршрутизація: %s\n",
			device.Timestamp.Format("2006-01-02 15:04:05"),
			device.Name, device.Type, device.IPAddress, device.RoutingType)

		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.WriteString(logEntry)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Запис успішно додано для пристрою: %s", device.Name)

	case http.MethodDelete:
		// Очищення лог-файлу
		err := os.Remove(logFile)
		if err != nil {
			fmt.Fprintf(w, "Лог-файл вже порожній.")
		} else {
			fmt.Fprintf(w, "Лог-файл успішно очищено.")
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", handleRequests)
	fmt.Println("Сервер групи 74 запущено на http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
