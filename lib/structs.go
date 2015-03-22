package lib

type Commit struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	URL     string `json:"url"`
	Author  Person `json:"author"`
}

type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Watchers    int    `json:"watchers"`
	Owner       Person `json:"owner"`
	Private     bool   `json:"private"`
}

type Person struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type CommitHook struct {
	Secret     string     `json:"secret"`
	Ref        string     `json:"ref"`
	Commits    []Commit   `json:"commits"`
	Repository Repository `json:"repository"`
	Pusher     Person     `json:"pusher"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	CompareURL string     `json:"compare_url"`
}
