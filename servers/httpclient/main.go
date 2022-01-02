package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Benyam-S/aseri/client/http/handler"
	"github.com/Benyam-S/aseri/client/http/session"
	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/log"
	"github.com/Benyam-S/aseri/tools"
	"github.com/go-redis/redis"

	cmRepository "github.com/Benyam-S/aseri/common/repository"
	cmService "github.com/Benyam-S/aseri/common/service"

	stRepository "github.com/Benyam-S/aseri/staff/repository"
	stService "github.com/Benyam-S/aseri/staff/service"

	pwRepository "github.com/Benyam-S/aseri/password/repository"
	pwService "github.com/Benyam-S/aseri/password/service"

	jaRepository "github.com/Benyam-S/aseri/jobapplication/repository"
	jaService "github.com/Benyam-S/aseri/jobapplication/service"

	jbRepository "github.com/Benyam-S/aseri/job/repository"
	jbService "github.com/Benyam-S/aseri/job/service"

	urRepository "github.com/Benyam-S/aseri/user/repository"
	urService "github.com/Benyam-S/aseri/user/service"

	fdRepository "github.com/Benyam-S/aseri/feedback/repository"
	fdService "github.com/Benyam-S/aseri/feedback/service"

	sbRepository "github.com/Benyam-S/aseri/subscription/repository"
	sbService "github.com/Benyam-S/aseri/subscription/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	configFilesDir string
	redisClient    *redis.Client
	mysqlDB        *gorm.DB
	sysConfig      SystemConfig
	err            error
	staffHandler   *handler.StaffHandler
)

// SystemConfig is a type that defines a server system configuration file
type SystemConfig struct {
	RedisClient          map[string]string `json:"redis_client"`
	MysqlClient          map[string]string `json:"mysql_client"`
	CookieName           string            `json:"cookie_name"`
	SecretKey            string            `json:"secret_key"`
	SuperAdminEmail      string            `json:"super_admin_email"`
	HTTPDomainAddress    string            `json:"http_domain_address"`
	BotDomainAddress     string            `json:"bot_domain_address"`
	HTTPClientServerPort string            `json:"http_client_server_port"`
	BotClientServerPort  string            `json:"bot_client_server_port"`
	LogFile              string            `json:"server_log_file"`
}

