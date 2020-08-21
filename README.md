To execute the bot (For discord):\
1- Get your discord key\
2- Build the project\
```go build```\
3- Set the environment variable\
```export DiscordBotToken="[REDACTED]"```\
4- Launch the executable\
```./craft-build-explore-protect-backend  -t=$DiscordBotToken```\

Note: Integration test that need the discord bot key will use the environement key

To be sure you can use the environnment key on the test, on vs code, go to settings, `Go: Test Env File` to copy the path to the secret env file
