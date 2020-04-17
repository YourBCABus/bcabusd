package db

type Meta map[string]interface{}

type User struct {
	ID           string `pg:"type:uuid,unique,pk,notnull,default:gen_random_uuid()"`
	IsBot        bool   `pg:"is_bot,notnull,use_zero"`
	BotTokenHash []byte
	UserToken    []byte
	IsSuperadmin bool `pg:",notnull,default:false"`
	Meta         Meta
}

type AuthProvider struct {
	tableName struct{} `pg:"authproviders"`

	ID       string `pg:"type:uuid,unique,pk,notnull,default:gen_random_uuid()"`
	UserID   string `pg:"type:uuid,notnull"`
	Provider string `pg:"type:varchar(255),notnull"`
	Email    string
	Subject  string
	Meta     Meta
}
