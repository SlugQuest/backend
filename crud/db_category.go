package crud

import "log"

// Get category by ID
func GetCatId(Cid int) (Category, bool, error) {
	rows, err := DB.Query("SELECT * FROM Category WHERE CatID=?;", Cid)
	var cat Category
	if err != nil {
		log.Println(err)
		return cat, false, err
	}
	counter := 0
	for rows.Next() {
		counter += 1
		log.Println(counter)
		rows.Scan(&cat.CatID, &cat.UserID, &cat.Name, &cat.Color)
		log.Println("finding")
	}
	rows.Close()
	log.Println("done finding")
	log.Println(counter)
	return cat, counter == 1, err
}

func CreateCategory(cat Category) (bool, int64, error) {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		log.Println("CreateCat(): breaky 1")
		return false, -1, err
	}
	defer tx.Rollback() //abort transaction if error

	//preparing statement to prevent SQL injection issues
	stmt, err := tx.Preparex("INSERT INTO Category ( UserID, Name, Color) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("cat(): breaky 2", err)
		return false, -1, err
	}

	defer stmt.Close() // Defer the closing of SQL statement to ensure it closes once the function completes
	res, err := stmt.Exec(cat.UserID, cat.Name, cat.Color)

	if err != nil {
		log.Println("Createcat(): breaky 3 ", err)
		return false, -1, err
	}

	catID, err := res.LastInsertId()
	if err != nil {
		log.Println("Createcat(): breaky 4 ", err)
		return false, -1, err
	}
	tx.Commit() //commit transaction to database

	return true, catID, nil
}
