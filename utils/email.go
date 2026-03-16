package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func KirimEmailReset(tujuanEmail string, resetLink string) error {
	// 1. Ambil Kunci Brevo dari Satpam Render
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return errors.New("Kunci BREVO_API_KEY belum dipasang di Render")
	}

	url := "https://api.brevo.com/v3/smtp/email"

	pesanHTML := `
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 2px solid #ffedd5; border-radius: 10px;">
			<h2 style="color: #ea580c;">Halo dari admin amdaay.scarf! </h2>
			<p>Kami menerima permintaan untuk mereset password Anda.</p>
			<p>Silakan klik tombol di bawah ini untuk membuat password baru:</p>
			<br>
			<a href="` + resetLink + `" style="background-color: #ea580c; color: white; padding: 12px 24px; text-decoration: none; border-radius: 8px; font-weight: bold; display: inline-block;">
				Ganti Password Saya
			</a>
			<br><br>
			<p style="font-size: 13px; color: #6b7280;">Atau copy-paste link rahasia berikut ke browser Anda:</p>
			<p style="font-size: 13px; color: #2563eb; word-break: break-all; background-color: #f3f4f6; padding: 10px; border-radius: 5px;">
				` + resetLink + `
			</p>
			<br>
			<p style="font-size: 12px; color: #ef4444;"><b>Perhatian:</b> Link ini hanya berlaku selama 5 menit.</p>
		</div>
	`

	// 2. Membungkus Paket Surat untuk Brevo
	payloadData := map[string]interface{}{
		"sender": map[string]string{
			"name":  "Admin Amdaay.Scarf",
			"email": "mandayusuf2728@gmail.com", // Menggunakan email Bundo sebagai pengirim
		},
		"to": []map[string]string{
			{
				"email": tujuanEmail,
			},
		},
		"subject":     "Reset Password - website resmi amdaay.scarf",
		"htmlContent": pesanHTML,
	}

	payloadBytes, _ := json.Marshal(payloadData)

	// 3. Mengirim Paket Lewat Jalur Utama (HTTP / Port 443)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Add("accept", "application/json")
	req.Header.Add("api-key", apiKey)
	req.Header.Add("content-type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 && res.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Ditolak Brevo! Status: %d, Balasan: %s", res.StatusCode, string(bodyBytes))
	}

	return nil
}

//kode untuk local testing

// package utils

// import (
// 	"crypto/tls"
// 	"gopkg.in/gomail.v2"
// )

// func KirimEmailReset(tujuanEmail string, resetLink string) error {
// 	m := gomail.NewMessage()

// 	m.SetHeader("From", "mandayusuf2728@gmail.com")
// 	m.SetHeader("To", tujuanEmail)
// 	m.SetHeader("Subject", "Reset Password - website resmi amdaay.scarf")

// 	pesanHTML := `
// 		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 2px solid #ffedd5; border-radius: 10px;">
// 			<h2 style="color: #ea580c;">Halo dari admin amdaay.scarf! </h2>
// 			<p>Kami menerima permintaan untuk mereset password Anda.</p>
// 			<p>Silakan klik tombol di bawah ini untuk membuat password baru:</p>
// 			<br>
// 			<a href="` + resetLink + `" style="background-color: #ea580c; color: white; padding: 12px 24px; text-decoration: none; border-radius: 8px; font-weight: bold; display: inline-block;">
// 				Ganti Password Saya
// 			</a>
// 			<br><br>
// 			<p style="font-size: 13px; color: #6b7280;">Atau copy-paste link rahasia berikut ke browser Anda:</p>
// 			<p style="font-size: 13px; color: #2563eb; word-break: break-all; background-color: #f3f4f6; padding: 10px; border-radius: 5px;">
// 				` + resetLink + `
// 			</p>
// 			<br>
// 			<p style="font-size: 12px; color: #ef4444;"><b>Perhatian:</b> Link ini hanya berlaku selama 5 menit.</p>
// 		</div>
// 	`

// 	m.SetBody("text/html", pesanHTML)

// 	dialer := gomail.NewDialer("smtp.gmail.com", 465, "mandayusuf2728@gmail.com", "tipbyxuusbrrmzwu")
// 	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

// 	return dialer.DialAndSend(m)
// }