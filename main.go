package main

func main() {
	app := App{router: SetUpRouter()}
	app.Run(":8080")
}
