package model


func MockUsers()  []*User{
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

//func MockPerson() []*Person {
//	persons := []*Person{
//		Person{
//			Users: []User{
//				{Name: "tom"},
//				{Name:"tom2"},
//			},
//			Jobs: []Job{
//				{JobID: 1, OrgID: 1, NameTH: "กรรมการผู้จัดการ"},
//				{JobID: 1, OrgID: 2, NameTH: "ผู้อำนวยการขาย"},
//			},
//			Emails: []Email{
//				{Email: "tom@nopadol.com"},
//				{Email: "mrtomyum@gmail.com"},
//			},
//			FirstName: "Kasem",
//			LastName: "Arnontavilas",
//			BirthDate: time.Date(1974,10,4,0,0,0,0,time.UTC),
//			CitizenID:"3509901377838",
//		},
//
//		Person{
//			Users: []User{{Name: "jip"}},
//			Jobs: []Job{
//				{JobID: 1, OrgID: 1, NameTH: "ผู้อำนวยการฝ่ายขาย", NameEN: "Sales Director"},
//				{JobID: 1, OrgID: 3, NameTH: "ผู้อำนวยการฝ่ายบริหารสินค้า", NameEN: "Merchandise Director"},
//			},
//			Emails: []Email{
//				{Email: "jip@nopadol.com"},
//				{Email: "jipjiraporn@gmail.com"},
//			},
//			FirstName: "Jiraporn",
//			LastName: "Arnontavilas",
//			BirthDate: time.Date(1976,8,12,0,0,0,0,time.UTC),
//		},
//
//		Person{
//			Users: []User{{Name: "tam"}},
//			Jobs: []Job{
//				{NameTH: "ผู้อำนวยการสาขา2", NameEN: "Sales Director Branch 2"},
//			},
//			FirstName: "Tanan",
//			LastName: "Arnontavilas",
//		},
//		Person{
//			Users: []User{{Name: "noi"}},
//			Jobs: []Job{
//				{JobID:3, OrgID: 9, NameTH: "ผู้จัดการฝ่าย IT", NameEN: "IT Manager"},
//			},
//			FirstName: "Satit",
//			LastName: "Chomwattana",
//		},
//		Person{
//			Users: []User{
//				{Name: "bee"},
//			},
//			Jobs: []Job{
//				{JobID: 4, OrgID: 9, NameTH: "เจ้าหน้าที่ปฏิบัติการสารสนเทศ", NameEN: "MIS"},
//			},
//			Emails: []Email{
//				{Email: "mis@nopadol.com"},
//			},
//			FirstName: "เอกชัย",
//			LastName:  "จันตะไพ",
//			BirthDate: time.Date(1984, 5, 5, 0, 0, 0, 0, time.UTC),
//			CitizenID: "3509901371234",
//		},
//	}
//	return persons
//}