package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

type CVData struct {
    Name        string   `json:"name"`
    Title       string   `json:"title"`
    Email       string   `json:"email"`
    Experience  []string `json:"experience"`
    Skills      []string `json:"skills"`
    Education   []string `json:"education"`
}

func main() {
    // Initialize Gin router
    router := gin.Default()
    
    // Security middleware
    router.Use(func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Next()
    })
    
    // Health endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "healthy",
            "service": "matteosalvatoreg",
        })
    })
    
    // CV API endpoint
    router.GET("/api/cv", func(c *gin.Context) {
        cv := CVData{
            Name: "Alex Chen",
            Title: "Senior Software Engineer",
            Email: "alex@example.com",
            Experience: []string{
                "10+ years in distributed systems",
                "Expert in microservices architecture",
                "Full-stack development experience",
            },
            Skills: []string{
                "Go, Python, JavaScript",
                "AWS, Docker, Kubernetes",
                "React, Node.js, GraphQL",
            },
            Education: []string{
                "B.S. Computer Science, MIT",
                "M.S. Software Engineering, Stanford",
            },
        }
        c.JSON(200, cv)
    })
    
    // Additional API endpoints
    router.GET("/api/about", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "company": "Tech Innovations Inc.",
            "founded": 2015,
            "employees": 500,
            "services": []string{"cloud", "ai", "consulting"},
        })
    })
    
    // Run the server
    if err := router.Run(":8081"); err != nil {
        fmt.Printf("Error starting server: %v\n", err)
    }
}