package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type ParkingTicket struct {
	Index         int
	VehicleNumber string
	VehicleType   string
	EntryTime     time.Time
	ExitTime      time.Time
	Fee           float64
}

type User struct {
	Username string
	Password string
}

type ParkingSystem struct {
	Tickets []ParkingTicket
	Users   []User
}

const (
	FeePerHourCar  = 5000.0
	FeePerHourBike = 2000.0
)

func (ps *ParkingSystem) AddTicket(ticket ParkingTicket) {
	ps.Tickets = append(ps.Tickets, ticket)
	ps.updateTicketIndices()
}

func (ps *ParkingSystem) ModifyTicket(index int, ticket ParkingTicket) {
	if index >= 0 && index < len(ps.Tickets) {
		ps.Tickets[index] = ticket
		ps.updateTicketIndices()
	} else {
		fmt.Println("Tiket Tidak Ditemukan!")
	}
}

func (ps *ParkingSystem) DeleteTicket(index int) {
	if index >= 0 && index < len(ps.Tickets) {
		ps.Tickets = append(ps.Tickets[:index], ps.Tickets[index+1:]...)
		ps.updateTicketIndices()
	} else {
		fmt.Println("Tiket Tidak Ditemukan!")
	}
}

func (ps *ParkingSystem) AddUser(user User) {
	ps.Users = append(ps.Users, user)
}

