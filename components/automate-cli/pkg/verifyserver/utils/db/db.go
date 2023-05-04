package db

import "database/sql"

type DBImpl struct{

}
type DB interface {
	InitPostgresDB(con string) error

}

func NewDBImpl()DB{
	return &DBImpl{}

}
func (di *DBImpl)InitPostgresDB(con string) error {

	db,err := sql.Open("postgres",con)
		if err != nil {
			return err
		}
	
		defer db.Close()
	
		resp := db.Ping()
		if resp != nil {
			return  err
		} 
		return nil
}