// initServer initialize the web server for takeoff
func initServer() {

	// Reading data from config.server.json file and creating the systemconfig  object
	sysConfigDir := filepath.Join(configFilesDir, "/config.server.json")
	sysConfigData, _ := ioutil.ReadFile(sysConfigDir)

	// Reading data from config.aseri.json file
	aseriConfig := make(map[string]interface{})
	aseriConfigDir := filepath.Join(configFilesDir, "/config.aseri.json")
	aseriConfigData, _ := ioutil.ReadFile(aseriConfigDir)

	err = json.Unmarshal(sysConfigData, &sysConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(aseriConfigData, &aseriConfig)
	if err != nil {
		panic(err)
	}

	apiAccessPoint, ok1 := aseriConfig["api_access_point"].(string)
	botAPIToken, ok2 := aseriConfig["bot_api_token"].(string)

	if !ok1 || !ok2 {
		panic(errors.New("unable to parse aseri config data"))
	}

	// Setting environmental variables so they can be used any where on the application
	os.Setenv("config_files_dir", configFilesDir)
	os.Setenv("super_admin_email", sysConfig.SuperAdminEmail)
	os.Setenv(entity.AppSecretKeyName, sysConfig.SecretKey)
	os.Setenv(entity.AppCookieName, sysConfig.CookieName)
	os.Setenv("http_domain_address", sysConfig.HTTPDomainAddress)
	os.Setenv("bot_domain_address", sysConfig.BotDomainAddress)
	os.Setenv("http_client_server_port", sysConfig.HTTPClientServerPort)
	os.Setenv("bot_client_server_port", sysConfig.BotClientServerPort)
	os.Setenv("log_file", sysConfig.LogFile)

	os.Setenv("api_access_point", apiAccessPoint)
	os.Setenv("bot_api_token", botAPIToken)

	// Initializing the database with the needed tables and values
	initDB()

	// file server
	path, _ := os.Getwd()
	templatePath := filepath.Join(path, "../../ui/templates/*.html")
	temp := template.Must(template.ParseGlob(templatePath))

	staffRepo := stRepository.NewStaffRepository(mysqlDB)
	passwordRepo := pwRepository.NewPasswordRepository(mysqlDB)
	jobApplicationRepo := jaRepository.NewJobApplicationRepository(mysqlDB)
	jobRepo := jbRepository.NewJobRepository(mysqlDB)
	userRepo := urRepository.NewUserRepository(mysqlDB)
	subscriptionRepo := sbRepository.NewSubscriptionRepository(mysqlDB)
	feedbackRepo := fdRepository.NewFeedbackRepository(mysqlDB)
	commonRepo := cmRepository.NewCommonRepository(mysqlDB)
	sessionRepo := stRepository.NewSessionRepository(mysqlDB)

	commonService := cmService.NewCommonService(commonRepo)
	staffService := stService.NewStaffService(staffRepo, passwordRepo, sessionRepo, commonRepo)
	passwordService := pwService.NewPasswordService(passwordRepo, commonRepo)
	userService := urService.NewUserService(userRepo, jobRepo, jobApplicationRepo, commonRepo)
	jobService := jbService.NewJobService(jobRepo, userRepo, commonService)
	jobApplicationService := jaService.NewJobApplicationService(jobApplicationRepo, commonRepo)
	subscriptionService := sbService.NewSubscriptionService(subscriptionRepo, commonService)
	feedbackService := fdService.NewFeedbackService(feedbackRepo, userRepo)

	// Creating push channel and queue
	pushChannel := make(chan string, 1000)
	pushQueue := cmService.NewPushQueue()
	logger := &log.Logger{ServerLogFile: filepath.Join(path, "../../log", os.Getenv("log_file"))}
	store := tools.NewRedisStore(redisClient)

	staffHandler = handler.NewStaffHandler(staffService, passwordService, userService,
		jobService, jobApplicationService, subscriptionService, feedbackService,
		commonService, temp, store, pushChannel, pushQueue, logger)
}

// initDB initialize the database for takeoff
func initDB() {

	redisDB, err := strconv.ParseInt(sysConfig.RedisClient["database"], 0, 0)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     sysConfig.RedisClient["address"] + ":" + sysConfig.RedisClient["port"],
		Password: sysConfig.RedisClient["password"], // no password set
		DB:       int(redisDB),                      // use default DB
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	if err != nil {
		panic(err)
	}

	mysqlDB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		sysConfig.MysqlClient["user"], sysConfig.MysqlClient["password"],
		sysConfig.MysqlClient["address"], sysConfig.MysqlClient["port"], sysConfig.MysqlClient["database"]))

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database: mysql @GORM")

	// Creating and Migrating tables from the structures
	mysqlDB.AutoMigrate(entity.Staff{})
	mysqlDB.AutoMigrate(session.ServerSession{})
	mysqlDB.AutoMigrate(entity.Password{})
	mysqlDB.AutoMigrate(&entity.Feedback{})
	mysqlDB.AutoMigrate(&entity.Subscription{})
	mysqlDB.AutoMigrate(&entity.JobApplication{})
	mysqlDB.AutoMigrate(&entity.Job{})
	mysqlDB.AutoMigrate(&entity.User{})
	mysqlDB.AutoMigrate(&entity.Referral{})
	mysqlDB.AutoMigrate(&entity.ReferralVisit{})

	// Setting foreign key constraint
	mysqlDB.Model(&entity.JobApplication{}).AddForeignKey("job_seeker_id", "users(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.JobApplication{}).AddForeignKey("job_id", "jobs(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.Feedback{}).AddForeignKey("user_id", "users(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.Subscription{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.ReferralVisit{}).AddForeignKey("code", "referrals(code)", "CASCADE", "CASCADE")
}

func main() {
	configFilesDir = "/aseri/config"

	// Initializing the server
	initServer()
	defer mysqlDB.Close()

	router := mux.NewRouter()

	staffRoutes(staffHandler, router)
	userRoutes(staffHandler, router)
	jobRoutes(staffHandler, router)
	feedbackRoutes(staffHandler, router)
	sessionRoutes(staffHandler, router)
	dashboardRoutes(staffHandler, router)
	externalRoutes(staffHandler, router)

	// Static file serving
	fileServer1 := http.FileServer(http.Dir("../../ui"))
	router.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", fileServer1))

	fileServer2 := http.FileServer(http.Dir("../../assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer2))

	go func() {
		staffHandler.HandlePushRequest()
	}()

	go func() {
		staffHandler.HandleCloseable()
	}()

	http.ListenAndServeTLS(":"+os.Getenv("http_client_server_port"),
		filepath.Join(configFilesDir, "/server.pem"),
		filepath.Join(configFilesDir, "/server.key"), router)
}

// dashboardRoutes is a function that defines job attribute and dashboard routes
func dashboardRoutes(handler *handler.StaffHandler, router *mux.Router) {

	router.HandleFunc("/staff/dashboard", tools.MiddlewareFactory(handler.HandleDashboard,
		handler.Authorization, handler.SessionAuthentication))

	router.HandleFunc("/staff/dashboard/init/{csrf}", tools.MiddlewareFactory(handler.HandleInitDashboard,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	/* +++++++++++++++++++++++++++ JOB ATTRIBUTE MANAGEMENT +++++++++++++++++++++++++++ */

	router.HandleFunc("/staff/attribute", tools.MiddlewareFactory(handler.HandleJobAttribute,
		handler.Authorization, handler.SessionAuthentication))

	router.HandleFunc("/staff/attributes/{csrf}", tools.MiddlewareFactory(handler.HandleViewJobAttributes,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/attribute/{category}/{csrf}", tools.MiddlewareFactory(handler.HandleAddJobAttribute,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/staff/attribute/{category}/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleDeleteJobAttribute,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("DELETE")

	/* +++++++++++++++++++++++++++ DEFAULT PAGE +++++++++++++++++++++++++++ */

	router.HandleFunc("/", tools.MiddlewareFactory(handler.HandleDashboard,
		handler.Authorization, handler.SessionAuthentication))
}

// sessionRoutes is a function that defines login, logout and dashboard routes
func sessionRoutes(handler *handler.StaffHandler, router *mux.Router) {

	router.HandleFunc("/staff/login", handler.HandleLogin)
	router.HandleFunc("/staff/logout", tools.MiddlewareFactory(handler.HandleLogout,
		handler.Authorization, handler.SessionAuthentication))
}

// staffRoutes is a function that defines all the routes related to staff
func staffRoutes(handler *handler.StaffHandler, router *mux.Router) {

	router.HandleFunc("/admin/staff", tools.MiddlewareFactory(handler.HandleStaffMember,
		handler.Authorization, handler.SessionAuthentication))

	router.HandleFunc("/admin/staff/{csrf}", tools.MiddlewareFactory(handler.HandleAddStaffMember,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/admin/register/init/{csrf}", tools.MiddlewareFactory(handler.HandleInitAddAdmin,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/admin/register/finish/{csrf}", tools.MiddlewareFactory(handler.HandleFinishAddAdmin,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	/* +++++++++++++++++++++++++++ STAFF MEMBER PROFILE MANAGEMENT +++++++++++++++++++++++++++ */

	router.HandleFunc("/admin/staff/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleViewStaffMember,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/admin/staff/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleDeleteStaffMember,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/admin/staff/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleUpdateStaffMember,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("PUT")

	router.HandleFunc("/admin/staff/{id}/profile/password/{csrf}", tools.MiddlewareFactory(handler.HandleChangeStaffMemberPassword,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("PUT")

	router.HandleFunc("/admin/staff/{id}/profile/pic/{csrf}", tools.MiddlewareFactory(handler.HandleUploadStaffMemberPhoto,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/admin/staff/{id}/profile/pic/{csrf}", tools.MiddlewareFactory(handler.HandleRemoveStaffMemberPhoto,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("DELETE")

	/* +++++++++++++++++++++++++++ STAFF MEMBER SEARCHING +++++++++++++++++++++++++++ */

	router.HandleFunc("/admin/staffs/all/{csrf}", tools.MiddlewareFactory(handler.HandleAllStaffMembers,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/admin/staffs/search/{csrf}", tools.MiddlewareFactory(handler.HandleSearchStaffMembers,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	/* +++++++++++++++++++++++++++ CURRENT STAFF MEMBER +++++++++++++++++++++++++++ */

	router.HandleFunc("/staff/profile/{csrf}", tools.MiddlewareFactory(handler.HandleViewProfile,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/profile/{csrf}", tools.MiddlewareFactory(handler.HandleUpdateProfile,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("PUT")

	router.HandleFunc("/staff/profile/{csrf}", tools.MiddlewareFactory(handler.HandleDeleteAccount,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/staff/profile/password/{csrf}", tools.MiddlewareFactory(handler.HandleChangePassword,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("PUT")

	router.HandleFunc("/staff/profile/pic/{csrf}", tools.MiddlewareFactory(handler.HandleUploadPhoto,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/staff/profile/pic/{csrf}", tools.MiddlewareFactory(handler.HandleRemovePhoto,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("DELETE")
}

// userRoutes is a function that defines all the routes related to user
func userRoutes(handler *handler.StaffHandler, router *mux.Router) {

	router.HandleFunc("/staff/user", tools.MiddlewareFactory(handler.HandleUser,
		handler.Authorization, handler.SessionAuthentication))

	router.HandleFunc("/staff/user/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleViewUser,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/user/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleDeleteUser,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	// /* +++++++++++++++++++++++++++ USER SEARCHING +++++++++++++++++++++++++++ */

	router.HandleFunc("/staff/users/all/{csrf}", tools.MiddlewareFactory(handler.HandleAllUsers,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/users/search/{csrf}", tools.MiddlewareFactory(handler.HandleSearchUsers,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")
}

// jobRoutes is a function that defines all the routes related to job
func jobRoutes(handler *handler.StaffHandler, router *mux.Router) {

	router.HandleFunc("/staff/job", tools.MiddlewareFactory(handler.HandleJob,
		handler.Authorization, handler.SessionAuthentication))

	router.HandleFunc("/staff/job/{csrf}", tools.MiddlewareFactory(handler.HandleAddJob,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/staff/jobs/{csrf}", tools.MiddlewareFactory(handler.HandleAddMultipleJobs,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("POST")

	router.HandleFunc("/staff/job/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleViewJob,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/job/initiator/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleViewJobInitiator,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/job/init/add/{csrf}", tools.MiddlewareFactory(handler.HandleInitAddJob,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/job/{id}/status/{csrf}", tools.MiddlewareFactory(handler.HandleChangeJobStatus,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("PUT")

	router.HandleFunc("/staff/job/{id}/push/{csrf}", tools.MiddlewareFactory(handler.HandleRePushJob,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	// /* +++++++++++++++++++++++++++ JOB SEARCHING +++++++++++++++++++++++++++ */

	router.HandleFunc("/staff/jobs/all/{csrf}", tools.MiddlewareFactory(handler.HandleAllJobs,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/jobs/search/{csrf}", tools.MiddlewareFactory(handler.HandleSearchJobs,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")
}

// feedbackRoutes is a function that defines all the routes related to feedback
func feedbackRoutes(handler *handler.StaffHandler, router *mux.Router) {

	router.HandleFunc("/staff/feedback", tools.MiddlewareFactory(handler.HandleFeedback,
		handler.Authorization, handler.SessionAuthentication))

	router.HandleFunc("/staff/feedback/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleViewFeedback,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/feedback/{id}/{csrf}", tools.MiddlewareFactory(handler.HandleMarkFeedbackAsSeen,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("PUT")

	// /* +++++++++++++++++++++++++++ FEEDBACK SEARCHING +++++++++++++++++++++++++++ */

	router.HandleFunc("/staff/feedbacks/all/{csrf}", tools.MiddlewareFactory(handler.HandleAllFeedbacks,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")

	router.HandleFunc("/staff/feedbacks/search/{csrf}", tools.MiddlewareFactory(handler.HandleSearchFeedbacks,
		handler.CSRFAuthentication, handler.Authorization, handler.SessionAuthentication)).Methods("GET")
}

// externalRoutes is a function that defines all the routes related external system
func externalRoutes(handler *handler.StaffHandler, router *mux.Router) {

	router.HandleFunc("/job/{user_id}/{access_token}", handler.HandleInitAddJobExternal).Methods("GET")

	router.HandleFunc("/job/{user_id}/{csrf}", tools.MiddlewareFactory(handler.HandleAddJobExternal,
		handler.CSRFAuthenticationExternal)).Methods("POST")

	router.HandleFunc("/referral/{code}", handler.HandleNewVisitReferral).Methods("POST")
}
