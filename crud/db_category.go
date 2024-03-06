package crud

import "log"

// GetCatID get category by its ID from the Category table.
// Inputs: Cid - an integer representing the category ID to be retrieved
// Outputs: Category - the retrieved category, bool - a success flag indicating whether the category was found, error- any error that occurred during the query
func GetCatId(Cid int) (Category, bool, error) {
	rows, err := DB.Query("SELECT * FROM Category WHERE CatID=?;", Cid)
	var cat Category
	if err != nil {
		log.Printf("GetCatId() #1: %v", err)
		return cat, false, err
	}
	counter := 0
	for rows.Next() {
		counter += 1
		rows.Scan(&cat.CatID, &cat.UserID, &cat.Name, &cat.Color)
	}
	rows.Close()

	return cat, counter == 1, err
}

// CreateCategory inserts a new category into the Category table.
// It starts a transaction, prepares the SQL statement, executes the statement, and commits the transaction.
// Inputs: cat - a Category struct representing the category to be created
// Outputs:
// bool  - a success flag indicating whether the category creation was successful
// int64 - the ID of the newly created category
// error - any error that occurred during the transaction or statement execution
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
