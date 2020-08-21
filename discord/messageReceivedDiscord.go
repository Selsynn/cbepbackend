package discord

import (
	"github.com/Selsynn/cbepbackend/business/user"
	"github.com/Selsynn/cbepbackend/communication"
	"github.com/Selsynn/cbepbackend/discord/discordreaction"
	"github.com/bwmarrin/discordgo"
)

type TextReceiveDiscord struct {
	Server  ServerDiscord
	Text    string
	Message communication.ActionID
	User    user.ID
}

type ReactionReceiveDiscord struct {
	Server   ServerDiscord
	Reaction discordreaction.ID
	Message  communication.ActionID
	User     user.ID
}

func (*ReactionReceiveDiscord) ThisStructIsAMessageReceive() {

}

func (*TextReceiveDiscord) ThisStructIsAMessageReceive() {

}

type MessageSentDiscord struct {
	Server      ServerDiscord
	Text        discordgo.MessageEmbed
	ReactionIDs []discordreaction.ID
	ParentErase *communication.ActionID
}

func (*MessageSentDiscord) ThisStructIsAMessageSent() {

}
