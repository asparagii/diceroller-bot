all: dice-cli dice-bot

dice-cli:
	pushd cmd/dice-cli && go build 

dice-bot:
	pushd cmd/dice-bot && go build
