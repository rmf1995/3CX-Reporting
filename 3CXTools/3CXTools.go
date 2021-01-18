package main

import (
	"bufio"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/spf13/viper"
)

// C config var from strut
var C config

var targetFailoverTimeInterval int

type config struct {
	Database database
	OpsGenie opsgenie
}

type database struct {
	User string `mapstructure:"User"`
	Pass string `mapstructure:"Pass"`
	Name string `mapstructure:"Name"`
	Host string `mapstructure:"Host"`
	Port string `mapstructure:"Port"`
}

type opsgenie struct {
	APIKey string `mapstructure:"APIKey"`
}

var servers struct {
	server []server `json:"data"`
}

type server struct {
	id                     int    `mapstructure:"ID"`
	Name                   string `mapstructure:"Name"`
	Location               string `mapstructure:"Location"`
	URL                    string `mapstructure:"URL"`
	UserName               string `mapstructure:"UserName"`
	PwStateID              string `mapstructure:"PwStateID"`
	lastUpdated            string `mapstructure:"lastUpdated"`
	SvrPasswordStateRecord PasswordStateRecord
	client                 *http.Client
	SvrSystemStatus        SystemStatus
	SvrAdditionalStatus    AdditionalStatus
	SvrUpdates             Updates
	SvrFailOver            FailOver
}

// FailOver represents the Updates response
type FailOver struct {
	ID           int `json:"Id"`
	ActiveObject struct {
		ID                int    `json:"Id"`
		Str               string `json:"_str"`
		IsNew             bool   `json:"IsNew"`
		LicenseRestricted bool   `json:"licenseRestricted"`
		Enabled           struct {
			Type  int  `json:"type"`
			Value bool `json:"_value"`
		} `json:"Enabled"`
		Mode struct {
			Type           int      `json:"type"`
			Selected       string   `json:"selected"`
			PossibleValues []string `json:"possibleValues"`
			Translatable   bool     `json:"translatable"`
		} `json:"Mode"`
		RemoteServer struct {
			Type     int    `json:"type"`
			Disabled bool   `json:"disabled"`
			Value    string `json:"_value"`
		} `json:"RemoteServer"`
		TestSIPServer struct {
			Type     int  `json:"type"`
			Disabled bool `json:"disabled"`
			Value    bool `json:"_value"`
		} `json:"TestSIPServer"`
		TestWebServer struct {
			Type     int  `json:"type"`
			Disabled bool `json:"disabled"`
			Value    bool `json:"_value"`
		} `json:"TestWebServer"`
		TestTunnel struct {
			Type     int  `json:"type"`
			Disabled bool `json:"disabled"`
			Value    bool `json:"_value"`
		} `json:"TestTunnel"`
		Interval struct {
			Type     int  `json:"type"`
			Disabled bool `json:"disabled"`
			Value    int  `json:"_value"`
		} `json:"Interval"`
		Condition struct {
			Type           int      `json:"type"`
			Disabled       bool     `json:"disabled"`
			Selected       string   `json:"selected"`
			PossibleValues []string `json:"possibleValues"`
			Translatable   bool     `json:"translatable"`
		} `json:"Condition"`
		PreStartScript struct {
			Type           int      `json:"type"`
			Disabled       bool     `json:"disabled"`
			Selected       string   `json:"selected"`
			PossibleValues []string `json:"possibleValues"`
			Translatable   bool     `json:"translatable"`
		} `json:"PreStartScript"`
		PostStartScript struct {
			Type           int      `json:"type"`
			Disabled       bool     `json:"disabled"`
			Selected       string   `json:"selected"`
			PossibleValues []string `json:"possibleValues"`
			Translatable   bool     `json:"translatable"`
		} `json:"PostStartScript"`
	} `json:"ActiveObject"`
}

