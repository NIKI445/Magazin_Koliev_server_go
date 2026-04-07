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
		// Если переменная не задана, разрешаем все origins (для разработки)
		allowOrigins = []string{"*"}
	} else {
		// Разбиваем строку по запятой, если там несколько origins
		// Пример: "http://localhost:5173,https://magazinkoliev.vercel.app"
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
	
	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}