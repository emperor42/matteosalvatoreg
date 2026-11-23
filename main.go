package main

import (
	"embed"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed index.html img.png
var assets embed.FS

type Experience struct {
	Company    string   `json:"company"`
	Role       string   `json:"role"`
	Location   string   `json:"location"`
	Period     string   `json:"period"`
	Highlights []string `json:"highlights"`
}

type Education struct {
	Degree   string `json:"degree"`
	School   string `json:"school"`
	Location string `json:"location"`
	Period   string `json:"period"`
	Details  string `json:"details"`
}

type CVData struct {
	Name        string              `json:"name"`
	Title       string              `json:"title"`
	Email       string              `json:"email"`
	Phone       string              `json:"phone"`
	Location    string              `json:"location"`
	Summary     string              `json:"summary"`
	Experience  []Experience        `json:"experience"`
	Education   []Education         `json:"education"`
	Skills      map[string][]string `json:"skills"`
	Credentials []string            `json:"credentials"`
	Links       map[string]string   `json:"links"`
}

// currentCV mirrors the content published in index.html. It is the single
// source of truth for both the page and the JSON API.
func currentCV() CVData {
	return CVData{
		Name:     "Matthew Giancola",
		Title:    "Software Development Expert",
		Email:    "matthew.giancola42@gmail.com",
		Phone:    "+1 514 430 9542",
		Location: "Hamilton, Ontario, Canada",
		Summary: "Software Engineering graduate (Concordia University) and EIT with PEO. " +
			"15+ years across Java, Python, C/C++, PHP, and Go, focused on AI, " +
			"cybersecurity, and Salesforce platform development.",
		Experience: []Experience{
			{
				Company: "Synpulse", Role: "Software Engineer", Location: "Toronto, Ontario", Period: "Current",
				Highlights: []string{
					"Financial services consulting",
					"Owned delivery of an internal inventory system",
					"Certified HyperExponential HX Core Model Developer",
				},
			},
			{
				Company: "Azzurro Technology Inc.", Role: "Founder", Location: "Remote", Period: "Current",
				Highlights: []string{
					"Open-source consulting and self-hosted stack integration",
					"Risk analysis and technical due diligence for engagements",
				},
			},
			{
				Company: "WSP Canada Inc", Role: "OT Cybersecurity Specialist", Location: "Thornhill, Ontario", Period: "Aug 2024 – Oct 2025",
				Highlights: []string{
					"Built Python-based SCADA reporting for water and wastewater compliance",
					"OT network cybersecurity analysis and standards development",
				},
			},
			{
				Company: "Deloitte Inc", Role: "Salesforce Developer", Location: "Ottawa, Ontario", Period: "Aug 2022 – Nov 2023",
				Highlights: []string{
					"LWC, Flows, Apex, VisualForce development",
					"Service, Sales, Non-Profit, Experience, OmniStudio clouds",
				},
			},
			{
				Company: "Nubik Inc", Role: "Junior Salesforce Developer", Location: "Montreal, Québec", Period: "Jan 2022 – Aug 2022",
				Highlights: []string{"Salesforce platform development: LWCs, Flows, Apex"},
			},
			{
				Company: "Riverview Cellars Estate Winery", Role: "Technical Consultant", Location: "Niagara-on-the-Lake, Ontario", Period: "Sep 2014 – Aug 2021",
				Highlights: []string{
					"Developed Canadian shipping pricing system (PHP, WordPress, SQL)",
					"Server migration, website administration, and technical support",
				},
			},
		},
		Education: []Education{
			{
				Degree: "B.Eng, Software Engineering", School: "Concordia University", Location: "Montreal, Québec", Period: "2016 – 2022",
				Details: "CEAB-accredited program. Project: Condo Management System with role-based interfaces for residents, owners, and managers.",
			},
		},
		Skills: map[string][]string{
			"languages":  {"Go (2016+)", "Java (2009+)", "Python (2015+)", "C/C++ (2015+)", "PHP (2014+)", "HTML/CSS/JS (2007+)", "Node.js (2018+)", "Apex (2022+)"},
			"salesforce": {"LWC", "Apex", "VisualForce", "Flow", "Service", "Sales", "Non-Profit", "Experience", "OmniStudio"},
			"tools":      {"Git", "Docker", "Jenkins", "MongoDB", "Linux (Arch, Debian)", "Windows", "WordPress", "Burp Suite"},
			"spoken":     {"English: Fluent", "French: B1", "Italian: B1"},
		},
		Credentials: []string{
			"Salesforce Platform Developer I", "Platform App Builder", "Administrator",
			"Sales Cloud Consultant", "AI Associate", "OmniStudio Developer",
			"IBM Cybersecurity Analyst", "Cyber Threat Intelligence",
			"Pen Testing & Incident Response", "Forensics", "HyperExponential HX Core Model Developer",
		},
		Links: map[string]string{
			"site":        "https://matteosalvatore.dev",
			"github":      "https://www.github.com/Emperor42/",
			"linkedin":    "https://www.linkedin.com/in/matthew-giancola-507072174/",
			"hackerrank":  "https://www.hackerrank.com/matthew_giancol1?hr_r=1",
			"trailblazer": "https://www.salesforce.com/trailblazer/mgiancola1/",
			"hackerone":   "https://www.hackerone.com/emperor42?type=user",
		},
	}
}

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "matteosalvatoreg",
		})
	})

	router.GET("/api/cv", func(c *gin.Context) {
		c.JSON(http.StatusOK, currentCV())
	})

	router.GET("/api/about", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"site":        "matteosalvatore.dev",
			"description": "Responsive single-page CV in English, French and Italian with client-side PDF resume generation.",
			"owner":       "Matthew Giancola",
			"languages":   []string{"en", "fr", "it"},
			"static_host": true,
		})
	})

	router.GET("/", func(c *gin.Context) {
		data, err := assets.ReadFile("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "CV page unavailable")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	router.GET("/img.png", func(c *gin.Context) {
		data, err := assets.ReadFile("img.png")
		if err != nil {
			c.String(http.StatusInternalServerError, "image unavailable")
			return
		}
		c.Data(http.StatusOK, "image/png", data)
	})

	return router
}

func envDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	addr := net.JoinHostPort(envDefault("CV_HOST", "127.0.0.1"), envDefault("CV_PORT", "8081"))
	server := &http.Server{
		Addr:              addr,
		Handler:           newRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("CV site and API listening on http://%s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
