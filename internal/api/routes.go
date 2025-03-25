package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rgehrsitz/me/internal/db"
	"github.com/rgehrsitz/me/internal/services"
)

// Server represents the API server
type Server struct {
	db              *db.DB
	dataDir         string
	router          *gin.Engine
	embeddingService *services.EmbeddingService
	searchService    *services.SearchService
	summarizeService *services.SummarizeService
}

// NewServer creates a new API server
func NewServer(db *db.DB, dataDir string) (*Server, error) {
	// Initialize services
	embeddingService, err := services.NewEmbeddingService()
	if err != nil {
		return nil, err
	}

	summarizeService, err := services.NewSummarizeService()
	if err != nil {
		return nil, err
	}

	searchService := services.NewSearchService(db, embeddingService)

	server := &Server{
		db:               db,
		dataDir:          dataDir,
		embeddingService: embeddingService,
		searchService:    searchService,
		summarizeService: summarizeService,
	}

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	}))

	// Serve frontend static files from the dist directory
	// Use a more specific path to avoid conflicts with API routes
	router.Static("/assets", "./web/dist/assets")
	router.StaticFile("/", "./web/dist/index.html")
	router.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	// API routes
	api := router.Group("/api")
	{
		// Content endpoints
		api.POST("/content", server.CreateContent)
		api.GET("/content", server.ListContent)
		api.GET("/content/:id", server.GetContent)
		api.PUT("/content/:id", server.UpdateContent)
		api.DELETE("/content/:id", server.DeleteContent)
		
		// Embedding endpoints
		api.POST("/content/:id/embed", server.GenerateEmbedding)
		
		// Search endpoints
		api.POST("/search", server.Search)
		
		// Summarization endpoints
		api.POST("/content/:id/summarize", server.SummarizeContent)
		
		// Tags endpoints
		api.GET("/tags", server.ListTags)
		api.POST("/tags", server.CreateTag)
	}

	server.router = router
	return server, nil
}

// Run starts the API server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
