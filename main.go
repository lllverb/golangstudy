package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Quiz struct {
	gorm.Model
	Question string
	Answer   string
	Choices  []Choice
}

type Choice struct {
	gorm.Model
	QuizId  int
	Text    string
	Correct int
}

//DB初期化
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInit）")
	}
	db.AutoMigrate(&Quiz{}, &Choice{})
	defer db.Close()
}

//DB追加
func create(question string, answer string, text string, correct int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInsert)")
	}
	db.Create(&Quiz{Question: question, Answer: answer, Choices: []Choice{{Text: text, Correct: correct}}})
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

func dbGetChoices(id int) []Choice {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetOne())")
	}
	var choices []Choice
	db.Where("quiz_id = ?", id).Find(&choices)
	db.Close()
	return choices
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//Index
	router.GET("/", func(ctx *gin.Context) {
		quizzes := dbGetAll()

		ctx.HTML(200, "index.html", gin.H{
			"quizzes": quizzes,
		})
	})
	// Sample
	// router.GET("/sample", func(ctx *gin.Context) {
	// 	quizzes := dbGetAll()
	// 	ctx.HTML(200, "base.html", gin.H{
	// 		"quizzes": quizzes,
	// 	})
	// })

	//Create
	router.POST("/new", func(ctx *gin.Context) {
		quiz := ctx.PostForm("quiz")
		answer := ctx.PostForm("answer")
		text := ctx.PostForm("text")
		c := ctx.PostForm("correct")
		correct, _ := strconv.Atoi(c)
		create(quiz, answer, text, correct)
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
		choices := dbGetChoices(id)
		ctx.HTML(200, "detail.html", gin.H{
			"quiz":    quiz,
			"choices": choices,
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
