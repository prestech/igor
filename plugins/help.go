package plugins

import (
	"bytes"
	"strings"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// HelpPlugin provides help functions
type HelpPlugin struct {
	name        string
	description string
}

// Help instantiates the HelpPlugin
func Help() IgorPlugin {
	plugin := HelpPlugin{
		name:        "help",
		description: "I provide help with the following commands",
	}
	return plugin
}

// Work parses the request and ensures a request comes through if any triggers
// are matched. Handled triggers:
//
//  * help
//  * introduce yourself
//  * tell me about yourself
func (HelpPlugin) Work(request slack.Request) (slack.Response, error) {
	response := slack.Response{}
	message := strings.ToLower(request.Text)
	switch message {
	case "help":
		tmpresponse, err := handleHelp(response)
		if err != nil {
			return tmpresponse, err
		}
		response = tmpresponse
	case "introduce yourself":
		response = handleIntroduction(response)
	case "tell me about yourself":
		response = handleTellMe(response)
	case "who am i?":
		response = handleWhoAmI(request, response)
	}
	if response.Text == "" {
		return response, CreateNoMatchError("Nothing found")
	}
	return response, nil
}

// Describe provides the triggers HelpPlugin can handle
func (HelpPlugin) Describe() map[string]string {
	descriptions := make(map[string]string)
	descriptions["help"] = "This help message, providing a list of available commands"
	descriptions["introduce yourself"] = "A public introduction of Igor"
	descriptions["tell me about yourself"] = "Information about Igor"
	descriptions["who am I?"] = "Development information about your account"
	return descriptions
}

func handleHelp(response slack.Response) (slack.Response, error) {
	response.Text = "I can see that you're trying to find an Igor, would you like some help with that?\n"
	response.Text += "You can choose from any of the below listed commands.\n"
	response.Text += "Also, if you prefix a command with an exclamation point, like *!help*, the result will be shown publicly."
	config, err := config.GeneralConfig()
	if err != nil {
		return response, err
	}
	allPlugins := GetPlugins(config)
	c := make(chan slack.Attachment)
	for _, plugin := range allPlugins {
		go func(plugin IgorPlugin) {
			var buffer bytes.Buffer
			for command, description := range plugin.Describe() {
				buffer.WriteString("- *" + command + "*: " + description + "\n")
			}
			attach := slack.Attachment{}
			attach.Title = plugin.Description()
			attach.Text = buffer.String()
			attach.EnableMarkdownFor("text")
			c <- attach
		}(plugin)
	}
	for i := 0; i < len(allPlugins); i++ {
		response.AddAttachment(<-c)
	}
	return response, nil
}

func handleIntroduction(response slack.Response) slack.Response {
	response.Text = "I am Igor, reprethenting We-R-Igors."
	response.SetPublic()
	attach := slack.Attachment{}
	attach.Title = "A Spare Hand When Needed"
	attach.Text = "We come from Überwald, but are alwayth where we are needed motht.\n"
	attach.Text += "Run */igor help* to see which Igors are currently available."
	attach.EnableMarkdownFor("text")
	response.AddAttachment(attach)
	return response
}

func handleTellMe(response slack.Response) slack.Response {
	response.Text = "Originally Igors come from Überwald, but in this world our home is on GitHub."
	attach := slack.Attachment{}
	attach.Title = "GitHub"
	attach.Text = "The main repository is on https://github.com/ArjenSchwarz/igor. Feel free to contribute"
	response.AddAttachment(attach)
	attach = slack.Attachment{}
	attach.Title = "ig.nore.me"
	attach.Text = "Articles written about Igor and its creation can be found on https://ig.nore.me/projects/igor"
	response.AddAttachment(attach)
	return response
}

func handleWhoAmI(request slack.Request, response slack.Response) slack.Response {
	response.Text = "You are not an Igor, and you are the one who commands me. Other than that I don't know much about you. Maybe I'll recognize you if you do some evil laughing?"
	attach := slack.Attachment{}
	attach.Title = "Account details"
	attach.AddField(slack.Field{Title: "Name", Value: request.UserName, Short: true})
	attach.AddField(slack.Field{Title: "UserID", Value: request.UserID, Short: true})
	attach.AddField(slack.Field{Title: "Channel", Value: request.ChannelName, Short: true})
	attach.AddField(slack.Field{Title: "ChannelID", Value: request.ChannelID, Short: true})
	attach.AddField(slack.Field{Title: "Team", Value: request.TeamDomain, Short: true})
	attach.AddField(slack.Field{Title: "TeamID", Value: request.TeamID, Short: true})
	response.AddAttachment(attach)
	return response
}

// Description returns a global description of the plugin
func (p HelpPlugin) Description() string {
	return p.description
}

// Name returns the name of the plugin
func (p HelpPlugin) Name() string {
	return p.name
}