func (ps *ParkingSystem) Authenticate(username, password string) bool {
	for _, user := range ps.Users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

func (ps *ParkingSystem) SearchTicket(vehicleNumber string) *ParkingTicket {
	for _, ticket := range ps.Tickets {
		if ticket.VehicleNumber == vehicleNumber {
			return &ticket
		}
	}
	return nil
}

func (ps *ParkingSystem) SortTickets() {
	sort.Slice(ps.Tickets, func(i, j int) bool {
		return strings.Compare(ps.Tickets[i].VehicleNumber, ps.Tickets[j].VehicleNumber) < 0
	})
	ps.updateTicketIndices()
}

func (ps *ParkingSystem) updateTicketIndices() {
	for i := range ps.Tickets {
		ps.Tickets[i].Index = i + 1
	}
}

func CalculateFee(entryTime, exitTime time.Time, vehicleType string) float64 {
	duration := exitTime.Sub(entryTime).Hours()
	var feePerHour float64
	if vehicleType == "mobil" {
		feePerHour = FeePerHourCar
	} else if vehicleType == "motor" {
		feePerHour = FeePerHourBike
	}
	return feePerHour * duration
}

func (ps *ParkingSystem) PrintReport() {
	carCount, bikeCount := 0, 0
	totalFee := 0.0

	fmt.Println("Daftar Kendaraan yang Sudah Ditambahkan Tiket Parkir:")
	for _, ticket := range ps.Tickets {
		fmt.Printf("%d. %s - %s\n", ticket.Index, ticket.VehicleNumber, ticket.VehicleType)
		fmt.Printf("   Waktu Masuk: %s, Waktu Keluar: %s, Total Pembayaran: %.2f\n", ticket.EntryTime.Format("2006-01-02 15:04"), ticket.ExitTime.Format("2006-01-02 15:04"), ticket.Fee)
		if ticket.VehicleType == "mobil" {
			carCount++
		} else if ticket.VehicleType == "motor" {
			bikeCount++
		}
		totalFee += ticket.Fee
	}

	fmt.Printf("\nTotal Mobil: %d\n", carCount)
	fmt.Printf("Total Motor: %d\n", bikeCount)
	fmt.Printf("Total Saldo: %.2f\n", totalFee)
}

func (ps *ParkingSystem) SaveReportToFile(laporan string) {
	file, err := os.Create(laporan)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	carCount, bikeCount := 0, 0
	totalFee := 0.0

	writer := bufio.NewWriter(file)
	writer.WriteString("Daftar Kendaraan yang Sudah Ditambahkan Tiket Parkir:\n")
	for _, ticket := range ps.Tickets {
		writer.WriteString(fmt.Sprintf("%d. %s - %s\n", ticket.Index, ticket.VehicleNumber, ticket.VehicleType))
		writer.WriteString(fmt.Sprintf("   Waktu Masuk: %s, Waktu Keluar: %s, Total Pembayaran: %.2f\n", ticket.EntryTime.Format("2006-01-02 15:04"), ticket.ExitTime.Format("2006-01-02 15:04"), ticket.Fee))
		if ticket.VehicleType == "mobil" {
			carCount++
		} else if ticket.VehicleType == "motor" {
			bikeCount++
		}
		totalFee += ticket.Fee 
	}

	writer.WriteString(fmt.Sprintf("\nTotal Mobil: %d\n", carCount))
	writer.WriteString(fmt.Sprintf("Total Motor: %d\n", bikeCount))
	writer.WriteString(fmt.Sprintf("Total Saldo: %.2f\n", totalFee))

	writer.Flush()
	fmt.Println("Laporan berhasil disimpan ke file:", laporan)
}

func main() {
	parkingSystem := ParkingSystem{}
	reader := bufio.NewReader(os.Stdin)

	parkingSystem.AddUser(User{Username: "admin", Password: "admin123"})

	var username, password string
	fmt.Println("Selamat Datang di Aplikasi Parking")
	fmt.Print("username: ")
	username, _ = reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Print("password: ")
	password, _ = reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if parkingSystem.Authenticate(username, password) {
		fmt.Println("Admin logged in successfully!")

		for {
			var choice int
			fmt.Println("\nParking System Menu")
			fmt.Println("1. Tambahkan tiket")
			fmt.Println("2. Modifikasi Tiket")
			fmt.Println("3. Hapus Tiket")
			fmt.Println("4. Search Tiket")
			fmt.Println("5. Sortir Tikets")
			fmt.Println("6. Cetak Laporan")
			fmt.Println("7. Simpan Laporan ke File")
			fmt.Println("8. Keluar")
			fmt.Print("Masukkan menu pilihan Anda: ")
			fmt.Scanln(&choice)

			switch choice {
			case 1:
				var ticket ParkingTicket
				fmt.Print("Masukan plat nomor kendaraan: ")
				ticket.VehicleNumber, _ = reader.ReadString('\n')
				ticket.VehicleNumber = strings.TrimSpace(ticket.VehicleNumber)

				fmt.Print("Masukan tipe kendaraan (mobil/motor): ")
				ticket.VehicleType, _ = reader.ReadString('\n')
				ticket.VehicleType = strings.TrimSpace(ticket.VehicleType)

				fmt.Print("Masukan waktu masuk: ")
				entryTimeStr, _ := reader.ReadString('\n')
				entryTimeStr = strings.TrimSpace(entryTimeStr)
				entryTime, _ := time.Parse("2006-01-02 15:04", entryTimeStr)
				ticket.EntryTime = entryTime

				fmt.Print("Masukan waktu keluar: ")
				exitTimeStr, _ := reader.ReadString('\n')
				exitTimeStr = strings.TrimSpace(exitTimeStr)
				exitTime, _ := time.Parse("2006-01-02 15:04", exitTimeStr)
				ticket.ExitTime = exitTime

				ticket.Fee = CalculateFee(ticket.EntryTime, ticket.ExitTime, ticket.VehicleType)

				parkingSystem.AddTicket(ticket)
				fmt.Println("Tiket berhasil ditambahkan!")
				fmt.Printf("Total Pembayaran: %.2f\n", ticket.Fee)
			case 2:
				var index int
				var ticket ParkingTicket
				fmt.Print("Masukan Index Tiket untuk dimodifikasi: ")
				fmt.Scanln(&index)

				fmt.Print("Masukan plat nomor kendaraan: ")
				ticket.VehicleNumber, _ = reader.ReadString('\n')
				ticket.VehicleNumber = strings.TrimSpace(ticket.VehicleNumber)

				fmt.Print("Masukan tipe kendaraan (mobil/motor): ")
				ticket.VehicleType, _ = reader.ReadString('\n')
				ticket.VehicleType = strings.TrimSpace(ticket.VehicleType)

				fmt.Print("Masukan Waktu Masuk: ")
				entryTimeStr, _ := reader.ReadString('\n')
				entryTimeStr = strings.TrimSpace(entryTimeStr)
				entryTime, _ := time.Parse("2006-01-02 15:04", entryTimeStr)
				ticket.EntryTime = entryTime

				fmt.Print("Masukan Waktu Keluar: ")
				exitTimeStr, _ := reader.ReadString('\n')
				exitTimeStr = strings.TrimSpace(exitTimeStr)
				exitTime, _ := time.Parse("2006-01-02 15:04", exitTimeStr)
				ticket.ExitTime = exitTime

				ticket.Fee = CalculateFee(ticket.EntryTime, ticket.ExitTime, ticket.VehicleType)

				parkingSystem.ModifyTicket(index-1, ticket)
				fmt.Println("Tiket Berhasil di modifikasi!")
				fmt.Printf("Total Pembayaran: %.2f\n", ticket.Fee)
			case 3:
				var index int
				fmt.Print("Masukan Index tiket yang ingin dihapus: ")
				fmt.Scanln(&index)
				parkingSystem.DeleteTicket(index-1)
				fmt.Println("Tiket berhasil dihapus!")
			case 4:
				var vehicleNumber string
				fmt.Print("Masukkan Nomor Kendaraan yang Akan Dicari: ")
				vehicleNumber, _ = reader.ReadString('\n')
				vehicleNumber = strings.TrimSpace(vehicleNumber)
				ticket := parkingSystem.SearchTicket(vehicleNumber)
				if ticket != nil {
					fmt.Println("Tiket ditemukan:", *ticket)
				} else {
					fmt.Println("Tiket tidak ditemukan!")
				}
			case 5:
				parkingSystem.SortTickets()
				fmt.Println("Tiket berhasil di sortir!")
			case 6:
				parkingSystem.PrintReport()
			case 7:
				var filename string
				fmt.Print("Masukkan nama file untuk menyimpan laporan: ")
				filename, _ = reader.ReadString('\n')
				filename = strings.TrimSpace(filename)
				parkingSystem.SaveReportToFile(filename)
			case 8:
				fmt.Println("Exiting the system...")
				return
			default:
				fmt.Println("Invalid choice. Please try again.")
			}
		}
	} else {
		fmt.Println("Invalid login!")
	}
}
