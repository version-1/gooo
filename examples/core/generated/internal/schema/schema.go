package schema

// This is a generated file. DO NOT EDIT manually.

type Error struct {
	Code    int
	Message string
}

type User struct {
	Username string
	Obj      struct {
		Hoge string
		Fuga string
	}
	ID int
}

type MutateUser struct {
	Username string
}

type Post struct {
	UserId  int
	Title   string
	Content string
	ID      int
}

type MutatePost struct {
	Title   string
	Content string
}
