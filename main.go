package main

import (
	routesAnalaytics "AuthenticationService/routes/Analaytics"
	routesAppointment "AuthenticationService/routes/Appointment"
	routes "AuthenticationService/routes/Authentication"
	routesProfile "AuthenticationService/routes/ProfileService"
	routesUser "AuthenticationService/routes/UserService"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	// Load the DotENV
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌Error loading .env file")
	}

	// ⚠️ Trust only localhost proxy (or none if you want)
	r.SetTrustedProxies(nil)

	// ✅ CORS configuration to allow only one origin
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"}, // Change to your allowed origin
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // allow all origins dynamically
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	//API calls 🚀

	fmt.Println()
	fmt.Println("*****************Authentication*****************")
	fmt.Println()

	//Authentication
	fmt.Println("=================Login=================")
	fmt.Println()
	routes.InitLoginRoutes(r)

	fmt.Println()
	fmt.Println("=================Forget Password=================")
	fmt.Println()
	routes.InitForgetPasswordRoutes(r)

	fmt.Println()
	fmt.Println("*****************UserService*****************")
	fmt.Println()

	//UserService
	fmt.Println("=================File Upload=================")
	fmt.Println()
	routesUser.InitFilesRoutes(r)

	fmt.Println()
	fmt.Println("=================Image Upload=================")
	fmt.Println()
	routesUser.InitImageRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Scan Center=================")
	fmt.Println()
	routesUser.InitScanCenterRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Technician=================")
	fmt.Println()
	routesUser.InitTechnicianRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Receptionist=================")
	fmt.Println()
	routesUser.InitReceptionistRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Patient=================")
	fmt.Println()
	routesUser.InitPatientRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Radiologist=================")
	fmt.Println()
	routesUser.InitRadiologistRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Doctor=================")
	fmt.Println()
	routesUser.InitDoctorRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Scribe=================")
	fmt.Println()
	routesUser.InitScribeRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Manager=================")
	fmt.Println()
	routesUser.InitManagerRoutes(r)

	fmt.Println()
	fmt.Println("=================Create Co-Doctor=================")
	fmt.Println()
	routesUser.InitCoDoctorRoutes(r)

	fmt.Println()
	fmt.Println("=================Create WellthgreenPerformingProvider=================")
	fmt.Println()
	routesUser.InitWellthgreenPerformingProviderRoutes(r)

	fmt.Println()
	fmt.Println("*****************ProfileService*****************")
	fmt.Println()

	//profileService
	fmt.Println("=================View Radiologist=================")
	fmt.Println()
	routesProfile.InitRadiologistRoutes(r)

	fmt.Println()
	fmt.Println("=================View Doctor=================")
	fmt.Println()
	routesProfile.InitDoctorRoutes(r)

	fmt.Println()
	fmt.Println("=================View Co-Doctor=================")
	fmt.Println()
	routesProfile.InitCoDoctorRoutes(r)

	fmt.Println()
	fmt.Println("=================View Receptionist=================")
	fmt.Println()
	routesProfile.InitReceptionistRoutes(r)

	fmt.Println()
	fmt.Println("=================Create WellthgreenPerformingProvider=================")
	fmt.Println()
	routesProfile.InitWellthgreenPerformingProviderRoutes(r)

	fmt.Println()
	fmt.Println("=================View Scribe=================")
	fmt.Println()
	routesProfile.InitScribeRoutes(r)

	fmt.Println()
	fmt.Println("=================View Manager=================")
	fmt.Println()
	routesProfile.InitManagerRoutes(r)

	fmt.Println()
	fmt.Println("=================View Technician=================")
	fmt.Println()
	routesProfile.InitTechnicianRoutes(r)

	fmt.Println()
	fmt.Println("=================View Scan Center=================")
	fmt.Println()
	routesProfile.InitScanCenterRoutes(r)

	fmt.Println()
	fmt.Println("=================View User=================")
	fmt.Println()
	routesProfile.InitUserRoutes(r)

	fmt.Println()
	fmt.Println("*****************AppointmentService*****************")
	fmt.Println()

	//profileService
	fmt.Println("=================Management Appointment=================")
	fmt.Println()
	routesAppointment.InitManageAppointmentRoutes(r)

	fmt.Println()
	fmt.Println("=================Intake Form=================")
	fmt.Println()

	routesAppointment.InitIntakeFormRoutes(r)

	fmt.Println()
	fmt.Println("=================Technician Intake Form=================")
	fmt.Println()

	routesAppointment.InitTechnicianIntakeFormRoutes(r)

	fmt.Println()
	fmt.Println("=================OverRide Form=================")
	fmt.Println()

	routesAppointment.InitOverrideRoutes(r)

	fmt.Println()
	fmt.Println("=================Report Intake Form=================")
	fmt.Println()

	routesAppointment.InitReportIntakeFormRoutes(r)

	fmt.Println()
	fmt.Println("=================Analaytics=================")
	fmt.Println()

	routesAnalaytics.InitAnalayticsRoutes(r)

	fmt.Println()
	fmt.Println("=================Training Material=================")
	fmt.Println()

	routesAnalaytics.InitTrainingMaterialRoutes(r)

	fmt.Println()
	fmt.Println("=================Invoice=================")
	fmt.Println()

	routesAnalaytics.InitInvoiceRoutes(r)

	fmt.Println()
	fmt.Println()

	//Ping 🎯API
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong from Authentication Service",
		})
	})

	//Run the Server and Log Message
	fmt.Println("✅Server is Running at Port:" + os.Getenv("PORT"))
	r.Run("0.0.0.0:" + os.Getenv("PORT"))
}
