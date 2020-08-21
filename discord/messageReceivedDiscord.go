package discord

import (
	"github.com/Selsynn/craft-build-explore-protect-backend /business/user"
	"github.com/Selsynn/craft-build-explore-protect-backend /discord/discordreaction"
	"github.com/bwmarrin/discordgo"
)

type TextReceiveDiscord struct {
	Server  ServerDiscord
	Text    string
	Message MessageID
	User    user.ID
}

type ReactionReceiveDiscord struct {
	Server   ServerDiscord
	Reaction discordreaction.ID
	Message  MessageID
	User     user.ID
}

type MessageID string

func (*ReactionReceiveDiscord) ThisStructIsAMessageReceive() {

}

func (*TextReceiveDiscord) ThisStructIsAMessageReceive() {

}

type MessageSentDiscord struct {
	Server      ServerDiscord
	Text        discordgo.MessageEmbed
	ReactionIDs []discordreaction.ID
}

func (*MessageSentDiscord) ThisStructIsAMessageSent() {

}