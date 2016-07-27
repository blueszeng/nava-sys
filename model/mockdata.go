package model

func MockUser()  []*User{
	users := []*User{
		{Name:"tom", Password: "2324"},
		{Name:"jip", Password: "2rwqf"},
		{Name:"tam", Password: "afafsd"},
		{Name:"satit", Password: "adfafq"},
		{Name:"somrod", Password: "af45"},
		{Name:"bee" , Password: "1qagg"},
	}

	for k, _ := range users{
		users[k].SetPass()
	}
	return users
}

