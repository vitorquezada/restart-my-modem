package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

var login string
var senha string

const urlLogin string = "http://192.168.1.1/login_inter.asp"
const urlPaginaReboot string = "http://192.168.1.1/management/reboot.asp"

func main() {

	err := preencherCredenciais()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	defer cancel()

	const selectorCampoLogin string = "#User"
	const selectorCampoSenha string = "#Passwd"
	const botaoEnviar string = `input[type="submit"][value="Login"]`
	const abaManagment string = "li#Management"
	const botaoReboot string = "input#reboot_apply"

	var ch chan bool = make(chan bool)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if _, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			go func() {
				err := chromedp.Run(ctx,
					page.HandleJavaScriptDialog(true),
				)
				if err != nil {
					ch <- false
				}

				ch <- true
			}()
		}
	})

	err = chromedp.Run(ctx,
		chromedp.Navigate(urlLogin),
		chromedp.WaitVisible(selectorCampoLogin),
		chromedp.SetValue(selectorCampoLogin, login),
		chromedp.SetValue(selectorCampoSenha, senha),
		chromedp.Click(botaoEnviar),
		chromedp.WaitVisible(abaManagment),
		chromedp.Navigate(urlPaginaReboot),
		chromedp.Click(botaoReboot),
		chromedp.ActionFunc(func(contexto context.Context) error {
			if <-ch {
				return nil
			}
			return errors.New("Não retornou")
		}),
		chromedp.Sleep(time.Second*5),
		chromedp.Stop(),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func preencherCredenciais() error {
	dados, err := godotenv.Read(".env")
	if err != nil {
		log.Fatalln(err)
	}

	login = dados["LOGIN"]
	senha = dados["SENHA"]

	if login == "" || senha == "" {
		return errors.New("Login ou senha não foram informados no arquivo .env")
	}

	return nil
}
