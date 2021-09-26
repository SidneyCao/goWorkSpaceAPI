package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

var credentails string = "../gsuiteServiceAccount.json"
var password string = "../createRandPwd.sh"

var (
	method       = flag.String("m", "list", "方法名\nlist 列出域下的所有用户\ncreate 创建用户\n")
	adminUser    = flag.String("a", "", "管理员账号 (默认为空)")
	firstName    = flag.String("f", "", "全名 (默认为空)")
	lastName     = flag.String("l", "", "姓氏 (默认为空)")
	primaryEmail = flag.String("p", "", "主邮箱 (默认为空)")
	group        = flag.String("g", "", "分组 (默认为空)")
	OU           = flag.String("o", "", "组织名 (默认为空)")
	domain       = flag.String("d", "", "域名 (默认为空)")
)

func getDirectoryService(adminUser string, mod string) (*admin.Service, error) {
	ctx := context.Background()
	jsonCredentials, err := ioutil.ReadFile(credentails)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var config *jwt.Config

	switch mod {
	case "user":
		config, err = google.JWTConfigFromJSON(jsonCredentials, admin.AdminDirectoryUserScope)
		if err != nil {
			return nil, fmt.Errorf("failed to JWTConfigFromJSON: %w", err)
		}
	case "group":
		config, err = google.JWTConfigFromJSON(jsonCredentials, admin.AdminDirectoryGroupScope)
		if err != nil {
			return nil, fmt.Errorf("failed to JWTConfigFromJSON: %w", err)
		}
	default:
		log.Panicf("mod %s error ", mod)
	}

	config.Subject = adminUser + "@" + *domain
	ts := config.TokenSource(ctx)

	srv, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("failed to create NewService: %w", err)
	}

	return srv, nil
}

func listUser(srv *admin.Service) {
	r, err := srv.Users.List().Domain(*domain).OrderBy("email").Do()
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

func createUser(srv *admin.Service) {

	//生成随机密码
	cmd := exec.Command("sh", password)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panicf("failed to get password stdout: %v", err)
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Panicf("failed to create password: %v", err)
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Panicf("failed to read stdout: %v", err)
	}
	password := strings.Split(string(opBytes), "\n")[0]

	//通过json unmarshal
	//创建 admin.user struct
	primaryEmail := *primaryEmail + "@" + *domain
	userJson := fmt.Sprintf(`{"name":{"GivenName":"%s","FamilyName":"%s"},"primaryEmail":"%s","Password":"%s","ChangePasswordAtNextLogin":false,"OrgUnitPath":"%s"}`, *firstName, *lastName, primaryEmail, password, *OU)
	userByte := []byte(userJson)
	u := admin.User{}
	err = json.Unmarshal(userByte, &u)
	if err != nil {
		log.Panicf("json unmarshal error: %v", err)
	}

	//创建用户
	_, err = srv.Users.Insert(&u).Do()
	if err != nil {
		log.Panicf("failed to create user %s: %v", u.Name.FamilyName, err)
	}
	log.Printf("用户创建成功！用户名: %s   密码: %s", u.PrimaryEmail, u.Password)
}

func update(srv *admin.Service) {
	//添加用户到组
	m := admin.Member{}
	m.Email = *primaryEmail + "@" + *domain
	group := *group + "@" + *domain
	_, err := srv.Members.Insert(group, &m).Do()
	if err != nil {
		log.Panicf("failed to add user %s to group %s: %v", m.Email, group, err)
	}
	log.Printf("用户 %s 成功添加到组 %s", m.Email, group)

}

func main() {
	//解析参数
	flag.Parse()

	userSrv, err := getDirectoryService(*adminUser, "user")
	if err != nil {
		log.Panicf("failed to get user service: %v", err)
	}

	groupSrv, err := getDirectoryService(*adminUser, "group")
	if err != nil {
		log.Panicf("failed to get group service: %v", err)
	}

	switch *method {
	case "list":
		listUser(userSrv)
	case "create":
		createUser(userSrv)
		update(groupSrv)
	default:
		log.Panic("method error")
	}

}
