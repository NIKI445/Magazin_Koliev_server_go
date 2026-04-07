package middleware

import (
	"strings"

	"RestApiGo/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	// Получаем origin из конфига
	originSite := cfg.CORSOriginSite
	
	// Разрешенные origins
	var allowOrigins []string
	
	if originSite == "" {
		// Если переменная не задана, разрешаем конкретные origins для безопасности
		allowOrigins = []string{
			"https://magazinkoliev.vercel.app",
			"https://magazinkoliev-6ynpmz6tn-nikitas-projects-026b8288.vercel.app",
			"http://localhost:5173",
			"http://localhost:3000",
		}
	} else {
		// Разбиваем строку по запятой, если там несколько origins
		origins := strings.Split(originSite, ",")
		allowOrigins = make([]string, 0, len(origins))
		
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				allowOrigins = append(allowOrigins, trimmed)
			}
		}
		
		// Добавляем localhost для разработки
		allowOrigins = append(allowOrigins, "http://localhost:5173", "http://localhost:3000")
	}
	
	// Удаляем дубликаты
	uniqueOrigins := make([]string, 0, len(allowOrigins))
	seen := make(map[string]bool)
	for _, origin := range allowOrigins {
		if !seen[origin] {
			seen[origin] = true
			uniqueOrigins = append(uniqueOrigins, origin)
		}
	}
	
	return cors.New(cors.Config{
		AllowOrigins:     uniqueOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           86400, // Кэшировать preflight запросы на 24 часа
	})
}