// Updates represents the Updates response
type Updates struct {
	ID           int `json:"Id"`
	ActiveObject struct {
		ID          string `json:"Id"`
		Str         string `json:"_str"`
		IsNew       bool   `json:"IsNew"`
		ReadyToSave struct {
			Type  string `json:"type"`
			Hide  bool   `json:"hide"`
			Value string `json:"_value"`
		} `json:"ReadyToSave"`
		ScheduleType struct {
			Type           string   `json:"type"`
			Disabled       bool     `json:"disabled"`
			Selected       string   `json:"selected"`
			PossibleValues []string `json:"possibleValues"`
			Translatable   bool     `json:"translatable"`
		} `json:"ScheduleType"`
		ScheduleDay struct {
			Type           string   `json:"type"`
			Disabled       bool     `json:"disabled"`
			Selected       string   `json:"selected"`
			PossibleValues []string `json:"possibleValues"`
			Translatable   bool     `json:"translatable"`
		} `json:"ScheduleDay"`
		ScheduleTime struct {
			Type     string `json:"type"`
			Disabled bool   `json:"disabled"`
			Value    string `json:"_value"`
		} `json:"ScheduleTime"`
		TcxPbxUpdates struct {
			Type  string `json:"type"`
			Value bool   `json:"_value"`
		} `json:"TcxPbxUpdates"`
	} `json:"ActiveObject"`
}

// SystemStatus represents the SystemStatus response
type SystemStatus struct {
	FQDN                      string
	Version                   string
	Activated                 bool
	MaxSimCalls               int
	MaxSimMeetingParticipants int
	CallHistoryCount          int
	ChatMessagesCount         int
	ExtensionsRegistered      int
	OwnPush                   bool
	ExtensionsTotal           int
	TrunksRegistered          int
	TrunksTotal               int
	CallsActive               int
	BlacklistedIPCount        int
	MemoryUsage               int
	PhysicalMemoryUsage       int
	FreeFirtualMemory         int64
	TotalVirtualMemory        int64
	FreePhysicalMemory        int64
	TotalPhysicalMemory       int64
	DiskUsage                 int
	FreeDiskSpace             int64
	TotalDiskSpace            int64
	CPUUsage                  int
	MaintenanceExpiresAt      *time.Time
	Support                   bool
	ExpirationDate            interface{}
	OutboundRules             int
	BackupScheduled           bool
	LastBackupDateTime        *time.Time
	ResellerName              string
	LicenseKey                string
	ProductCode               string
}

// AdditionalStatus represents the AdditionalStatus response
type AdditionalStatus struct {
	RecordingUsedSpace    uint
	RecordingQuota        uint
	RecordingStopped      bool
	RecordingQuotaReached bool
}

// PasswordStateRecord represents the PasswordStateRecord Record
type PasswordStateRecord struct {
	PasswordID     int
	Title          string
	Domain         string
	HostName       string
	UserName       string
	Description    string
	GenericField1  string
	GenericField2  string
	GenericField3  string
	GenericField4  string
	GenericField5  string
	GenericField6  string
	GenericField7  string
	GenericField8  string
	GenericField9  string
	GenericField10 string
	AccountTypeID  int
	Notes          string
	URL            string
	Password       string
	ExpiryDate     string
	AllowExport    bool
	AccountType    string
}

