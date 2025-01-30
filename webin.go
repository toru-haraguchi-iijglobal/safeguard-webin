// webin
//
// このソフトウェア webin は One Identity Safeguard のユーティリティです。
//
// One Identity Safeguard は RemoteAppLauncher というツールを使用して
// RDS 上の RemoteApp として、この webin を起動します。
//
// そして webin は ブラウザを起動し、URL やログインID、パスワードなどを
// 自動的に入力します。
//
// RemoteAppLauncher は webin に、資産名、アカウント名、パスワードを引数
// として渡します。
// 当初 webin は、URLやボタンの名前など、すべての必要情報をコマンドライン
// 引数として受け取る設計でした。
// しかし、RemoteAppLauncher が webin を起動する際に webin に渡せる引数が
// Windows UI の制約により 255 文字に制限されており、そのせいで接続できない
// ウェブサイトが出現しました。
// そこで、資産名に対応する設定内容をCSV、JSONに記述して、RDS 上にあらかじめ
// 配置しておく設計に変更されました。
//
// This software, webin, is a utility of One Identity Safeguard.

// One Identity Safeguard uses a tool called RemoteAppLauncher to launch
// this webin as a RemoteApp on RDS.

// Then webin launches a browser and automatically enters the URL, login ID,
// password, etc.

// RemoteAppLauncher passes the asset name, account name, and password as
// arguments to webin.
// Initially, webin was designed to receive all necessary information, such
// as the URL and button name, as command line arguments.
// However, when RemoteAppLauncher launches webin, the arguments that can be
// passed to webin are limited to 255 characters due to Windows UI restrictions,
// and as a result, some websites that could not be connected to have appeared.
// Therefore, the design was changed to write the settings corresponding to
// the asset name in CSV or JSON and place them on RDS in advance.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	// Driver to talk to Chrome-based browsers leveraging the
	// Chrome DevTools protocol
	"github.com/chromedp/chromedp"
)

// Arguments structure
type ArgsStruct struct {
	config                   string // -conf=[config file path] default: "webin.conf"
	jump                     bool   // -jump=[true or false] default: false
	step                     bool   // -step=[true or false] default: false
	navigateOnly             bool   // -navonly=[true or false] default:false
	account                  string // -act=[sign-in account]
	password                 string // -pwd=[password]
	asset                    string // -asset=[asset name of SPP]
	targetUrl                string // -url=[target url] Exp: "https://portal.azure.com/"
	jumpButtonElement        string // -jmp-f=[name of jump button element]
	accountInputElement      string // -act-f=[name of account input element]
	passwordInputElement     string // -pwd-f=[name of password input element]
	submitForAccountElement  string // -smt-a=[name of submit button for account element]
	submitForPasswordElement string // -smt-p=[name of submit button for password element]
	submitForSignInElement   string // -smt-f=[name of submit button for sign-in element]
	waitingTime              int    // -wait=[from 100 to 10000] [msec]
	useEdge                  bool   // -edge=[true or false] default: false
	secret                   bool   // -secret=[true or false] default: true
	skipCertValidation       bool   // -insec=[true or false] default: false
	debug                    bool   // -debug=[true or false] default: false
}

var args ArgsStruct

// Initialze arguments
func init_args() bool {
	flag.StringVar(&args.config, "conf", "webin.conf", "Specify the configuration file")
	// Actions
	flag.BoolVar(&args.jump, "jump", false, "Enable to jump button when before sign-in")
	flag.BoolVar(&args.step, "step", false, "Enable to step by step submit button when Sign-in (account/password)")
	flag.BoolVar(&args.navigateOnly, "navonly", false, "Disable auto sign-in feature and only navigate to the -url")
	// Account, Password and Asset
	flag.StringVar(&args.account, "act", "", "Name of account")
	flag.StringVar(&args.password, "pwd", "", "Password")
	flag.StringVar(&args.asset, "asset", "", "Name of asset")
	// URL, Input elements and Button elements
	flag.StringVar(&args.targetUrl, "url", "", "URL of target system")
	flag.StringVar(&args.jumpButtonElement, "jmp-f", "", "Click first before sign-in form jump button CSS selector")
	flag.StringVar(&args.accountInputElement, "act-f", "", "Account form field CSS selector")
	flag.StringVar(&args.passwordInputElement, "pwd-f", "", "Password form field CSS selector")
	flag.StringVar(&args.submitForAccountElement, "smt-a", "", "Submit for account form button CSS selector")
	flag.StringVar(&args.submitForPasswordElement, "smt-p", "", "Submit for password form button CSS selector")
	flag.StringVar(&args.submitForSignInElement, "smt-f", "", "Submit for login form button CSS selector")
	// Other options
	flag.IntVar(&args.waitingTime, "wait", 1000, "waiting time before inputs are submitted [msec]")
	flag.BoolVar(&args.useEdge, "edge", false, "Using MS Edge instead of Chrome")
	flag.BoolVar(&args.secret, "secret", true, "Using secret mode")
	flag.BoolVar(&args.skipCertValidation, "insec", false, "Skip certificate validation")
	flag.BoolVar(&args.debug, "debug", false, "Enable to debug mode")
	return true
}

