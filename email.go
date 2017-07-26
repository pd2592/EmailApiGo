package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	gomail "gopkg.in/gomail.v2"
)

type HotelData struct {
	hotelEmail    string
	hotelPassword string
	host          string
	port          int
}

type BookingData struct {
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	NumberOfRoom  int
	HotelID       int
}

// type FullDetail struct {
// 	booking BookingData
// 	hotel   HotelData
// }

func main() {

	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", EmailPage)
	// log.Fatal(http.ListenAndServe(":8080", router))

	EmailPage()

}

func EmailPage() {
	//hotelid := 1
	bookingid := 2
	emailtypeid := 2

	BookingDet := getBookingDet(bookingid)

	hotelid := BookingDet.HotelID
	//fmt.Println("Hotel Id is : ", hotelid)
	HotelCred := getHotelCred(hotelid)

	EmailTemplate := getEmailTemplate(emailtypeid)

	EmailAction(EmailTemplate, HotelCred, BookingDet)
	//fmt.Fprint(w, "mail sent")

}

func EmailAction(EmailTemplate string, HotelCredential HotelData, BookingDetail BookingData) {
	// ContentEmail := emaildet{
	// 	mid: r.FormValue("mId"),
	// 	//customerName:  r.FormValue("customerName"),
	// 	//phoneNumber:   r.FormValue("phoneNumber"),
	// 	customerEmail: r.FormValue("customerEmail"),
	// 	message:       r.FormValue("message"),
	// }
	from := HotelCredential.hotelEmail
	subject := "Testing"
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", BookingDetail.CustomerEmail)
	m.SetHeader("subject", subject)
	m.SetBody("text/html", getTemplate(EmailTemplate, BookingDetail))
	d := gomail.NewDialer(HotelCredential.host, HotelCredential.port, from, HotelCredential.hotelPassword)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	fmt.Println("mail sent")

}

func getHotelCred(hotelId int) HotelData {
	db, err := sql.Open("mysql", "root:@/user")
	checkErr(err)
	rows, err := db.Query("SELECT hotelEmail, hotelPassword, host, port FROM hotel_detail WHERE hotelId !=" + strconv.Itoa(hotelId))
	var HData HotelData
	for rows.Next() {
		err := rows.Scan(&HData.hotelEmail, &HData.hotelPassword, &HData.host, &HData.port)
		checkErr(err)
		fmt.Println("////"+HData.hotelEmail+"????"+HData.hotelPassword+"???? "+HData.host+"???? ", HData.port)

		HData = HotelData{
			hotelEmail:    HData.hotelEmail,
			hotelPassword: HData.hotelPassword,
			host:          HData.host,
			port:          HData.port,
		}

	}
	return HData

}

func getBookingDet(bookingId int) BookingData {

	db, err := sql.Open("mysql", "root:@/user")
	checkErr(err)

	err = db.Ping() //checking Db connection

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Ping to database successful, connection is still alive")

	Id := strconv.Itoa(bookingId)
	//fmt.Println(reflect.TypeOf(bookingId))
	//fmt.Println(Id)
	rows, err := db.Query("SELECT CustomerName, CustomerEmail, CustomerPhone, NumberOfRoom, HotelID FROM `booking_detail` WHERE BookingId =" + Id)

	var BData BookingData
	for rows.Next() {
		err := rows.Scan(&BData.CustomerName, &BData.CustomerEmail, &BData.CustomerPhone, &BData.NumberOfRoom, &BData.HotelID)
		checkErr(err)
		//fmt.Println("////" + BData.CustomerName + "????" + BData.CustomerEmail + "????" + BData.CustomerPhone + "????" + BData.NumberOfRoom + "???? " + BData.HotelID)

		BData = BookingData{
			CustomerName:  BData.CustomerName,
			CustomerEmail: BData.CustomerEmail,
			CustomerPhone: BData.CustomerPhone,
			NumberOfRoom:  BData.NumberOfRoom,
			HotelID:       BData.HotelID,
		}
		fmt.Println("Data from db : ", BData)

	}
	return BData

}

func getEmailTemplate(templateId int) string {
	db, err := sql.Open("mysql", "root:@/user")
	checkErr(err)
	var tempUrl string
	err = db.QueryRow("SELECT templateUrl FROM template_detail WHERE templateId =" + strconv.Itoa(templateId)).Scan(&tempUrl)
	//defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("template url is : ", tempUrl)
	return tempUrl
}

func getTemplate(file string, data BookingData) string {
	//function to make body of the email (template + Data)
	t, _ := template.ParseFiles(file)
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()

}

func checkErr(err error) {

	if err != nil {
		//fmt.Println(err)
		log.Fatal(err)
		os.Exit(1)
	}
}
