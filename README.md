# Diceroller Bot (aka Giorgio)
A Discord Bot that rolls dice for you

#### Disclaimer
I started this project with the main goal of learning some Go.
This is my first time with Go, thus the code is (at this point in time) a big mess.
I look forward to improving the code, especially the parsing part.

If you want a new feature or want to give some advice, file a new issue.

## Commands
`!r 1d20` - Roll a d20

`!r 2d20+2` - Roll two d20 + 2

`!r d20-1` - Roll a d20 - 1

`!r 2d20k1` - Roll two d20 but keep only the highest one (d&d 5e advantage)

`!r 5d20k4` - Roll five d20 but the keep only the highest 4

`!r 99d100k1+40d4-3d2-4` - You got it

## Usage
```
go get github.com/MicheleLambertucci/diceroller-bot
cd $GOPATH/github.com/MicheleLambertucci/diceroller-bot
go build
./diceroller-bot -t <your-long-bot-token>
```

If you want to try the bot without deploying your own, you can use my instance.

Add it to your server via [this link](https://discordapp.com/api/oauth2/authorize?client_id=573599563051434025&permissions=0&scope=bot)
