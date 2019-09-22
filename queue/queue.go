package queue

type Queue struct {
	path string

	queueUsers    []*User
	queueUsersMap map[string]*User
	skipUsers     []*User
	skipUsersMap  map[string]*User
	usersAm       int
}

type User struct {
	TgNick      string
	ChatID      int64
	PriorityNum int
}

func NewQueue(path string) (*Queue, error) {
	//add loading backup
	return &Queue{
		path:          path,
		queueUsers:    make([]*User, 0),
		queueUsersMap: make(map[string]*User),
		skipUsers:     make([]*User, 0),
		skipUsersMap:  make(map[string]*User),
	}, nil
}

func (q *Queue) Next() *User {
	if len(q.queueUsers) != 0 {
		return q.queueUsers[0]
	}
	return nil
}

func (q *Queue) All() []*User {
	return q.queueUsers
}

func (q *Queue) Add(u User) {
	if _, ok := q.queueUsersMap[u.TgNick]; ok {
		return
	}
	user := u
	q.queueUsersMap[user.TgNick] = &user
	q.queueUsers = append(q.queueUsers, &user)
}

// func (q *Queue) Skip(tgName string) {
// 	if _, ok := q.queueUsersMap[u.TgNick]; ok {
// 		return
// 	}
// 	user := u
// 	q.queueUsersMap[user.TgNick] = &user
// 	q.queueUsers = append(q.queueUsers, &user)
// }

func (q *Queue) Save() {

}

const (
	UserNotRegisterd = iota
	UserInQueue
	UserSkipping
)

func (q *Queue) GetUser(tgNick string) (*User, int) {
	if user, ok := q.queueUsersMap[tgNick]; ok {
		return user, UserInQueue
	}
	if user, ok := q.skipUsersMap[tgNick]; ok {
		return user, UserSkipping
	}
	return &User{TgNick: tgNick}, UserNotRegisterd
}
