package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

var credentails string = "../gsuiteServiceAccount.json"

var (
	method     = flag.String("m", "list", "方法名\nlist 列出域下的所有用户\nupload 创建用户\n")
	adminEmail = flag.String("a", "", "管理员账号 (默认为空)")
)

func getDirectoryService(adminEmail string) (*admin.Service, error) {
	ctx := context.Background()
	jsonCredentials, err := ioutil.ReadFile(credentails)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, admin.AdminDirectoryUserScope)
	if err != nil {
		return nil, fmt.Errorf("failed to JWTConfigFromJSON: %w", err)
	}

	config.Subject = adminEmail
	ts := config.TokenSource(ctx)

	srv, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("failed to create NewService: %w", err)
	}

	return srv, nil
}

func List(srv *admin.Service) {
	r, err := srv.Users.List().Domain("17996.com").OrderBy("email").Do()
	if err != nil {
		log.Panicf("failed to list user in domain: %v", err)
	}

	if len(r.Users) == 0 {
		fmt.Println("No user found.")
	} else {
		fmt.Println("Users:")
		for _, u := range r.Users {
			fmt.Printf("%s (%s)\n", u.PrimaryEmail, u.Name.FullName)
		}
	}
}

func main() {
	flag.Parse()
	srv, err := getDirectoryService(*adminEmail)
	if err != nil {
		log.Panicf("failed to get service: %v", err)
	}
	switch *method {
	case "list":
		List(srv)
	case "create":
		fmt.Println("coming soon...")
	default:
		log.Panic("method error")
	}
}
