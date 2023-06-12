package client

import (
	"context"
	"fmt"
	"log"

	"github.com/AbramovArseniy/GophKeeper/internal/server/handlers"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
	"github.com/manifoldco/promptui"
	"google.golang.org/grpc/metadata"
)

func NewAction(address string, md *metadata.MD) (*MDAct, error) {
	mda := &MDAct{md: md}
	if err := mda.Connection(address); err != nil {
		return nil, fmt.Errorf("error: wrong connection: %w", err)
	}

	log.Println("Connected successfully!")
	return mda, nil
}

func NewCLI(action *MDAct) *CommandLine {
	return &CommandLine{action: action}
}

func (cli *CommandLine) StartCLI(ctx context.Context) (err error) {
	// получаем токен в следующей строке, он нужен будет? пока стоит как прочерк
	_, err = cli.Authentication(ctx)
	if err != nil {
		if err == ErrExitCLI {
			return nil
		}
		return fmt.Errorf("Authentication error: %w", err)
	}
	fmt.Println("Authenticated successfully!")

	if err := cli.Action(ctx); err != nil {
		if err == ErrExitCLI {
			return nil
		}
		return fmt.Errorf("Action error: %w", err)
	}

	return nil
}

func (cli *CommandLine) Authentication(ctx context.Context) (string, error) {
	prompt := promptui.Select{
		Label: "Welcome to GophKeeper! What would you like to do?",
		Items: []string{"Register", "Authorize", "Exit"},
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("error choose action prompt failed: %w", err)
	}
	if idx == 0 {
		token, err := register(ctx)
		if err != nil {
			return "", fmt.Errorf("error: can't authorize: %w", err)
		}
		return token, nil
	}
	if idx == 1 {
		token, err := authorize(ctx)
		if err != nil {
			return "", fmt.Errorf("error: can't authorize: %w", err)
		}
		return token, nil
	}
	if idx == 2 {
		return "", exitCLI(ctx)
	}
	cli.Action(ctx)
	return "", nil
}

func register(ctx context.Context) (string, error) {
	login, err := getLogin()
	if err != nil {
		return "", fmt.Errorf("error: can't get username: %w", err)
	}
	password, err := getPassword()
	if err != nil {
		return "", fmt.Errorf("error: can't get password: %w", err)
	}
	request := types.User{
		Login:        login,
		HashPassword: password,
	}
	var authjwt handlers.AuthJWT
	token, err := authjwt.GenerateToken(request)
	// var ustore storage.UserStorage
	// response, err := ustore.RegisterNewUser(request.Login, request.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("error: can't login: %w", err)
	}

	return token, nil
}

func authorize(ctx context.Context) (string, error) {
	login, err := getLogin()
	if err != nil {
		return "", fmt.Errorf("error: can't get username: %w", err)
	}
	password, err := getPassword()
	if err != nil {
		return "", fmt.Errorf("error: can't get password: %w", err)
	}
	request := types.User{
		Login:        login,
		HashPassword: password,
	}
	var authjwt handlers.AuthJWT
	// Тут тоже генерация токена?
	token, err := authjwt.GenerateToken(request)
	// var ustore storage.UserStorage
	// response, err := ustore.FindUser(request.Login)
	if err != nil {
		return "", fmt.Errorf("error: can't login: %w", err)
	}

	return token, nil
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

func addInfo(ctx context.Context, client ClientAction) {
	infoType := getInfoType()
	infoName := getInfoName()

	switch infoType {
	case storage.LoginPassword:
		req := storage.InfoLoginPass{
			Name:     infoName,
			Login:    getValueFromUser("Enter login"),
			Password: getValueFromUser("Enter password"),
		}
		err := client.SavePassword(ctx, req)
		if err != nil {
			fmt.Println("Cant save your info!")
		}
		fmt.Println("Password Saved!")
		return
	case storage.Card:
		req := storage.InfoCard{
			CardName:   infoName,
			CardNumber: getValueFromUser("Enter cardnumber"),
			Holder:     getValueFromUser("Enter cardholder name"),
			Date:       getValueFromUser("Enter expiration date"),
			CVCcode:    getValueFromUser("Enter cvc code"),
		}
		err := client.SaveCard(ctx, req)
		if err != nil {
			fmt.Println("Cant save your info!")
		}
		fmt.Println("Card Saved!")
		return
	case storage.Text:
		req := storage.InfoText{
			Name: infoName,
			Text: getValueFromUser("Enter text"),
		}
		err := client.SaveText(ctx, req)
		if err != nil {
			fmt.Println("Cant save your info!")
		}
		fmt.Println("Text Saved!")
		return
	}
}

func getInfo(ctx context.Context, client ClientAction) {
	infoType := getInfoType()
	infoName := getInfoName()

	req := GetRequest{Name: infoName}

	switch infoType {
	case storage.LoginPassword:
		resp, err := client.GetPassword(ctx, req)
		if err != nil {
			fmt.Println("Cant get your info!")
		}
		fmt.Printf("Login: %s\n", resp.Login)
		fmt.Printf("Password: %s\n", resp.Password)
		return
	case storage.Card:
		resp, err := client.GetCard(ctx, req)
		if err != nil {
			fmt.Println("Cant get your info!")
		}
		fmt.Printf("CardNumber: %s\n", resp.CardNumber)
		fmt.Printf("Holder: %s\n", resp.Holder)
		fmt.Printf("Date: %s\n", resp.Date)
		fmt.Printf("CVCcode: %s\n", resp.CVCcode)
	case storage.Text:
		resp, err := client.GetText(ctx, req)
		if err != nil {
			fmt.Println("Can't get your info!")
		}
		fmt.Printf("Text: %s\n", resp.Text)
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
	return ErrExitCLI
}
