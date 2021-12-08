package discord

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cjtim/cjtim-backend-go/configs"
)

var (
	ErrorsChannel = configs.Config.DISCORD_ERROR_CHANNEL
	ServerID      = configs.Config.DISCORD_SERVER_ID
	token         = "Bot " + configs.Config.DISCORD_TOKEN
)

type DiscordInterface interface {
	SendMsg(channel, title, msg string) error
	// JoinVoice(channel DiscordChannel) (*discordgo.VoiceConnection, error)
	Disconnect() error
}

type DiscordImpl struct {
	client *discordgo.Session
}

func newClient() (DiscordInterface, error) {
	if configs.Config.DISCORD_TOKEN == "" {
		return nil, errors.New("NO DISCORD TOKEN")
	}
	discord, err := discordgo.New(token)
	return &DiscordImpl{client: discord}, err
}

func (s *DiscordImpl) Disconnect() error {
	return s.client.Close()
}

func (s *DiscordImpl) SendMsg(channel, title, msg string) error {
	_, err := s.client.ChannelMessageSendEmbed(string(channel), &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("```%s```", msg),
	})
	return err
}

// func (s *DiscordImpl) JoinVoice(channel DiscordChannel) (*discordgo.VoiceConnection, error) {
// 	err := s.client.Open()
// 	if err != nil {
// 		return nil, err
// 	}
// 	voice, err := s.client.ChannelVoiceJoin(ServerID, string(channel), true, false)
// 	if err != nil {
// 		return voice, err
// 	}
// 	// defer voice.Close()
// 	return voice, nil
// }
