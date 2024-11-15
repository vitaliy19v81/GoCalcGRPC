package main

import (
	pb "GRPCADDER/pkg/api/proto" //"path/to/your/proto" // Укажите путь к сгенерированному proto-коду
	"context"
	"fmt"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net/http"
	"strconv"
	// "flag"
)

// Подключение к gRPC-серверу
var grpcClient pb.CalculatorClient

func init() {
	conn, err := grpc.Dial("192.168.31.214:50051", grpc.WithInsecure()) // Подключение к gRPC-серверу на порту 50051
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC-серверу: %v", err)
	}
	grpcClient = pb.NewCalculatorClient(conn)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/calculator", calculator)
	http.HandleFunc("/doCalc", doCalc)

	http.ListenAndServe("localhost:8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("cmd/client/tmpl/home.html")

	Marketing := struct {
		Message string
	}{
		Message: "Наше сообщение",
	}
	t.Execute(w, Marketing)
}

type dataForCalc struct {
	Answer        int
	Error         string
	IsAnswerExist bool
	Num1          int
	Num2          int
	Operation     string
}

func calculator(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("cmd/client/tmpl/calc.html")
	err := t.Execute(w, dataForCalc{Num1: 0, Num2: 0, Operation: "add"})
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func doCalc(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("cmd/client/tmpl/calc.html")
	r.ParseForm()
	data := dataForCalc{}
	var err error

	// Получаем числа и операцию из формы
	data.Num1, err = strconv.Atoi(r.FormValue("number1"))
	if err != nil {
		data.Error = "Ошибка в первом числе"
		t.Execute(w, data)
		return
	}

	data.Num2, err = strconv.Atoi(r.FormValue("number2"))
	if err != nil {
		data.Error = "Ошибка во втором числе"
		t.Execute(w, data)
		return
	}

	data.Operation = r.FormValue("operation")

	// Отправляем данные на gRPC-сервер
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем gRPC-запрос в зависимости от операции
	request := &pb.CalculationRequest{X: int32(data.Num1), Y: int32(data.Num2), Operation: data.Operation}
	response, err := grpcClient.Calculate(ctx, request)
	if err != nil {
		data.Error = fmt.Sprintf("Ошибка при вызове gRPC-сервера: %v", err)
		t.Execute(w, data)
		return
	}

	// Устанавливаем ответ в структуру и отправляем на страницу
	data.Answer = int(response.Result)
	data.Error = response.Error
	data.IsAnswerExist = true

	err = t.Execute(w, data)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

//package main
//
//import (
//	"GRPCADDER/pkg/api"
//	"context"
//	"flag"
//	"google.golang.org/grpc"
//	"log"
//	"strconv"
//)
//
//func main() {
//	flag.Parse()
//	args := flag.Args()
//	if len(args) < 2 {
//		log.Fatal("not enough arguments")
//	}
//
//	x, err := strconv.Atoi(args[0])
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	y, err := strconv.Atoi(args[1])
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	c := api.NewAdderClient(conn)
//	res, err := c.Add(context.Background(), &api.AddRequest{X: int32(x), Y: int32(y)})
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println(res.GetResult())
//}
