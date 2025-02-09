package schema

// This is a generated file. DO NOT EDIT manually.

type Error struct {
	Code    int
	Message string
}

type User struct {
	ID       int
	Username string
	Obj      struct {
		Hoge string
		Fuga string
	}
}

type MutateUser struct {
	Username string
}

type Post struct {
	ID      int
	UserId  int
	Title   string
	Content string
}

type MutatePost struct {
	Content string
	Title   string
}
