package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper" // https://github.com/spf13/viper#
	"github.com/zpatrick/go-bytesize"
)

// C config var from strut
var C config

type config struct {
	Database database
}

type database struct {
	User string `mapstructure:"User"`
	Pass string `mapstructure:"Pass"`
	Name string `mapstructure:"Name"`
	Host string `mapstructure:"Host"`
	Port string `mapstructure:"Port"`
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

// ServerRecordingUsage receives Server ID and Returns String (1GB/5GB)
func ServerRecordingUsage(id int) string {
	RecordingQuotaBytes := bytesize.Bytesize(servers.server[id].SvrAdditionalStatus.RecordingQuota)
	RecordingQuotaGibibytes := RecordingQuotaBytes.Gibibytes()

	RecordingUsedSpaceBytes := bytesize.Bytesize(servers.server[id].SvrAdditionalStatus.RecordingUsedSpace)
	RecordingUsedSpaceGibibytes := RecordingUsedSpaceBytes.Gibibytes()
	RecordingUsedSpaceMebibytes := RecordingUsedSpaceBytes.Mebibytes()
	RecordingUsedSpaceKibibytes := RecordingUsedSpaceBytes.Kibibytes()
	FristPart := ""

	if RecordingUsedSpaceKibibytes == 0 {
		FristPart = "-"
	} else if RecordingUsedSpaceGibibytes >= 1 {
		FristPart = fmt.Sprintf("%.2f", RecordingUsedSpaceGibibytes) + "GB"
	} else if RecordingUsedSpaceMebibytes >= 1 {
		FristPart = fmt.Sprintf("%.2f", RecordingUsedSpaceMebibytes) + "MiB"
	} else {
		FristPart = fmt.Sprintf("%.2f", RecordingUsedSpaceKibibytes) + "KiB"
	}

	LastPart := fmt.Sprintf("%.2f", RecordingQuotaGibibytes) + "GB"
	return FristPart + " / " + LastPart
}

// ServerTotalMemory receives Server ID and Returns Total RAM round
func ServerTotalMemory(id int) string {
	TotalVirtualMemoryBytes := bytesize.Bytesize(servers.server[id].SvrSystemStatus.TotalVirtualMemory)
	TotalVirtualMemoryGibibytes := TotalVirtualMemoryBytes.Gibibytes()

	result := fmt.Sprintf("%.0f", (math.Round(TotalVirtualMemoryGibibytes)))
	return result
}

// ServerTotalSwap receives Server ID and Returns Total RAM round
func ServerTotalSwap(id int) string {
	TotalPhysicalMemoryBytes := bytesize.Bytesize(servers.server[id].SvrSystemStatus.TotalPhysicalMemory)
	TotalPhysicalMemoryGibibytes := TotalPhysicalMemoryBytes.Gibibytes()

	result := fmt.Sprintf("%.0f", (math.Round(TotalPhysicalMemoryGibibytes)))
	return result
}

// ServerTotalDiskSpace receives Server ID and Returns Total Disk round
func ServerTotalDiskSpace(id int) string {
	TotalDiskSpaceBytes := bytesize.Bytesize(servers.server[id].SvrSystemStatus.TotalDiskSpace)
	TotalDiskSpaceGibibytes := TotalDiskSpaceBytes.Gibibytes()

	result := fmt.Sprintf("%.0f", (math.Round(TotalDiskSpaceGibibytes/10) * 10))
	return result
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

// Get3CXData creates a user session
func Get3CXData(id int) {

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
		SystemStatusResp, SystemStatusErr := client.R().
			EnableTrace().
			Get("https://" + servers.server[id].URL + ":5001/api/SystemStatus")

		if SystemStatusErr != nil {
			log.Println(SystemStatusErr)
		}
		err := json.Unmarshal(SystemStatusResp.Body(), &servers.server[id].SvrSystemStatus)
		if err != nil {
			log.Println(err)
		}

		AdditionalStatusResp, AdditionalStatusErr := client.R().
			EnableTrace().
			Get("https://" + servers.server[id].URL + ":5001/api/SystemStatus/AdditionalStatus")
		if SystemStatusErr != nil {
			log.Println(AdditionalStatusErr)
		}
		err = json.Unmarshal(AdditionalStatusResp.Body(), &servers.server[id].SvrAdditionalStatus)
		if err != nil {
			log.Println(err)
		}

		UpdatesResp, UpdatesErr := client.R().
			EnableTrace().
			Post("https://" + servers.server[id].URL + ":5001/api/UpdateChecker/set")
		if SystemStatusErr != nil {
			log.Println(UpdatesErr)
		}
		err = json.Unmarshal(UpdatesResp.Body(), &servers.server[id].SvrUpdates)
		if err != nil {
			//log.Println(err)
		}

	case 401:
		// session expired
		log.Printf("Login Session Expired code: %d", LoginStatusCode)
	default:
		log.Printf("Unexpected Login Status code: %d", LoginStatusCode)
	}

}

// GetData Retruns All the 3cx Data
func GetData() {
	db := dbConn()

	var queryStr string = "SELECT id, Name, Location, URL, UserName, PwStateID, lastUpdated from servers order by Name asc"
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
		Get3CXData(i)
	}

	for i, r2 := range servers.server {
		RecordingUsage := ServerRecordingUsage(i)
		//TotalRAM := ServerTotalMemory(i)
		TotalSwap := ServerTotalSwap(i)
		TotalDiskSpace := ServerTotalDiskSpace(i)
		s := servers.server[i]
		License := ProductCodeToLicense(s.SvrSystemStatus.ProductCode)
		LicenseExpiration, e := time.Parse(time.RFC3339, fmt.Sprintf("%v", s.SvrSystemStatus.ExpirationDate))
		if e != nil {
			//log.Println("Can't parse time format")
		}

		LicenseExpirationEpoch := LicenseExpiration.Unix()
		timenow := fmt.Sprintf("%d", time.Now().Unix())
		if s.SvrSystemStatus.FQDN == "" {
			//Log
			log.Println("FAILED Getting data from -> " + s.Name)
			continue
		}
		update3cx, err := db.Prepare("UPDATE servers SET Version=?, FQDN=?, CallRecordingUsage=?, MaxSimCalls=?, ExtTotal=?, OSswap=?, OSDiskSpace=?, AutoUpdate=?, License=?, LicenseExpiration=?, ResellerName=?,lastUpdated=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		update3cx.Exec(s.SvrSystemStatus.Version, s.SvrSystemStatus.FQDN, RecordingUsage, s.SvrSystemStatus.MaxSimCalls, s.SvrSystemStatus.ExtensionsTotal, TotalSwap, TotalDiskSpace, s.SvrUpdates.ActiveObject.TcxPbxUpdates.Value, License, LicenseExpirationEpoch, s.SvrSystemStatus.ResellerName, timenow, r2.id)
	}

	defer db.Close()

}

func main() {
	log.Printf("3CX Reporting Tool (3CXtoDatabase) by rmf1995")
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

	GetData()

}