// Run
func run() bool {
	// Setting up browser options
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("window-size", "1280,800"),
	)
	if args.useEdge {
		opts = append(opts,
			chromedp.ExecPath("C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"),
		)
	}
	if args.secret {
		opts = append(opts,
			chromedp.Flag("incognito", true),
		)
	}
	if args.skipCertValidation {
		opts = append(opts,
			chromedp.Flag("ignore-certificate-errors", true),
		)
	}
	// for Debug
	// browserOpts := make([]chromedp.ContextOption, 0)
	// if args.debug {
	// 	browserOpts = append(browserOpts,
	// 		chromedp.WithDebugf(log.Printf),
	// 	)
	// }
	// Initializing contexts
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	runCtx, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	// prepare actions
	actions := []chromedp.Action{
		chromedp.Navigate(args.targetUrl),
		chromedp.Sleep(time.Millisecond * time.Duration(args.waitingTime)),
	}

	if !args.navigateOnly {
		if args.jump {
			actions = append(actions,
				chromedp.Click(args.jumpButtonElement, chromedp.ByQuery, chromedp.NodeVisible),
				chromedp.Sleep(time.Millisecond*time.Duration(args.waitingTime)),
			)
		}

		if args.step {
			actions = append(actions,
				// Type in the account and wait
				chromedp.SendKeys(args.accountInputElement, args.account, chromedp.ByQuery, chromedp.NodeVisible),
				chromedp.Sleep(time.Millisecond*time.Duration(args.waitingTime)),
				// Click on submit for account button and wait
				chromedp.Click(args.submitForAccountElement, chromedp.ByQuery, chromedp.NodeVisible),
				chromedp.Sleep(time.Millisecond*time.Duration(args.waitingTime)),
				// Type in Password and wait
				chromedp.SendKeys(args.passwordInputElement, args.password, chromedp.ByQuery, chromedp.NodeVisible),
				chromedp.Sleep(time.Millisecond*time.Duration(args.waitingTime)),
				// Click on submit for password button
				chromedp.Click(args.submitForPasswordElement, chromedp.ByQuery, chromedp.NodeVisible),
			)
		} else {
			actions = append(actions,
				chromedp.SendKeys(args.accountInputElement, args.account, chromedp.ByQuery, chromedp.NodeVisible),
				chromedp.Sleep(time.Millisecond*time.Duration(args.waitingTime)),
				// Type in Password and wait
				chromedp.SendKeys(args.passwordInputElement, args.password, chromedp.ByQuery, chromedp.NodeVisible),
				chromedp.Sleep(time.Millisecond*time.Duration(args.waitingTime)),
				// Click on submit for login button
				chromedp.Click(args.submitForSignInElement, chromedp.ByQuery, chromedp.NodeVisible),
			)
		}

	}
	err := chromedp.Run(runCtx, actions...)

	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Done")
	}

	return true
}

