package main

//Library structure
type Library struct {
	BooksBorrowed []Books
	Members       []Member
}

//Checks if the user wishing to borrow is registered and eligible
func checkUserValidity(name string, lib Library) (bool, *Member) {
	b := false         //Denotes user validity
	var member *Member //Used to return member object of user that wishes to borrow or return
	for i := range lib.Members {
		if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
			member = &lib.Members[i]
			b = true
			break
		}
	}
	if b { //used for out of index error handling
		if len(member.BooksBorrowed) == 5 {
			b = false
		}
	}
	return b, member
}

//Checks if the user wishing to return is eligible
func checkUserValidityReturn(name string, lib Library) (bool, *Member) {
	b := false         //Denotes user validity
	var member *Member //Used to return member object of user that wishes to borrow or return
	for i := range lib.Members {
		if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
			member = &lib.Members[i]
			b = true
			break
		}
	}
	return b, member
}

//checks if username is unique as it acts as primary key
func checkUserNameValidity(username string, lib Library) bool {
	b := true
	for i := range lib.Members {
		if username == lib.Members[i].Name {
			b = false
			break
		}
	}
	return b
}

//checks if book borrowed and present in directory
func checkUserBookValidity(bname string, lib Library, member Member) (bool, *Books) {
	bfound := false //Denotes book validity
	borrowed := false
	var bookFound *Books //Used to return book object that user wishes to borrow or return
	for i := range lib.BooksBorrowed {
		if lib.BooksBorrowed[i].B_Name == bname {
			bookFound = &lib.BooksBorrowed[i]
			bfound = true
			break
		}
	}
	if bfound { //used for out of index error handling
		for i := range member.BooksBorrowed { // checks if user has borrowed same book
			if member.BooksBorrowed[i].B_Name == bookFound.B_Name {
				borrowed = true
			}
		}
	}
	if !borrowed { //incase book not borrowed by user
		bfound = false
	}
	return bfound, bookFound
}