func dbConn() (db *sql.DB) {

	dbDriver := "mysql"
	dbUser := C.Database.User
	dbPass := C.Database.Pass
	dbName := C.Database.Name
	dbHost := C.Database.Host
	dbPort := C.Database.Port

	connStr := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	db, err := sql.Open(dbDriver, connStr)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// ProductCodeToLicense receives ProductCode and gives the License
func ProductCodeToLicense(license string) string {
	switch license {
	case "3CXPSPROFENTSPLA":
		return "Enterprise Annual"
	case "3CXPSPROFSPLA":
		return "Professional Annual"
	case "3CXPSPROF":
		return "Professional Perpetual"
	case "3CXPSSPLA":
		return "Standard Annual"
	default:
		return license + " - Who Knows?"
	}

}

// GetPasswordFormPWState Retruns Password from passwordID
func GetPasswordFormPWState(pwid string) string {
	// Create a Resty Client
	client := resty.New()
	// Disable Security Check (HTTPS)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	db := dbConn()

	sqlStatement := "SELECT PasswordStateURL, PasswordStateAPIKey FROM config where id=1"
	var PasswordStateURL string
	var PasswordStateAPIKey string

	row := db.QueryRow(sqlStatement)
	errStatement := row.Scan(&PasswordStateURL, &PasswordStateAPIKey)

	if errStatement != nil {
		log.Println(errStatement)
	}

	// Login With Username And Password
	PasswordStateresp, PasswordStateErr := client.R().
		SetHeader("APIKey", PasswordStateAPIKey).
		Get(PasswordStateURL + "/api/passwords/" + pwid)
	if PasswordStateErr != nil {
		log.Println(PasswordStateErr)
	}

	//fmt.Println("  Body       :\n", PasswordStateresp)
	var PasswordRecord []PasswordStateRecord
	err := json.Unmarshal(PasswordStateresp.Body(), &PasswordRecord)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	return PasswordRecord[0].Password
}

// GetFailoverTimeIntervalData gets data
func GetFailoverTimeIntervalData(id int) {

	// Create a Resty Client
	client := resty.New()
	// Disable Security Check (HTTPS)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// Login With Username And Password
	LoginResp, LogiErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"` + servers.server[id].UserName + `", "password":"` + GetPasswordFormPWState(servers.server[id].PwStateID) + `"}`).
		Post("https://" + servers.server[id].URL + ":5001/api/login")

	if LogiErr != nil {
		log.Println(LogiErr)
	}

	LoginStatusCode := LoginResp.StatusCode()
	switch LoginStatusCode {
	case 200:

		FailOverResp, FailOverErr := client.R().
			EnableTrace().
			SetHeader("Content-Type", "application/json").
			Post("https://" + servers.server[id].URL + ":5001/api/BackupAndRestoreList/failover_settings")
		if FailOverErr != nil {
			log.Println(FailOverErr)
		}
		err := json.Unmarshal(FailOverResp.Body(), &servers.server[id].SvrFailOver)
		if err != nil {
			log.Println(err)
		}

	case 401:
		// session expired
		log.Printf("Login Session Expired code: %d", LoginStatusCode)
	default:
		log.Printf("Unexpected Login Status code: %d", LoginStatusCode)
	}

}

// SetFailoverTimeIntervalData gets data
func SetFailoverTimeIntervalData(id int) {

	// Create a Resty Client
	client := resty.New()
	// Disable Security Check (HTTPS)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// Login With Username And Password
	LoginResp, LogiErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"` + servers.server[id].UserName + `", "password":"` + GetPasswordFormPWState(servers.server[id].PwStateID) + `"}`).
		Post("https://" + servers.server[id].URL + ":5001/api/login")

	if LogiErr != nil {
		log.Println(LogiErr)
	}

	LoginStatusCode := LoginResp.StatusCode()
	switch LoginStatusCode {
	case 200:
		FailOverResp, FailOverErr := client.R().
			EnableTrace().
			SetHeader("Content-Type", "application/json").
			Post("https://" + servers.server[id].URL + ":5001/api/BackupAndRestoreList/failover_settings")
		if FailOverErr != nil {
			log.Println(FailOverErr)
		}
		err := json.Unmarshal(FailOverResp.Body(), &servers.server[id].SvrFailOver)
		if err != nil {
			log.Println(err)
		}

		_, FailOverUpdateErr := client.R().
			EnableTrace().
			SetBody(`{"Path":{"ObjectId":`+strconv.Itoa(servers.server[id].SvrFailOver.ID)+`,"PropertyPath":[{"Name":"Interval"}]},"PropertyValue":`+strconv.Itoa(targetFailoverTimeInterval)+`}`).
			SetHeader("Content-Type", "application/json").
			Post("https://" + servers.server[id].URL + ":5001/api/edit/update")
		if FailOverUpdateErr != nil {
			log.Println(FailOverUpdateErr)
		}

		_, FailOverSaveErr := client.R().
			EnableTrace().
			SetBody(strconv.Itoa(servers.server[id].SvrFailOver.ID)).
			SetHeader("Content-Type", "application/json").
			Post("https://" + servers.server[id].URL + ":5001/api/edit/save")
		if FailOverSaveErr != nil {
			log.Println(FailOverSaveErr)
		}

	case 401:
		// session expired
		log.Printf("Login Session Expired code: %d", LoginStatusCode)
	default:
		log.Printf("Unexpected Login Status code: %d", LoginStatusCode)
	}

}

