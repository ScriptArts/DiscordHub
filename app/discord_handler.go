package app

import (
	"encoding/json"
	"errors"
	"github.com/ScriptArts/DiscordHub/models"
	"github.com/ScriptArts/DiscordHub/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/txgruppi/parseargs-go"
	"log"
	"os"
	"strings"
	"time"
)

type DiscordHubCommand struct {
	Command   string   `json:"command"`
	Args      []string `json:"args"`
	GuildID   string   `json:"guild_id"`
	ChannelID string   `json:"channel_id"`
	AuthorID  string   `json:"author_id"`
}

type DiscordHubType int

const (
	TypeSimpleMessage = iota
)

type DiscordHubResponse struct {
	Type    DiscordHubType `json:"type"`
	Message string         `json:"message"`
}

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := os.Getenv("DISCORD_HUB_PREFIX")
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	err := handlerValidate(m.Content)
	if err != nil {
		log.Println(err.Error())
		return
	}

	args, _ := parseargs.Parse(m.Content)
	cmd := args[0]

	c, err := getCommandData(cmd)
	if err != nil {
		log.Println(err.Error())
		return
	}

	dc := getDiscordCommandData(cmd, args, m)

	v, err := json.Marshal(&dc)
	if err != nil {
		log.Println(err.Error())
		return
	}

	nc := utils.GetNatsConnection()
	msg, err := nc.Request(c.Subject, v, 10*time.Second)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var response DiscordHubResponse
	err = json.Unmarshal(msg.Data, &response)
	if err != nil {
		log.Println(err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, response.Message)
}

func handlerValidate(c string) error {
	parsed, err := parseargs.Parse(c)
	if err != nil {
		return err
	}

	v := parsed[0]
	if len(v) < 4 {
		return errors.New("DiscordHubのコマンドではありません")
	}

	return nil
}

func getCommandData(cmd string) (*models.CommandData, error) {
	repo := new(models.CommandRepository)
	return repo.GetCommandData(cmd)
}

func getDiscordCommandData(cmd string, args []string, m *discordgo.MessageCreate) DiscordHubCommand {
	return DiscordHubCommand{
		Command:   cmd,
		Args:      args,
		GuildID:   m.GuildID,
		ChannelID: m.ChannelID,
		AuthorID:  m.Author.ID,
	}
}
