package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/go-playground/validator.v8"
)

type Quiz struct {
	gorm.Model
	Question    string `validate:"required"`
	Explanation string `validate:"required"`
	Choices     []Choice
}

type Choice struct {
	gorm.Model
	QuizId  int    `validate:"required"`
	Text    string `validate:"required"`
	Correct int    `validate:"required"`
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

var validate *validator.Validate

//DB追加
func create(question string, explanation string, text1 string, correct1 int, text2 string, correct2 int, text3 string, correct3 int, text4 string, correct4 int) {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	quiz := &Quiz{Question: question, Explanation: explanation, Choices: []Choice{{Text: text1, Correct: correct1}, {Text: text2, Correct: correct2}, {Text: text3, Correct: correct3}, {Text: text4, Correct: correct4}}}
	errs := validate.Struct(quiz)
	if errs != nil {

		fmt.Println(errs) // output: Key: "User.Age" Error:Field validation for "Age" failed on the "lte" tag
		//	                         Key: "User.Addresses[0].City" Error:Field validation for "City" failed on the "required" tag
		// err := errs.(validator.ValidationErrors)
		// fmt.Println(err.Field) // output: City
		// fmt.Println(err.Tag)   // output: required
		// fmt.Println(err.Kind)  // output: string
		// fmt.Println(err.Type)  // output: string
		// fmt.Println(err.Param) // output:
		// fmt.Println(err.Value) // output:

		// from here you can create your own error messages in whatever language you wish
		return
	}
	if err != nil {
		panic("データベース開けず！（create)")
	}
	db.Create(&quiz)
	defer db.Close()
}

// //DB更新
func dbUpdate(id int, question string, explanation string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbUpdate)")
	}
	var quiz Quiz
	db.First(&quiz, id)
	quiz.Question = question
	quiz.Explanation = explanation
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
	db.Preload("Choices").Find(&quizzes)
	// fmt.Println(quizzes.Choices)
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
	db.Preload("Choices").Find(&quiz)
	fmt.Println(quiz)
	db.Close()
	return quiz
}

func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	// Index
	router.GET("/", func(ctx *gin.Context) {
		quizzes := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{
			"quizzes": quizzes,
		})
	})
	// New
	router.GET("/new", func(ctx *gin.Context) {
		ctx.HTML(200, "new.html", gin.H{})
	})

	//Create
	router.POST("/new", func(ctx *gin.Context) {
		quiz := ctx.PostForm("quiz")
		explanation := ctx.PostForm("explanation")
		text1 := ctx.PostForm("text1")
		c1 := ctx.PostForm("correct")
		correct1, _ := strconv.Atoi(c1)
		text2 := ctx.PostForm("text2")
		c2 := ctx.PostForm("correct")
		correct2, _ := strconv.Atoi(c2)
		text3 := ctx.PostForm("text3")
		c3 := ctx.PostForm("correct")
		correct3, _ := strconv.Atoi(c3)
		text4 := ctx.PostForm("text4")
		c4 := ctx.PostForm("correct")
		correct4, _ := strconv.Atoi(c4)
		create(quiz, explanation, text1, correct1, text2, correct2, text3, correct3, text4, correct4)
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
		ctx.HTML(200, "detail.html", gin.H{
			"quiz": quiz,
		})
		fmt.Println(quiz.Question)
	})

	// Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		question := ctx.PostForm("question")
		explanation := ctx.PostForm("explanation")
		dbUpdate(id, question, explanation)
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
	// render html template
	// router.Use(render.Renderer(render.Options{
	// 	Funcs: []template.FuncMap{ // レンダラにテンプレート関数を登録します。
	// 		{
	// 			"add": func(a, b int) int { return a + b },
	// 			"sub": func(a, b int) int { return a - b },
	// 			"mul": func(a, b int) int { return a * b },
	// 			"div": func(a, b int) int { return a / b },
	// 		},
	// 	},
	// }))
	// ExampleTemplateCalculator()
	router.Run()
}

// func ExampleTemplateCalculator(w http.ResponseWriter, r *http.Request) {
// 	funcMap := template.FuncMap{
// 		"add": func(a, b int) int { return a + b },
// 		"sub": func(a, b int) int { return a - b },
// 		"mul": func(a, b int) int { return a * b },
// 		"div": func(a, b int) int { return a / b },
// 	}
// 	tp := template.Must(template.New("index.html").Funcs(funcMap).Parse("index.html"))
// 	err := tp.Execute(w, params)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Output:
// 	// 10
// }
