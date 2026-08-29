package strs

import (
	"mime/multipart"
)

type User struct {
	UserId   int    `db:"userId"   form:"UserId"`
	UserName string `db:"userName" form:"UserName" validate:"required,min=3,max=50"`
	Password string `db:"password" form:"Password" validate:"required,min=6,max=12"`
	Access   string `db:"access"   form:"Access"   validate:"max=4"`
}

type Category struct {
	CategoryId   int                   `db:"categoryId"   form:"CategoryId"`
	CategoryName string                `db:"categoryName" form:"CategoryName" validate:"required,min=2,max=150"`
	ImagePath    string                `db:"imagePath"    form:"ImagePath"`
	Image        *multipart.FileHeader `db:"-"            form:"Image"`
}

type Project struct {
	ProjectId   int                     `db:"projectId"    form:"ProjectId"`
	UserId      int                     `db:"userId"       form:"UserId"`
	CategoryId  int                     `db:"categoryId"   form:"CategoryId"`
	ProjectName string                  `db:"projectName"  form:"ProjectName" validate:"required,min=2,max=150"`
	Description string                  `db:"description"  form:"Description" validate:"required,min=5,max=5000"`
	ImagePaths  string                  `db:"imagePaths"   form:"ImagePaths"`
	Images      []*multipart.FileHeader `db:"-"            form:"Images"`
	MImagePath  string                  `db:"mImagePath"   form:"MImagePath"`
	MImage      *multipart.FileHeader   `db:"-"            form:"MImage"`
}

type ProjectMod struct {
	ProjectId   int                     `db:"projectId"    form:"ProjectId"`
	UserId      int                     `db:"userId"       form:"UserId"`
	CategoryId  int                     `db:"categoryId"   form:"CategoryId"`
	ProjectName string                  `db:"projectName"  form:"ProjectName" validate:"min=2,max=150"`
	Description string                  `db:"description"  form:"Description" validate:"min=5,max=5000"`
	ImagePaths  string                  `db:"imagePaths"   form:"ImagePaths"`
	Images      []*multipart.FileHeader `db:"-"            form:"Images"`
	MImagePath  string                  `db:"mImagePath"   form:"MImagePath"`
	MImage      *multipart.FileHeader   `db:"-"            form:"MImage"`
}

type Rating struct {
	RatingId int    `db:"ratingId" form:"RatingId"`
	Name     string `db:"name"     form:"Name"    validate:"required,min=2,max=60"`
	Comment  string `db:"comment"  form:"Comment" validate:"omitempty,max=1000"`
	Rating   int    `db:"rating"   form:"Rating"  validate:"required,gte=0,lte=5"`
}

type Member struct {
	MemberId          int                   `db:"memberId"          form:"MemberId"`
	MemberName        string                `db:"memberName"        form:"MemberName"        validate:"required,min=2,max=100"`
	MemberTitle       string                `db:"memberTitle"       form:"MemberTitle"       validate:"omitempty,max=100"`
	MemberDescription string                `db:"memberDescription" form:"MemberDescription" validate:"omitempty,max=2000"`
	MemberImagePath   string                `db:"memberImagePath"   form:"MemberImagePath"`
	MemberImage       *multipart.FileHeader `db:"-"                 form:"MemberImage"`
}

type Contact struct {
	Address   string                 `db:"address"   form:"Address"   validate:"omitempty,max=255"`
	Baladya   string                 `db:"baladya"   form:"Baladya"   validate:"omitempty,max=100"`
	Wilaya    string                 `db:"wilaya"    form:"Wilaya"    validate:"omitempty,max=100"`
	Email     string                 `db:"email"     form:"Email"     validate:"omitempty,email,max=100"`
	Number    string                 `db:"number"    form:"Number"    validate:"omitempty,max=30"`
	Fax       string                 `db:"fax"       form:"Fax"       validate:"omitempty,max=30"`
	Location  string                 `db:"location"  form:"Location"  validate:"omitempty,max=255"`
	ImagePath string                 `db:"image"     form:"ImagePath"`
	Image     *multipart.FileHeader  `db:"-"         form:"Image"`
}

type Message struct {
	MessageId int    `db:"messageId" form:"MessageId"`
	FirstName string `db:"firstName" form:"FirstName" validate:"required,min=2,max=60"`
	LastName  string `db:"lastName"  form:"LastName"  validate:"required,min=2,max=60"`
	Email     string `db:"email"     form:"Email"     validate:"required,email,max=100"`
	Object    string `db:"object"    form:"Object"    validate:"required,min=2,max=150"`
	Message   string `db:"message"   form:"Message"   validate:"required,min=5,max=3000"`
}

type JobApplication struct {
	ApId      int    `db:"apId"      form:"ApId"`
	FirstName string `db:"firstName" form:"FirstName" validate:"required,min=2,max=60"`
	LastName  string `db:"lastName"  form:"LastName"  validate:"required,min=2,max=60"`
	Email     string `db:"email"     form:"Email"     validate:"required,email,max=100"`
	Object    string `db:"object"    form:"Object"    validate:"required,min=2,max=150"`
	Message   string `db:"message"   form:"Message"   validate:"required,min=5,max=3000"`
	Software  string `db:"software"  form:"Software"  validate:"required,min=5,max=3000"`
	Cv        *multipart.FileHeader `form:"cv"`
	CvPath    string             
	Letter    *multipart.FileHeader `form:"letter"`
	LetterPath string
}

type Partenair struct {
	PartenairID   int64                 `db:"partenairId"   form:"PartenairId"`
	PartenairName string                `db:"partenairName" form:"PartenairName" validate:"required,min=2,max=100"`
	ImagePath     string                `db:"imagePath"     form:"ImagePath"`
	PartenairImg  *multipart.FileHeader `db:"-"             form:"PartenairImg"  validate:"required"`
}

type Software struct {
	SoftwareId   int    `db:"softwareId"   form:"SoftwareId"`
	SoftwareName string `db:"softwareName" form:"SoftwareName" validate:"required,min=2,max=150"`
	Required     int    `db:"required"     form:"Required"     validate:"oneof=0 1"`
}

type CanNot struct {
	CanNotId  int  `db:"canNotId"   form:"CanNotId"`
	CandidatId int `db:"candidatId" form:"CandidatId"`
}

type MessNot struct {
    MessNotId int `db:"messNotId" form:"MessNotId"`
    MessageId int `db:"messageId" form:"MessageId"`
}