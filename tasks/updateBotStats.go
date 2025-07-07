package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SoundBoardBot/server-counter/config"
	"github.com/SoundBoardBot/server-counter/db"
	"github.com/SoundBoardBot/server-counter/utils"
)

var Last_guild_count = 0
var Last_member_count = 0

func UpdateBotStats() {
	guild_count, err := db.GetGuildCount(context.Background())
	if err != nil {
		utils.Logger.Sugar().Errorf("An error while fetching server count: %w", err)
		return
	}
	member_count, err := db.GetMemberCount(context.Background())
	if err != nil {
		utils.Logger.Sugar().Errorf("An error while fetching member count: %w", err)
		return
	}
	if guild_count == Last_guild_count && member_count == Last_member_count {
		return
	}
	Last_guild_count = guild_count
	Last_member_count = member_count

	ctx := context.Background()
	if config.Conf.Auth.TopGG != "" {
		postStats(ctx, fmt.Sprintf("https://top.gg/api/bots/%s/stats", config.Conf.ClientId), config.Conf.Auth.TopGG, map[string]int{
			"server_count": guild_count,
		})
	}
	if config.Conf.Auth.DiscordBotList != "" {
		postStats(ctx, fmt.Sprintf("https://discordbotlist.com/api/v1/bots/%s/stats", config.Conf.ClientId), config.Conf.Auth.DiscordBotList, map[string]int{
			"guilds": guild_count,
			"users":  member_count,
		})
	}
	if config.Conf.Auth.BotListMe != "" {
		postStats(ctx, fmt.Sprintf("https://api.botlist.me/api/v1/bots/%s/stats", config.Conf.ClientId), config.Conf.Auth.BotListMe, map[string]int{
			"server_count": guild_count,
			// "shard_count": 0
		})
	}
	if config.Conf.Auth.VoidBots != "" {
		postStats(ctx, fmt.Sprintf("https://api.voidbots.net/bot/stats/%s", config.Conf.ClientId), config.Conf.Auth.VoidBots, map[string]int{
			"server_count": guild_count,
			// "shard_count": 0
		})
	}

	utils.Logger.Sugar().Infof("Guild Count Updated to %d, Member Count Updated to %d", guild_count, member_count)
}

func postStats(ctx context.Context, url, token string, payload map[string]int) {
	body, err := json.Marshal(payload)
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to marshal payload: %v", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SoundBoardBot/1.0 (+https://github.com/SoundBoardBot/server-counter)")
	req.Header.Set("Authorization", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		utils.Logger.Sugar().Errorf("Request to %s failed: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		utils.Logger.Sugar().Errorf("Bad response from %s: %s", url, resp.Status)
	}
}
