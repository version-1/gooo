package fixtures

import "time"

type User struct {
	ID           int       `json:"id" gooo:"primary_key,immutable"`
	Username     string    `json:"username" gooo:"unique"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	Timezone     string    `json:"timezone"`
	TimeDiff     int       `json:"time_diff"`
	CreatedAt    time.Time `json:"created_at" gooo:"immutable"`
	UpdatedAt    time.Time `json:"updated_at" gooo:"immutable"`

	Profile *Profile `json:"profile" gooo:"association"`
	Posts   []Post   `json:"posts" gooo:"association"`
}

type Post struct {
	ID        int       `json:"id" gooo:"primary_key,immutable"`
	UserID    int       `json:"user_id" gooo:"index"`
	Title     string    `json:"title"`
	Body      string    `json:"body" gooo:"type=text"`
	CreatedAt time.Time `json:"created_at" gooo:"immutable"`
	UpdatedAt time.Time `json:"updated_at" gooo:"immutable"`

	User  User   `json:"user" gooo:"association"`
	Likes []Like `json:"likes" gooo:"association"`
}

type Profile struct {
	ID        int       `json:"id" gooo:"primary_key,immutable"`
	UserID    int       `json:"user_id" gooo:"index"`
	Bio       string    `json:"bio" gooo:"type=text"`
	CreatedAt time.Time `json:"created_at" gooo:"immutable"`
	UpdatedAt time.Time `json:"updated_at" gooo:"immutable"`
}

type Like struct {
	ID           int       `json:"id" gooo:"primary_key,immutable"`
	LikeableID   int       `json:"likeable_id" gooo:"index"`
	LikeableType string    `json:"likeable_type" gooo:"index"`
	CreatedAt    time.Time `json:"created_at" gooo:"immutable"`
	UpdatedAt    time.Time `json:"updated_at" gooo:"immutable"`
}
