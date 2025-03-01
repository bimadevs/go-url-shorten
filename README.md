# URL Shortener API

*Dibuat oleh Bimadev*

URL Shortener adalah layanan API untuk memendekkan URL panjang menjadi URL pendek yang lebih mudah dibagikan. Aplikasi ini dibangun menggunakan Go dengan framework Gin dan database SQLite.

## Fitur Utama

- **Pemendekkan URL**: Mengubah URL panjang menjadi URL pendek
- **Custom Alias**: Opsi untuk menentukan nama kustom untuk URL pendek
- **Sistem API Key**: Keamanan dengan API key untuk mengakses layanan
- **Statistik Klik**: Melacak berapa kali URL pendek diakses
- **CORS Support**: Mendukung integrasi dengan aplikasi frontend

## Instalasi dan Penggunaan

### Prasyarat

- Go (versi 1.16 atau lebih baru)
- SQLite
- Git

### Langkah Instalasi

1. Clone repositori
   ```bash
   git clone https://github.com/bimadevs/go-url-shorten.git
   cd go-url-shorten
   ```

2. Install dependensi
   ```bash
   # Inisialisasi modul Go jika belum ada file go.mod
   go mod init url-shortener

   # Install dependensi yang diperlukan
   go get github.com/gin-gonic/gin
   go get gorm.io/gorm
   go get gorm.io/driver/sqlite
   ```

3. Jalankan aplikasi
   ```bash
   go run main.go
   ```

Aplikasi akan berjalan di `http://localhost:8080`

## Penggunaan API

### 1. Mendapatkan API Key

**Endpoint**: `POST /generate-key`

**Contoh Request**:
```bash
curl -X POST http://localhost:8080/generate-key
```

**Contoh Response**:
```json
{
  "api_key": "Bimadevxxxxxxxx"
}
```

### 2. Memendekkan URL

**Endpoint**: `POST /shorten`

**Headers**:
- `X-API-Key`: API key yang didapatkan sebelumnya
- `Content-Type`: application/json

**Body**:
```json
{
  "original_url": "https://example.com/very/long/url/that/needs/shortening",
  "custom_alias": "myurl"  // Opsional
}
```

**Contoh Request**:
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -H "X-API-Key: Bimadevxxxxxxxx" \
  -d '{"original_url": "https://example.com/path", "custom_alias": "myurl"}'
```

**Contoh Response**:
```json
{
  "short_url": "http://localhost:8080/myurl"
}
```

### 3. Mengakses URL Pendek

**Endpoint**: `GET /:short`

Cukup akses URL pendek melalui browser atau aplikasi HTTP client, dan Anda akan diarahkan ke URL asli.

### 4. Melihat Statistik URL

**Endpoint**: `GET /stats/:short`

**Contoh Request**:
```bash
curl http://localhost:8080/stats/myurl
```

**Contoh Response**:
```json
{
  "short_url": "http://localhost:8080/myurl",
  "original_url": "https://example.com/path",
  "click_count": 5
}
```

## Aturan dan Batasan

- **URL**: Harus dalam format yang valid
- **Alias**: Hanya huruf dan angka diperbolehkan, maksimal 15 karakter
- **API Key**: Wajib digunakan untuk membuat URL pendek

## Integrasi dengan Frontend

API ini mendukung CORS, sehingga dapat dengan mudah diintegrasikan dengan aplikasi frontend seperti React, Vue, atau Angular.

Contoh penggunaan dengan JavaScript:

```javascript
// Mendapatkan API Key
async function getApiKey() {
  const response = await fetch('http://localhost:8080/generate-key', {
    method: 'POST'
  });
  return await response.json();
}

// Memendekkan URL
async function shortenUrl(apiKey, originalUrl, customAlias = '') {
  const response = await fetch('http://localhost:8080/shorten', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': apiKey
    },
    body: JSON.stringify({
      original_url: originalUrl,
      custom_alias: customAlias
    })
  });
  return await response.json();
}
```

## Struktur Database

### Tabel URL
- `ID`: Primary key
- `ShortCode`: Kode pendek unik untuk URL
- `OriginalURL`: URL asli
- `ClickCount`: Jumlah klik
- `APIKey`: API key yang digunakan untuk membuat URL pendek

### Tabel User
- `ID`: Primary key
- `APIKey`: API key unik untuk pengguna

## Lisensi

Kode ini open source. Silakan gunakan dengan tetap mencantumkan kredit kepada developer asli.

---

Happy Coding :_)

*Source Code ini open source untuk orang. Jangan hapus ini untuk menghargai Developer.*
