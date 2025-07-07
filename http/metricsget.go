package http

import (
	"fmt"
	"strings"

	"github.com/SoundBoardBot/server-counter/tasks"
	"github.com/gin-gonic/gin"
)

func (s *Server) metricsGetHandler(ctx *gin.Context) {
	accept := ctx.GetHeader("Accept")

	switch {
	case strings.Contains(accept, "text/plain"):
		ctx.Header("Content-Type", "text/plain; version=0.0.4")
		ctx.String(200, fmt.Sprintf(
			"soundboard_guilds %d\nsoundboard_members %d\n",
			tasks.Last_guild_count, tasks.Last_member_count,
		))
	default:
		ctx.JSON(200, gin.H{
			"guilds":  tasks.Last_guild_count,
			"members": tasks.Last_member_count,
		})
	}
}
