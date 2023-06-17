package structs

type Post struct {
	ID            int
	Username      string
	Title         string
	Description   string
	CreationDate  string
	Tags          []string
	ImageFileName string
	Likes         int
	Dislikes      int
}

type Comment struct {
	ID           string
	Content      string
	PostID       int
	UserID       int
	Username     string
	CreationDate string
	Likes        int
	Dislikes     int
}

type User struct {
	ID       string
	Username string // Display the name of the user who is logged in
	LoggedIn bool
}

type ErrorMessage struct {
	Error   bool
	Message string
	Field1  string
	Field2  string
	Field3  []string
}

type Tag struct {
	ID   int
	Name string
}

type ForPage struct {
	Error    ErrorMessage
	User     User
	Posts    []Post
	Tags     []Tag
	Comments []Comment
}
