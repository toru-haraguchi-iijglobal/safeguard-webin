package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Params struct {
	Asset                 string `json:"asset"`                   // -asset=[asset name of SPP]
	Jump                  bool   `json:"jump"`                    // -jump=[true or false] default: false
	Step                  bool   `json:"step"`                    // -step=[true or false] default: false
	NavigateOnly          bool   `json:"navigate_only"`           // -navonly=[true or false] default:false
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

// var params_from_json Params

func run_by_asset_json(json_filename string) bool {
	// Open error logfile and redirect log
	exe_path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exe_path = filepath.Dir(exe_path)
	log_filename := exe_path + "\\webin_" + strconv.Itoa(os.Getpid()) + ".log"
	logfile, err := os.OpenFile(log_filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	// Open config jsonfile
	jsonfile, err := os.Open(json_filename)
	if err != nil {
		log.Fatalln("'" + json_filename + "' file could not be opened.")
	}
	defer jsonfile.Close()
	decoder := json.NewDecoder(io.Reader(jsonfile))

	// Scan jsonfile for the asset specified
	for {
		var params_from_json Params
		err := decoder.Decode(&params_from_json)
		if err != nil {
			if err == io.EOF {
				log.Println("Asset \"" + args.asset + "\" not found")
				log.Fatal(err)
			} else {
				log.Println("Last correct asset: " + params_from_json.Asset)
				log.Fatal(err)
			}
		}

		if params_from_json.Asset == args.asset {
			// The asset specified exists
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
		// Everything correct including JSON file. No log needed.
		logfile.Close()
		os.Remove(log_filename)
		// Run
		return run()
	}
	return false
}
