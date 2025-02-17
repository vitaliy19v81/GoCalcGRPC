package main

import (
	pb "GRPCADDER/pkg/api/proto" //"path/to/your/proto" // Укажите путь к сгенерированному proto-коду
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	// "flag"
)

// Подключение к gRPC-серверу
var grpcClient pb.CalculatorClient

func init() {
	// адрес удаленного компьютера (можно узнать на удаленном компьютере
	// или ip route get 1.1.1.1
	// или hostname -I
	// или ip -o -f inet addr show
	// или ip a
	// )

	// Подключение к gRPC-серверу на порту 50051
	conn, err := grpc.NewClient("192.168.31.214:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//conn, err := grpc.Dial("192.168.31.214:50051", grpc.WithInsecure()) // Устарело
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC-серверу: %v", err)
	}
	grpcClient = pb.NewCalculatorClient(conn)
}

// secretKey - ключ для подписи JWT
var secretKey = []byte("my-secret-key") //TODO перенести в секреты

func generateJWT(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // Токен на 1 час
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Запускать на клиенте
func main() {
	http.HandleFunc("/", calculator)
	http.HandleFunc("/doCalc", doCalc)

	http.ListenAndServe("localhost:8080", nil)
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

	// Получаем JWT-токен (должен быть сгенерирован заранее)
	token, _ := generateJWT("user123")

	// Создаём метаданные с токеном
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))

	// Отправляем данные на gRPC-сервер
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

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
