package main

import (
	"strconv"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// type Todo struct {
// 	gorm.Model
// 	Text   string
// 	Status string
// }

type Quiz struct {
	gorm.Model
	Question string
	Answer   string
}

//DB初期化
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInit）")
	}
	db.AutoMigrate(&Quiz{})
	defer db.Close()
}

//DB追加
func create(question string, answer string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInsert)")
	}
	db.Create(&Quiz{Question: question, Answer: answer})
	defer db.Close()
}

// //DB更新
func dbUpdate(id int, question string, answer string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbUpdate)")
	}
	var quiz Quiz
	db.First(&quiz, id)
	quiz.Question = question
	quiz.Answer = answer
	db.Save(&quiz)
	db.Close()
}

// //DB削除
// func dbDelete(id int) {
// 	db, err := gorm.Open("sqlite3", "test.sqlite3")
// 	if err != nil {
// 		panic("データベース開けず！（dbDelete)")
// 	}
// 	var todo Todo
// 	db.First(&todo, id)
// 	db.Delete(&todo)
// 	db.Close()
// }

//DB全取得
func dbGetAll() []Quiz {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetAll())")
	}
	var quizzes []Quiz
	db.Order("created_at desc").Find(&quizzes)
	db.Close()
	return quizzes
}

//DB一つ取得
func dbGetOne(id int) Quiz {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetOne())")
	}
	var quiz Quiz
	db.First(&quiz, id)
	db.Close()
	return quiz
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//Index
	router.GET("/", func(ctx *gin.Context) {
		quizzes := dbGetAll()

		html := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
		router.SetHTMLTemplate(html)
		ctx.HTML(200, "base.html", gin.H{
			"quizzes": quizzes,
		})
	})
	// Sample
	router.GET("/sample", func(ctx *gin.Context) {
		quizzes := dbGetAll()

		html := template.Must(template.ParseFiles("templates/base.html", "templates/sample.html"))
		router.SetHTMLTemplate(html)
		ctx.HTML(200, "base.html", gin.H{
			"quizzes": quizzes,
		})
	})

	//Create
	router.POST("/new", func(ctx *gin.Context) {
		quiz := ctx.PostForm("quiz")
		answer := ctx.PostForm("answer")
		create(quiz, answer)
		ctx.Redirect(302, "/")
	})

	// Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		quiz := dbGetOne(id)

		html := template.Must(template.ParseFiles("templates/base.html", "templates/detail.html"))
		router.SetHTMLTemplate(html)
		ctx.HTML(200, "base.html", gin.H{
			"quiz": quiz,
		})
	})

	// Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		question := ctx.PostForm("question")
		answer := ctx.PostForm("answer")
		dbUpdate(id, question, answer)
		ctx.Redirect(302, "/")
	})

	// 削除確認
	// router.GET("/delete_check/:id", func(ctx *gin.Context) {
	// 	n := ctx.Param("id")
	// 	id, err := strconv.Atoi(n)
	// 	if err != nil {
	// 		panic("ERROR")
	// 	}
	// 	todo := dbGetOne(id)
	// 	ctx.HTML(200, "delete.html", gin.H{"todo": todo})
	// })

	// Delete
	// router.POST("/delete/:id", func(ctx *gin.Context) {
	// 	n := ctx.Param("id")
	// 	id, err := strconv.Atoi(n)
	// 	if err != nil {
	// 		panic("ERROR")
	// 	}
	// 	dbDelete(id)
	// 	ctx.Redirect(302, "/")

	// })

	router.Run()
}