func AnsibleUpdate(srvType string, rootPassword string) {
	db := dbConn()
	var queryStr string = "SELECT id, Name, Location, URL, UserName, PwStateID, lastUpdated FROM servers WHERE name LIKE \"%" + srvType + "%\" AND AnsibleUpdates=1 ORDER BY Name ASC"
	rows, err := db.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		// Scan one Server record
		var r server
		if err := rows.Scan(&r.id, &r.Name, &r.Location, &r.URL, &r.UserName, &r.PwStateID, &r.lastUpdated); err != nil {
			// handle error
		}
		servers.server = append(servers.server, r)
	}

	for _, r := range servers.server {
		out, err := exec.Command("/usr/bin/bash", "/opt/3CX-Reporting/Scripts/AnsibleUpdate.sh", r.URL, rootPassword).Output()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("The date is %s\n", string(out))
	}

	defer db.Close()

}

func WordingAutomaticUpdatesSingleHost(rootPassword string, hostname string) {

	out, err := exec.Command("/usr/bin/bash", "/opt/3CX-Reporting/Scripts/AnsibleWordingAutomaticUpdates.sh", hostname, rootPassword).Output()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))

}

func WordingAutomaticUpdates(rootPassword string) {
	db := dbConn()
	var queryStr string = "SELECT id, Name, Location, URL, UserName, PwStateID, lastUpdated FROM servers WHERE AnsibleUpdates=1 ORDER BY Name ASC"
	rows, err := db.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		// Scan one Server record
		var r server
		if err := rows.Scan(&r.id, &r.Name, &r.Location, &r.URL, &r.UserName, &r.PwStateID, &r.lastUpdated); err != nil {
			// handle error
		}
		servers.server = append(servers.server, r)
	}

	for _, r := range servers.server {
		out, err := exec.Command("/usr/bin/bash", "/opt/3CX-Reporting/Scripts/AnsibleWordingAutomaticUpdates.sh", r.URL, rootPassword).Output()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(out))
	}

	defer db.Close()

}

// GetFailoverTimeInterval Sets Failover Time Interval
func SetFailoverTimeInterval() {
	db := dbConn()

	var queryStr string = "SELECT id, Name, Location, URL, UserName, PwStateID, lastUpdated FROM servers WHERE name LIKE \"%-3cx-2%\" AND AnsibleUpdates=1 ORDER BY Name ASC"
	rows, err := db.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		// Scan one Server record
		var r server
		if err := rows.Scan(&r.id, &r.Name, &r.Location, &r.URL, &r.UserName, &r.PwStateID, &r.lastUpdated); err != nil {
			// handle error
		}
		servers.server = append(servers.server, r)
	}

	for i, r := range servers.server {
		log.Println("Setting FailOver Time -> " + r.Name)
		SetFailoverTimeIntervalData(i)
	}

	defer db.Close()

}

// GetFailoverTimeInterval gets Failover Time Interval
func GetFailoverTimeInterval() {
	db := dbConn()

	var queryStr string = "SELECT id, Name, Location, URL, UserName, PwStateID, lastUpdated FROM servers WHERE name LIKE \"%-3cx-2%\" AND AnsibleUpdates=1 ORDER BY Name ASC"
	rows, err := db.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		// Scan one Server record
		var r server
		if err := rows.Scan(&r.id, &r.Name, &r.Location, &r.URL, &r.UserName, &r.PwStateID, &r.lastUpdated); err != nil {
			// handle error
		}
		servers.server = append(servers.server, r)
	}

	for i, r := range servers.server {
		log.Println("Getting data from -> " + r.Name)
		GetFailoverTimeIntervalData(i)
	}

	for i, _ := range servers.server {
		s := servers.server[i]
		if s.SvrFailOver.ActiveObject.Interval.Value == 0 {
			//Log
			log.Println("FAILED Getting data from -> " + s.Name)
			continue
		}

		if s.SvrFailOver.ActiveObject.Enabled.Value == true {
			log.Println(s.Name + " -> FailOver Interval " + strconv.Itoa(s.SvrFailOver.ActiveObject.Interval.Value))
		} else {
			log.Println(s.Name + " -> FailOver Disabled")
		}
	}

	defer db.Close()

}

// opsgenieCreateAlarm gets Failover Time Interval
func opsgenieCreateAlarm(message string, boxURL string) {
	alertClient, err := alert.NewClient(&client.Config{
		ApiKey:         C.OpsGenie.APIKey,
		OpsGenieAPIURL: "api.eu.opsgenie.com",
	})
	if err != nil {
		//log.Println(err)
	}
	_, err = alertClient.Create(nil, &alert.CreateAlertRequest{
		Message:     message,
		Alias:       boxURL,
		Description: message,
		Source:      "3CX-Reporting",
		Priority:    alert.P3,
	})
	if err != nil {
		//log.Println(err)
	}
	//log.Println(createResult)
}

