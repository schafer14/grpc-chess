package client

// Client is a test implementation for a grpc client
type Client interface {
	NewGameRequest()
}

// Option is a object representing the UCI option object
type Option struct {
	// The option has the name id.
	// 	Certain options have a fixed value for , which means that the semantics of this option is fixed.
	// 	Usually those options should not be displayed in the normal engine options window of the GUI but
	// 	get a special treatment. "Pondering" for example should be set automatically when pondering is
	// 	enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
	// 	automatically by the GUI. All those certain options have the prefix "UCI_" except for the
	// 	first 6 options below. If the GUI get an unknown Option with the prefix "UCI_", it should just
	// 	ignore it and not display it in the engine's options dialog.
	Name string
	// The option has type t.
	// There are 5 different types of options the engine can send
	// * check
	// 	a checkbox that can either be true or false
	// * spin
	// 	a spin wheel that can be an integer in a certain range
	// * combo
	// 	a combo box that can have different predefined strings as a value
	// * button
	// 	a button that can be pressed to send a command to the engine
	// * string
	// 	a text field that has a string as a value,
	// 	an empty string has the value ""
	// todo: make this a enum instead of a string
	Type string
	// the default value of this parameter is x
	Default string
	// the minimum value of this parameter is x
	Min int32
	// the maximum value of this parameter is x
	Max int32
	// a predefined value of this parameter is x
	Var []string
}

// EngineIdent is the identity of an engine including the name and author name
type EngineIdent struct {
	// The name of the engine
	Name string
	// The author of the chess engine
	Author string
}

// Engine defines the required specification for interfacing with the UCI over gRPC protocol
type Engine interface {
	// Id returns the engine name and the engine author
	Init() (EngineIdent, []Option, error)
}
