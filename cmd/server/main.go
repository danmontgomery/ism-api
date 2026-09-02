package main

import (
	"log"
	"os"

	"dmontgomery/ism-api/api"
	"dmontgomery/ism-api/internal/guidance"
	"dmontgomery/ism-api/internal/guidance/resolvers"
	"dmontgomery/ism-api/internal/handler"
	"dmontgomery/ism-api/internal/refdata"
	"dmontgomery/ism-api/internal/validation"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	reg := refdata.NewRegistry()
	validator := validation.NewEngine(reg)
	guider := guidance.NewEngine(reg,
		&resolvers.ClassificationResolver{},
		&resolvers.CUIResolver{},
		&resolvers.DisseminationResolver{},
		&resolvers.DistributionResolver{},
		&resolvers.AuthorityResolver{},
		&resolvers.DeclassResolver{},
		&resolvers.SCIResolver{},
	)

	r := gin.New()
	h := handler.New(reg, validator, guider, api.Content)
	h.Register(r)

	log.Printf("Starting ISM API server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
