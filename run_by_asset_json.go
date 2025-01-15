package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Params struct {
	Asset                 string `json:"asset"`                   // -asset=[asset name of SPP]
	Jump                  bool   `json:"jump"`                    // -jump=[true or false] default: false
	Step                  bool   `json:"step"`                    // -step=[true or false] default: false
	NavigateOnly          bool   `json:"navigate_only"`           // new
	Account               string `json:"account"`                 // -act=[sign-in account]
	Password              string `json:"password"`                // -pwd=[password]
	TargetUrl             string `json:"target_url"`              // -url=[target url] Exp: "https://portal.azure.com/"
	JumpButtonElement     string `json:"jump_button_element"`     // -jmp-f=[name of jump button element]
	AccountInputElement   string `json:"account_input_element"`   // -act-f=[name of account input element]
	PasswordInputElement  string `json:"password_input_element"`  // -pwd-f=[name of password input element]
	SubmitAccountElement  string `json:"submit_account_element"`  // -smt-a=[name of submit button for account element]
	SubmitPasswordElement string `json:"submit_password_element"` // -smt-p=[name of submit button for password element]
	SubmiSigninElement    string `json:"submit_signin_element"`   // -smt-f=[name of submit button for sign-in element]
	WaitingMs             int    `json:"wait_ms"`                 // -wait=[from 100 to 10000] [msec]
	UseEdge               bool   `json:"use_edge"`                // -edge=[true or false] default: false
	Secret                bool   `json:"secret"`                  // -secret=[true or false] default: true
	SkipCertValidation    bool   `json:"skip_cert_validation"`    // -insec=[true or false] default: false
	Debug                 bool   `json:"debug"`                   // -debug=[true or false] default: false
}

var params_from_json Params

func run_by_asset_json(filename string) bool {
	// Open config file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("'" + filename + "' file could not be opened.")
	}
	// Read config file
	defer file.Close()

	decoder := json.NewDecoder(io.Reader(file))

	for {
		err := decoder.Decode(&params_from_json)
		if err != nil {
			log.Fatal(err)
		}
		if params_from_json.Asset == args.asset {
			// Exist asset
			args.jump = params_from_json.Jump
			args.step = params_from_json.Step
			args.navigateOnly = params_from_json.NavigateOnly
			args.targetUrl = params_from_json.TargetUrl
			args.jumpButtonElement = params_from_json.JumpButtonElement
			args.accountInputElement = params_from_json.AccountInputElement
			args.passwordInputElement = params_from_json.PasswordInputElement
			args.submitForAccountElement = params_from_json.SubmitAccountElement
			args.submitForPasswordElement = params_from_json.SubmitPasswordElement
			args.submitForSignInElement = params_from_json.SubmiSigninElement
			args.waitingTime = params_from_json.WaitingMs
			args.useEdge = params_from_json.UseEdge
			args.secret = params_from_json.Secret
			args.skipCertValidation = params_from_json.SkipCertValidation
			args.debug = params_from_json.Debug
			break
		}
	}
	// Validate arguments
	if validate_args() {
		// Run
		return run()
	}
	return false
}