// Get3CXwithAutoUpdate gets Failover Time Interval
func Get3CXwithAutoUpdate() {
	db := dbConn()

	var queryStr string = "SELECT id, Name, AutoUpdate, URL FROM servers WHERE lastUpdated IS NOT NULL ORDER BY Name ASC"
	rows, err := db.Query(queryStr)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	for rows.Next() {
		// Scan one Server record
		var r server
		if err := rows.Scan(&r.id, &r.Name, &r.SvrUpdates.ActiveObject.TcxPbxUpdates.Value, &r.URL); err != nil {
			log.Println(err)
		}
		servers.server = append(servers.server, r)
	}

	for i, _ := range servers.server {
		s := servers.server[i]

		if s.SvrUpdates.ActiveObject.TcxPbxUpdates.Value == true {
			opsgenieCreateAlarm(s.Name+" Has Auto Update Enabled", s.URL)
			//log.Println(s.Name + "Auto Update Enabled")
		} else {
			//log.Println(s.Name + " -> Auto Update Disabled")
		}
	}

	defer db.Close()

}

func main() {
	log.Printf("       3CX Tools by Ricardo Ferreira")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/opt/3CX-Reporting")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Unable to read config file")
		viper.GetViper().ConfigFileUsed()

		os.Exit(1)
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Unable to Decode into Struct, %v", err)
	}

	// Check if there is at least 1 arg
	if len(os.Args) <= 1 {
		log.Println("Error: No Arg given!")
		os.Exit(1)
	} else {
		arg := os.Args[1]
		switch arg {
		case "--Get3CXwithAutoUpdate":
			Get3CXwithAutoUpdate()
		case "--GetFailoverTimeInterval":
			GetFailoverTimeInterval()
		case "--SetFailoverTimeInterval":
			arg2 := os.Args[2]
			if _, err := strconv.Atoi(arg2); err == nil {
				log.Println("Set Failover Time Interval")
				targetFailoverTimeInterval, _ = strconv.Atoi(arg2)
				log.Println("Setting Set Failover Time Interval to " + arg2)
				SetFailoverTimeInterval()
			} else {
				log.Println("Set Failover Time Interval NOT VALID")
			}
		case "--AnsibleUpdate3CX-1":
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Root Password: ")
			text, _ := reader.ReadString('\n')
			AnsibleUpdate("3cx-1", text)
		case "--AnsibleUpdate3CX-2":
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Root Password: ")
			text, _ := reader.ReadString('\n')
			AnsibleUpdate("3cx-2", text)
		case "--WordingAutomaticUpdates":
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Root Password: ")
			text, _ := reader.ReadString('\n')
			WordingAutomaticUpdates(text)
		case "--WordingAutomaticUpdatesSingleHost":
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Root Password: ")
			password, _ := reader.ReadString('\n')
			fmt.Print("Hostname: ")
			hostname, _ := reader.ReadString('\n')
			WordingAutomaticUpdatesSingleHost(password, hostname)
		case "--help":
			fmt.Println("DESCRIPTION")
			fmt.Println("       3CXTools is set of tools used to mass manage 3CX Systems.")
			fmt.Println("")
			fmt.Println("COMMAND-LINE OPTIONS")
			fmt.Println("       --GetFailoverTimeInterval					Displays the Failover Time Interval of all 3CX-2")
			fmt.Println("       --SetFailoverTimeInterval XX					Sets the Failover Time Interval of all 3CX-2 to XX secons")
			fmt.Println("       --AnsibleUpdate3CX-1						Updates All 3CX-1 Servers")
			fmt.Println("       --AnsibleUpdate3CX-2						Updates All 3CX-2 Servers")
			fmt.Println("       --WordingAutomaticUpdates					Updates All l10n Files With Warning")
			fmt.Println("       --WordingAutomaticUpdatesSingleHost				Updates All l10n Files With Warning On Single Host")
			fmt.Println("       --help								Display a help message and exit")

			fmt.Println("")
			fmt.Println("AUTHOR")
			fmt.Println("       Written by Ricardo Ferreira ( https://github.com/rmf1995 )")
		default:
			log.Println("Arg Not Valid!")
		}

	}

}
