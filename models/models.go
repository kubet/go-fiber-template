package models

type Product struct {
	Id    uint   `json:"id" gorm:"primarykey;not null;default:null`
	Name  string `json:"name" gorm:"text;not null;default:null`
	Price uint   `json:"price" gorm:"number;not null;default:null`
}
