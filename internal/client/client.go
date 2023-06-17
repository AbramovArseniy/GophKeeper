package client

import (
	"context"
	"fmt"
	"log"

	"github.com/AbramovArseniy/GophKeeper/internal/client/httpclient"
	clienttypes "github.com/AbramovArseniy/GophKeeper/internal/client/utils/types"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/manifoldco/promptui"
	"google.golang.org/grpc/metadata"
)

type CommandLine struct {
	action *MDAct
}

type MDAct struct {
	act clienttypes.ClientAction
	md  *metadata.MD
}

func NewAction(address string, md *metadata.MD) (*MDAct, error) {
	mda := &MDAct{md: md, act: httpclient.NewHTTPClient(address)}
	log.Println("Connected successfully!")
	return mda, nil
}

func NewCLI(action *MDAct) *CommandLine {
	return &CommandLine{action: action}
}

func (cli *CommandLine) StartCLI(ctx context.Context) (err error) {
	err = cli.Authentication(ctx)
	if err != nil {
		if err == clienttypes.ErrExitCLI {
			return nil
		}
		return fmt.Errorf("Authentication error: %w", err)
	}
	fmt.Println("Authenticated successfully!")

	if err := cli.Action(ctx); err != nil {
		if err == clienttypes.ErrExitCLI {
			return nil
		}
		return fmt.Errorf("Action error: %w", err)
	}

	return nil
}

func (cli *CommandLine) Authentication(ctx context.Context) error {
	prompt := promptui.Select{
		Label: "Welcome to GophKeeper! What would you like to do?",
		Items: []string{"Register", "Log in", "Exit"},
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error choose action prompt failed: %w", err)
	}
	if idx == 0 {
		err := register(ctx, cli)
		if err != nil {
			return fmt.Errorf("error: can't authorize: %w", err)
		}
		return nil
	}
	if idx == 1 {
		err := authorize(ctx, cli)
		if err != nil {
			return fmt.Errorf("error: can't authorize: %w", err)
		}
		return nil
	}
	if idx == 2 {
		return exitCLI(ctx)
	}
	cli.Action(ctx)
	return nil
}

func register(ctx context.Context, cli *CommandLine) error {
	login, err := getLogin()
	if err != nil {
		return fmt.Errorf("error: can't get username: %w", err)
	}
	password, err := getPassword()
	if err != nil {
		return fmt.Errorf("error: can't get password: %w", err)
	}
	authReq := clienttypes.AuthRequest{
		Login:    login,
		Password: password,
	}
	err = cli.action.act.Register(ctx, authReq)
	if err != nil {
		return fmt.Errorf("error: can't register: %w", err)
	}

	return nil
}

func authorize(ctx context.Context, cli *CommandLine) error {
	login, err := getLogin()
	if err != nil {
		return fmt.Errorf("error: can't get username: %w", err)
	}
	password, err := getPassword()
	if err != nil {
		return fmt.Errorf("error: can't get password: %w", err)
	}
	authReq := clienttypes.AuthRequest{
		Login:    login,
		Password: password,
	}
	err = cli.action.act.Login(ctx, authReq)
	if err != nil {
		return fmt.Errorf("error: can't login: %w", err)
	}

	return nil
}

func getLogin() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter your login: ",
	}
	return prompt.Run()
}

func getPassword() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter your password: ",
		Mask:  '*',
	}
	return prompt.Run()
}

func (cli *CommandLine) Action(ctx context.Context) error {
	prompt := promptui.Select{
		Label: "What would you like to do?",
		Items: []string{"Add secret info", "Get secret info", "Exit"},
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error choose action prompt failed: %w", err)
	}
	if idx == 0 {
		addInfo(ctx, cli.action.act)
	}
	if idx == 1 {
		getInfo(ctx, cli.action.act)
	}
	if idx == 2 {
		exitCLI(ctx)
	}
	cli.Action(ctx)
	return nil
}

func addInfo(ctx context.Context, client clienttypes.ClientAction) {
	infoType := getInfoType()
	infoName := getInfoName()

	switch infoType {
	case storage.LoginPassword:
		req := storage.InfoLoginPass{
			Login:    getValueFromUser("Enter login"),
			Password: getValueFromUser("Enter password"),
		}
		err := client.SaveData(ctx, &req, infoType, infoName)
		if err != nil {
			fmt.Println("Cant save your info!")
		}
		fmt.Println("Password Saved!")
		return
	case storage.Card:
		req := storage.InfoCard{
			CardNumber: getValueFromUser("Enter cardnumber"),
			Holder:     getValueFromUser("Enter cardholder name"),
			Date:       getValueFromUser("Enter expiration date"),
			CVCcode:    getValueFromUser("Enter cvc code"),
		}
		err := client.SaveData(ctx, &req, infoType, infoName)
		if err != nil {
			fmt.Println("Cant save your info!")
		}
		fmt.Println("Card Saved!")
		return
	case storage.Text:
		req := storage.InfoText{
			Text: getValueFromUser("Enter text"),
		}
		err := client.SaveData(ctx, &req, infoType, infoName)
		if err != nil {
			fmt.Println("Cant save your info!")
		}
		fmt.Println("Text Saved!")
		return
	}
}

func getInfo(ctx context.Context, client clienttypes.ClientAction) {
	infoType := getInfoType()
	infoName := getInfoName()

	req := clienttypes.GetRequest{Name: infoName, Type: infoType}

	switch infoType {
	case storage.LoginPassword:
		resp, err := client.GetData(ctx, req)

		if err != nil {
			fmt.Println("Cant get your info!")
		}
		info, ok := resp.(*storage.InfoLoginPass)
		if !ok {
			fmt.Println("Cant get your info!")
		}
		fmt.Printf("Login: %s\n", info.Login)
		fmt.Printf("Password: %s\n", info.Password)
		return
	case storage.Card:
		resp, err := client.GetData(ctx, req)

		if err != nil {
			fmt.Println("Cant get your info!")
		}
		info, ok := resp.(*storage.InfoCard)
		if !ok {
			fmt.Println("Cant get your info!")
		}
		fmt.Printf("CardNumber: %s\n", info.CardNumber)
		fmt.Printf("Holder: %s\n", info.Holder)
		fmt.Printf("Date: %s\n", info.Date)
		fmt.Printf("CVCcode: %s\n", info.CVCcode)
	case storage.Text:
		resp, err := client.GetData(ctx, req)

		if err != nil {
			fmt.Println("Cant get your info!")
		}
		info, ok := resp.(*storage.InfoText)
		if !ok {
			fmt.Println("Cant get your info!")
		}
		fmt.Printf("Text: %s\n", info.Text)
	}
	fmt.Println(infoType, infoName)
}

func getValueFromUser(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}
	value, err := prompt.Run()
	if err != nil {
		log.Fatal("Failed getting login for password secret")
	}

	return value
}

func getInfoName() string {
	prompt := promptui.Prompt{
		Label: "Enter secret name",
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatal("Failed choose secret type prompt")
	}

	return result
}

func getInfoType() storage.InfoType {
	infoTypes := []storage.InfoType{
		storage.LoginPassword,
		storage.Card,
		storage.Text,
	}
	prompt := promptui.Select{
		Label: "Select type of info",
		Items: infoTypes,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		log.Fatal("Failed choose secret type prompt")
	}

	return infoTypes[idx]
}

func exitCLI(ctx context.Context) error {
	return clienttypes.ErrExitCLI
}
