package structures


// import "go.mongodb.org/mongo-driver/bson/primitive"

type Questions struct{
	QuestionsText string  `json:"questiontext"`
	QuestionsTitle string `json:"questiontitle"`
	QuestionsAuthorName string `json:"questionauthorname"`
}
type Signin struct{
	Userid string `bson:"_id"`
	Permissions string `json:"permissions"`
	Phone int32 `json:"phone"`
	Password string `json:"password"`
}
type Signup struct{
	Userid string `bson:"_id"`
	Phone int32 `json:"phone"`
	Password string `json:"password"`
	Email string `json:"email"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Lastname string `json:"lastname"`
	Birth string `json:"birth"`
	Gender string `json:"gender"`
	Disabilaties string `json:"disabilaties"`
	Blood string `json:"blood"`
	Adress string `json:"adress"`
	Workplace string `json:"workplace"`
	ImgUrl string `json:"imgurl"`
	Permissions string `json:"permissions"`
}
type SignupDoctor struct{
	Userid string `bson:"_id"`
	Phone int32 `json:"phone"`
	Password string `json:"password"`
	Email string `json:"email"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Lastname string `json:"lastname"`
	Position string `json:"position"`
	Biography string `json:"biography"`
	Adress string `json:"adress"`
	Permissions string `json:"permissions"`
	ImgUrl string `json:"imgurl"`
	History []History
}
type History struct{
	Year string `json:"year"`
	Position string `json:"position"`
}
type Admin struct{
	Userid string `bson:"_id"`
	Phone int32 `json:"phone"`
	Password string `json:"password"`
	Email string `json:"email"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Lastname string `json:"lastname"`	
	Permissions string `json:"permissions"`
	ImgUrl string `json:"imgurl"`
}
type Views struct{
	Id string `bson:"_id"`
	ClientFLSname string `json:"clientFLSname"`
	ClientId string `json:"clientid"`
	DoctorId string `json:"doctorid"`
	Sickness string `json:"sickness"`
	Date string `json:"date"`
	Phone int32 `json:"phone"`
}
type File struct{
	Id string `json:"id" bson:"_id"`
	ClientId string `json:"clientid"`
	DoctorId string `json:"doctorid"`
	ClientFLSname string `json:"clientFLSname"`
	DoctorFLSname string `json:"doctorFLSname"`
	Title string `json:"title"`
	ImgUrl string `json:"imgurl"`
}
type GlobeStruct struct{
	// Userid string`bson:"_id"`
	QuestionsText string  `json:"questiontext"`
	QuestionsTitle string `json:"questiontitle"`
	QuestionsAuthorName string `json:"questionauthorname"`
	// ----------------------------------------------------
	Userid string `bson:"_id"`
	Phone int32 `json:"phone"`
	Password string `json:"password"`
	Email string `json:"email"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Lastname string `json:"lastname"`
	Birth string `json:"birth"`
	Gender string `json:"gender"`
	Disabilaties string `json:"disabilaties"`
	Adress string `json:"adress"`
	Workplace string `json:"workplace"`
	Permissions string `json:"permissions"`
	// ----------------------------------------------------
	Position string `json:"position"`
	Biography string `json:"biography"`
	ImgUrl string `json:"imgurl"`
	// History []History
}