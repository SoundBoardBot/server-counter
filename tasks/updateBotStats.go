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

func UpdateBotStats() {
	count, err := db.GetGuildCount(context.Background())
	if err != nil {
		utils.Logger.Sugar().Errorf("An error while fetching server count: %w", err)
		return
	}

	PostTopGG(context.Background(), config.Conf.ClientId, config.Conf.Auth.TopGG, count)
	PostDiscordBotList(context.Background(), config.Conf.ClientId, config.Conf.Auth.DiscordBotList, count)

	utils.Logger.Sugar().Infof("Guild Count Updated to %d", count)
}

func PostTopGG(ctx context.Context, botID string, token string, serverCount int) error {
	endpoint := fmt.Sprintf("https://top.gg/api/bots/%s/stats", botID)

	payload := map[string]int{
		"server_count": serverCount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SoundBoardBot/1.0 (+https://github.com/SoundBoardBot/server-counter)")
	req.Header.Set("Authorization", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status posting to Top.gg: %s", resp.Status)
	}

	return nil
}

func PostDiscordBotList(ctx context.Context, botSlug string, token string, guilds int /*, users int*/) error {
	endpoint := fmt.Sprintf("https://discordbotlist.com/api/v1/bots/%s/stats", botSlug)

	payload := map[string]int{
		"guilds": guilds,
		// "users":  users,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SoundBoardBot/1.0 (+https://github.com/SoundBoardBot/server-counter)")
	req.Header.Set("Authorization", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bad status posting to DBL: %s", resp.Status)
	}

	return nil
}
