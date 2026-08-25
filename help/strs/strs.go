package strs

import (
	"mime/multipart"
)

type User struct {
	UserId   int    `db:"userId"`
	UserName string `db:"userName" form:"UserName" validate:"required,min=3,max=50"`
	Password string `db:"password" form:"Password" validate:"required,min=6,max=12"`
	Access   string `db:"access"   form:"Access"   validate:"required,max=4"`
}

type Category struct {
	CategoryId   int                   `db:"categoryId"`
	CategoryName string                `db:"categoryName" form:"CategoryName" validate:"required,min=2,max=150"`
	ImagePath    string                `db:"imagePath"`
	Image        *multipart.FileHeader `db:"-" form:"Image" validate:"required"`
}

type Project struct {
	ProjectId   int                     `db:"projectId"`
	UserId      int                     `db:"userId"`
	CategoryId  int                     `db:"categoryId"`
	ProjectName string                  `db:"projectName" form:"ProjectName" validate:"required,min=2,max=150"`
	Description string                  `db:"description" form:"Description" validate:"required,min=5,max=5000"`
	ImagePaths  string                  `db:"imagePaths"`
	Images      []*multipart.FileHeader `db:"-"           form:"Images"      validate:"required,gt=0"`
	MImagePath  string                  `db:"mImagePath"`
	MImage      *multipart.FileHeader   `db:"-"           form:"MImage"      validate:"required"`
}

type Rating struct {
	RatingId int    `db:"ratingId"`
	Name     string `db:"name"     form:"Name"    validate:"required,min=2,max=60"`
	Comment  string `db:"comment"  form:"Comment" validate:"omitempty,max=1000"`
	Rating   int    `db:"rating"   form:"Rating"  validate:"required,gte=0,lte=50"`
}

type Member struct {
	MemberId          int                   `db:"memberId"`
	MemberName        string                `db:"memberName"        form:"MemberName"        validate:"required,min=2,max=100"`
	MemberTitle       string                `db:"memberTitle"       form:"MemberTitle"       validate:"omitempty,max=100"`
	MemberDescription string                `db:"memberDescription" form:"MemberDescription" validate:"omitempty,max=2000"`
	MemberImagePath   string                `db:"memberImagePath"`
	MemberImage       *multipart.FileHeader `db:"-"                 form:"MemberImage"       validate:"required"`
}

type Contact struct {
	Address  string `db:"address"  form:"Address"  validate:"omitempty,max=255"`
	Baladya  string `db:"baladya"  form:"Baladya"  validate:"omitempty,max=100"`
	Wilaya   string `db:"wilaya"   form:"Wilaya"   validate:"omitempty,max=100"`
	Email    string `db:"email"    form:"Email"    validate:"omitempty,email,max=100"`
	Number   string `db:"number"   form:"Number"   validate:"omitempty,max=30"`
	Location string `db:"location" form:"Location" validate:"omitempty,max=255"`
}

type Message struct {
	MessageId int    `db:"messageId"`
	FirstName string `db:"firstName" form:"FirstName" validate:"required,min=2,max=60"`
	LastName  string `db:"lastName"  form:"LastName"  validate:"required,min=2,max=60"`
	Email     string `db:"email"     form:"Email"     validate:"required,email,max=100"`
	Object    string `db:"object"    form:"Object"    validate:"required,min=2,max=150"`
	Message   string `db:"message"   form:"Message"   validate:"required,min=5,max=3000"`
}

type JobApplication struct {
	ApId      int    `db:"apId"`
	FirstName string `db:"firstName" form:"FirstName" validate:"required,min=2,max=60"`
	LastName  string `db:"lastName"  form:"LastName"  validate:"required,min=2,max=60"`
	Email     string `db:"email"     form:"Email"     validate:"required,email,max=100"`
	Object    string `db:"object"    form:"Object"    validate:"required,min=2,max=150"`
	Message   string `db:"message"   form:"Message"   validate:"required,min=5,max=3000"`
}

type Partenair struct {
	PartenairID   int64                 `db:"partenairId"   form:"partenairId"`
	PartenairName string                `db:"partenairName" form:"partenairName" validate:"required,min=2,max=100"`
	ImagePath     string                `db:"imagePath"`
	PartenairImg  *multipart.FileHeader `db:"-"             form:"PartenairImg"  validate:"required"`
}