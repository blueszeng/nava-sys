package models

func MockUser()  []*User{
	users := []*User{
		{Name:"tom"},
		{Name:"jip"},
		{Name:"tam"},
		{Name:"satit"},
		{Name:"somrod"},
		{Name:"bee" },
	}
	passwords := []string{
		"1234",
		"1234",
		"23rsafasf",
		";alsjdfl",
		"a;dlfjka",
		"abc",
	}
	for k, _ := range users{
		users[k].SetPass(passwords[k])
	}
	return users
}

