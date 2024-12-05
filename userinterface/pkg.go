package userinterface

type Command string

const (
	HELP     Command = "help"
	UPDATE   Command = "update"
	EXCHANGE Command = "exchange"
	EXIT     Command = "exit"
)