// Validate arguments
func validate_args() bool {
	if args.targetUrl == "" {
		log.Fatalln("Either argument '-url=[target URL]' or JSON entry 'target_url' is required")
	}
	if args.jump {
		if args.jumpButtonElement == "" {
			log.Fatalln("Either argument '-jump-f=[button to sign-in page]' or JSON entry 'jump_button_element' is required")
		}
	}
	if !args.navigateOnly {
		if args.account == "" {
			log.Fatalln("Either argument '-act=[account name]' or JSON entry 'account' is required")
		}
		if args.password == "" {
			log.Fatalln("Either argument '-pwd=[password]' or JSON entry 'password' is required")
		}
		if args.accountInputElement == "" {
			log.Fatalln("Either argument '-act-f=[account input field]' or JSON entry 'account_input_element' is required")
		}
		if args.passwordInputElement == "" {
			log.Fatalln("Either argument '-pwd-f=[password input field]' or JSON entry 'password_input_element' is required")
		}
		if args.step {
			if args.submitForAccountElement == "" {
				log.Fatalln("Either argument '-smt-a=[button to submit account]' or JSON entry 'submit_account_element' is required")
			}
			if args.submitForPasswordElement == "" {
				log.Fatalln("Either argument '-smt-p=[button to submit password]' or JSON entry 'submit_password_element' is required")
			}
		} else {
			if args.submitForSignInElement == "" {
				log.Fatalln("Either argument '-smt-f=[button to submit account/password]' or JSON entry 'submit_signin_element' is required")
			}
		}

	}
	return true
}

// Run (standard mode)
func run_standard() bool {
	// Validate arguments
	if validate_args() {
		// Run
		return run()
	}
	return false
}

// Run (by asset)
// func run_by_asset() bool {
// 	// Validate arguments
// 	if args.config == "" {
// 		log.Fatalln("'-conf=' option [config file] is required")
// 	}
// 	// Open config file
// 	filename := args.config
// 	fp, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatalln("'" + filename + "' file could not be opened.")
// 	}
// 	// Read config file
// 	defer fp.Close()
// 	scanner := bufio.NewScanner(fp)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		count := len(line)
// 		if count == 0 {
// 			continue
// 		}
// 		// Skip comment
// 		f := line[0:1]
// 		if f == "#" {
// 			continue
// 		}
// 		// Split colums
// 		colums := strings.Split(line, ",")
// 		// Find asset
// 		if colums[0] != args.asset {
// 			continue
// 		}
// 		// Exist asset
// 		for index, value := range colums {
// 			switch index {
// 			case 1: // 'jump'
// 				if value == "true" {
// 					args.jump = true
// 				} else if value == "false" {
// 					args.jump = false
// 				}
// 			case 2: // 'step'
// 				if value == "true" {
// 					args.step = true
// 				} else if value == "false" {
// 					args.step = false
// 				}
// 			case 3: // 'url'
// 				if value != "" {
// 					args.targetUrl = value
// 				}
// 			case 4: // 'jmp-f'
// 				if value != "" {
// 					args.jumpButtonElement = value
// 				}
// 			case 5: // 'act-f'
// 				if value != "" {
// 					args.accountInputElement = value
// 				}
// 			case 6: // 'pwd-f'
// 				if value != "" {
// 					args.passwordInputElement = value
// 				}
// 			case 7: // 'smt-f'
// 				if value != "" {
// 					args.submitForSignInElement = value
// 				}
// 			case 8: // 'smt-a'
// 				if value != "" {
// 					args.submitForAccountElement = value
// 				}
// 			case 9: // 'smt-p'
// 				if value != "" {
// 					args.submitForPasswordElement = value
// 				}
// 			case 10: // 'wait'
// 				if value != "" {
// 					var num int
// 					num, _ = strconv.Atoi(value)
// 					args.waitingTime = num
// 				}
// 			case 11: // 'edge'
// 				if value == "true" {
// 					args.useEdge = true
// 				} else if value == "false" {
// 					args.useEdge = false
// 				}
// 			case 12: // 'secret'
// 				if value == "true" {
// 					args.secret = true
// 				} else if value == "false" {
// 					args.secret = false
// 				}
// 			case 13: // 'isec'
// 				if value == "true" {
// 					args.skipCertValidation = true
// 				} else if value == "false" {
// 					args.skipCertValidation = false
// 				}
// 			}
// 		}
// 	}
// 	// Is read error
// 	if err = scanner.Err(); err != nil {
// 		log.Fatalln("'" + filename + "' file read failed.")
// 	}
// 	// Validate arguments
// 	if validate_args() {
// 		// Run
// 		return run()
// 	}
// 	return false
// }

// Main routine
func main() {
	// Initialze arguments
	init_args()
	// Analyze  arguments
	flag.Parse()
	// Is set asset
	if args.asset == "" {
		// Run (standard mode)
		run_standard()
	} else {
		// Run (by asset)
		// run_by_asset()
		run_by_asset_json(args.config)
	}
}
