# Diceroller Bot 
A Discord Bot that rolls dice for you

If you want a new feature or want to give some advice, please open an new issue. It will be very much appreciated.

## Commands
`!r 1d20` - Roll a d20

`!r 2d20+2` - Roll two d20 + 2

`!r d20-1` - Roll a d20 - 1

`!r 2d20k1` - Roll two d20 but keep only the highest one (d&d 5e advantage)

`!r 5d20k4` - Roll five d20 but the keep only the highest 4

`!r (4d8+3)*3d2k1` - Roll four d8, add three and multiply the result for three d2 (keeping only one)

`!r (8d12k1+40)*(1-(2+9d4)-3d2)-4` - You got it

## Installation

```
go get github.com/MicheleLambertucci/diceroller-bot
cd $GOPATH/src/github.com/MicheleLambertucci/diceroller-bot
cd cmd/dice-bot && go install && cd -
cd ../dice-cli && go install && cd -
```

Both `dice-bot` and `dice-cli` will be installed in your $GOPATH/bin directory. Add the folder to your PATH if you want to launch the programs from anywhere.

### Bot
To start the bot use
```
dice-bot -t <your-long-bot-token>
```

### Cli
```
dice-cli
```


If you want to try the bot without deploying your own, you can use my instance.

Add it to your server via [this link](https://discordapp.com/api/oauth2/authorize?client_id=573599563051434025&permissions=0&scope=bot)
