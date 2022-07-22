# Library Management System

## Table of Contents
* [General Info](#general-information)
* [Technologies Used](#technologies-used)
* [Features](#features)

* [Usage](#usage)
* [Room for Improvement](#room-for-improvement)
* [Acknowledgements](#acknowledgements)
* [Contact](#contact)
* [License](#license)


## General Information
- Library Management System built using GOLang
- Data stored in DB in encoded format using GOB encoding
- Persistent application using BadgerDB
- Rest API implemented to accept JSON inputs
- Week(1-4) Final Project during SDE Internship at Sarva Labs



## Technologies Used
- GOLang - version 1.18.3 amd64
- BadgerDB - version v3.2103.2
- Gorilla Mux - version 1.8.0


## Features
Ready to use features:
- Add Book
- Add User
- Borrow book
- Return book
- Each user can borrow a maximum of 5 books
- Audiobook and eBook are only available in digital format
- Hardback and Paperback are only available in physical format
- Encyclopedia, Comic and Magazine are available in both digital and physical format
## Usage
Application Features:
- File can be downloaded and run using go by running the main package.
- Application will enter API mode on pressing 3 from main menu.
- API endpoints are available on port:10000
- Sample input provided in input.txt file to enter 7 books, 3 users and 1 borrow and 1 return into the system

### REST API JSON data formats:
#### ``1.BookType Enum``
Book type in following fields must be entered using the associated int value:
```
0. eBook
1. AudioBook
2. HardBack
3. PaperBack
4. Encyclopedia
5. Magazine
6. Comic
 
```
#### ``2.Book``
The object requires booktype, name, author and capacity in respective fields.
```
{
    "booktype":<int>(between 0 to 6),
    "name":<str>,
    "author":<str>,
    "capacity":<int>(To set >1 type must be digital),
    "borrowed":0
}
```
#### ``3.Member``
The object requires name and age in respective fields.
```
{
    "name":<str>,
    "age":<int>,
    "books":[]
}
```
#### ``4.Borrow/Return``
The object requires name, type(booktype) and bookname in respective fields.
```
{
    "name":<str>,
    "type":<int>(between 0 to 6),
    "bookname":<str>
}
```
### REST API end points:
The application has 4 API end points
#### **1. Add Book**
This functionality is available at the */book* endpoint of the application URL.
- Accepts only **POST** requests.
- Request JSON must follow the following format
    ```
    {
    "booktype":<int>(between 0 to 6),
    "name":<str>,
    "author":<str>,
    "capacity":<int>(To set >1 type must be digital),
    "borrowed":0
    }
    ```
The application returns any error or problem at server side while also returning the received json object to verify input.
#### **2. Register User**
This functionality is available at the */user* endpoint of the application URL.
- Accepts only **POST** requests.
- Request JSON must follow the following format
    ```
    {
    "name":<str>,
    "age":<int>,
    "books":[]
    }
    ```
  
The application returns any error or problem at server side while also returning the received json object to verify input.
#### **3. Borrow Book**
This functionality is available at the */borrow* endpoint of the application URL.
- Accepts only **POST** requests.
- Request JSON must follow the following format
    ```
    {
    "name":<str>,
    "type":<int>(between 0 to 6),
    "bookname":<str>
    }
    ```
The application returns any error or problem at server side while also returning the received json object to verify input.
#### **4. Return Book**
This functionality is available at the */return* endpoint of the application URL.
- Accepts only **POST** requests.
- Request JSON must follow the following format
    ```
    {
    "name":<str>,
    "type":<int>(between 0 to 6),
    "bookname":<str>
    }
    ```
The application returns any error or problem at server side while also returning the received json object to verify input.
## Room for Improvement
Room for improvement:
- Errors caused by API endpoints to be returned to client
- Improvement in efficiency using concurrent methods
- Improved read-write operations to reduce redundant operations when using API end points


## Acknowledgements

- This readme file was based on [this template](https://github.com/ritaly/README-cheatsheet/blob/master/README.md)
- The REST API with MUX implementation was inspired by [this tutorial](https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql).



## Contact
Created by [@sarthak815](https://www.github.com/sarthak815) - feel free to contact me!



## License
 This project is open source and available under the [MIT License](https://raw.githubusercontent.com/sarthak815/Library-Management-System/main/LICENSE). 

