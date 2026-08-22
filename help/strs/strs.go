package strs

import (
	"mime/multipart"
	_ "modernc.org/sqlite"
)

// structs
type User struct {
	UserId   int    `db:"userId"`
	UserName string `db:"userName" form:"UserName" validate:"required"`
	Password string `db:"password" form:"Password" validate:"required"`
	Access   int    `db:"access"   form:"Access"   validate:"number"`
}
 
type Project struct {
	ProjectId   int                     `db:"projectId"`
	UserId      int                     `db:"userId"`
	ProjectName string                  `db:"projectName" form:"ProjectName" validate:"required"`
	Description string                  `db:"description" form:"Description" validate:"required"`
	ImagePaths  string                  `db:"imagePaths"`
	Images      []*multipart.FileHeader `form:"Images" db:"-"`
}
 
type Rating struct {
	RatingId int    `db:"ratingId"`
	Name     string `db:"name" form:"Name" validate:"required"`
	Comment  string `db:"comment" form:"Comment"`
	Rating int `db:"rating" form:"Rating" validate:"required,gte=0,lte=50"`
}
 
type Member struct {
	MemberId          int                   `db:"memberId"`
	MemberName        string                `db:"memberName"        form:"MemberName" validate:"required"`
	MemberTitle       string                `db:"memberTitle"       form:"MemberTitle"`
	MemberDescription string                `db:"memberDescription" form:"MemberDescription"`
	MemberImagePath   string                `db:"memberImagePath"`
	MemberImage       *multipart.FileHeader `form:"MemberImage" db:"-" validate:"required"`
}
 
type Contact struct {
	Address  string `db:"address"  form:"Address"`
	Baladya  string `db:"baladya"  form:"Baladya"`
	Wilaya   string `db:"wilaya"   form:"Wilaya"`
	Email    string `db:"email"    form:"Email" validate:"omitempty,email"`
	Number   string `db:"number"   form:"Number"`
	Location string `db:"location" form:"Location"`
}
 
type Message struct {
	MessageId int    `db:"messageId"`
	FirstName string `db:"firstName" form:"FirstName" validate:"required"`
	LastName  string `db:"lastName"  form:"LastName" validate:"required"`
	Email     string `db:"email"     form:"Email" validate:"required,email"`
	Object    string `db:"object"    form:"Object" validate:"required"`
	Message   string `db:"message"   form:"Message" validate:"required"`
}
 
type JobApplication struct {
	ApId int    `db:"apId"`
	FirstName     string `db:"firstName" form:"FirstName" validate:"required"`
	LastName      string `db:"lastName"  form:"LastName" validate:"required"`
	Email         string `db:"email"     form:"Email" validate:"required,email"`
	Object        string `db:"object"    form:"Object" validate:"required"`
	Message       string `db:"message"   form:"Message" validate:"required"`
}
