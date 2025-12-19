# UAS-Backend
Project UAS Mata Kuliah Backend Pemrograman Lanjutan Praktikum 

Akbar Zahron Jiwa Yanu
434231028 - TI C2

# 🎓 Sistem Pelaporan Prestasi Mahasiswa

![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=flat&logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v2.52-00ACD7?style=flat&logo=fiber)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)
![MongoDB](https://img.shields.io/badge/MongoDB-7.0-47A248?style=flat&logo=mongodb)
![License](https://img.shields.io/badge/License-MIT-green.svg)

Backend API untuk sistem pelaporan dan verifikasi prestasi mahasiswa dengan arsitektur hybrid database (PostgreSQL + MongoDB) dan Role-Based Access Control (RBAC).

---

## 📋 Daftar Isi

- [Fitur Utama](#-fitur-utama)
- [Tech Stack](#-tech-stack)
- [Arsitektur Database](#-arsitektur-database)
- [Instalasi](#-instalasi)
- [Konfigurasi](#-konfigurasi)
- [Migrasi Database](#-migrasi-database)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Project Structure](#-project-structure)
- [Role & Permissions](#-role--permissions)
- [Kontribusi](#-kontribusi)
- [Author](#-author)
- [License](#-license)

---

## ✨ Fitur Utama

- 🔐 **Authentication & Authorization** dengan JWT Token
- 👥 **Role-Based Access Control (RBAC)** - Admin, Mahasiswa, Dosen Wali
- 📝 **Dynamic Achievement Fields** - Fleksibel sesuai tipe prestasi
- ✅ **Workflow Approval** - Draft → Submitted → Verified/Rejected
- 📊 **Statistics & Reports** - Dashboard prestasi mahasiswa
- 🗄️ **Hybrid Database** - PostgreSQL untuk relasional, MongoDB untuk dokumen dinamis
- 📱 **RESTful API** - Clean architecture dengan Fiber framework
- 🔒 **Data Validation** - Input validation & error handling
- 📝 **Structured Logging** - Zerolog untuk monitoring

---

## 🛠 Tech Stack

### Backend Framework
- **Go 1.25.5** - Programming language
- **Fiber v2.52** - Web framework (Express-inspired)

### Database
- **PostgreSQL 14+** - Relational data (users, roles, references)
- **MongoDB 7.0** - Document store (dynamic achievement data)

### Libraries
- **JWT (golang-jwt/jwt)** - Token-based authentication
- **Bcrypt** - Password hashing
- **Zerolog** - Structured logging
- **SQLX** - PostgreSQL driver
- **Mongo Driver** - MongoDB official driver

### Tools
- **golang-migrate** - Database migrations
- **godotenv** - Environment configuration
- **Postman** - API testing

---

## 🗃 Arsitektur Database

### PostgreSQL (Relational Data)
Menyimpan data yang terstruktur dan memerlukan relasi:

```
users ──┬── students ──── achievement_references ──→ MongoDB
        │
        └── lecturers
        
roles ── role_permissions ── permissions
```

**Tables:**
- `users` - Data user (admin, mahasiswa, dosen)
- `roles` - Role definitions
- `permissions` - Permission resources
- `role_permissions` - Many-to-many relation
- `students` - Profil mahasiswa
- `lecturers` - Profil dosen wali
- `achievement_references` - Referensi ke MongoDB dengan status workflow

### MongoDB (Dynamic Documents)
Menyimpan detail prestasi dengan field yang fleksibel:

```javascript
{
  _id: ObjectId,
  studentId: UUID,
  achievementType: "competition", // or "publication", "organization", etc.
  title: String,
  description: String,
  details: {
    // Dynamic fields based on achievement type
    competitionName: String,
    competitionLevel: String,
    rank: Number,
    // ... other fields
  },
  attachments: [...],
  tags: [...],
  points: Number
}
```

---

## 🚀 Instalasi

### Prerequisites

Pastikan sudah terinstall:
- Go 1.21 atau lebih baru
- PostgreSQL 14+
- MongoDB 7.0+
- golang-migrate (untuk migrasi database)

### Clone Repository

```bash
git clone https://github.com/bomboskuy/UAS-Backend.git
cd UAS-Backend
```

### Install Dependencies

```bash
go mod download
```

### Install golang-migrate

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

**Windows:**
```bash
# Download dari https://github.com/golang-migrate/migrate/releases
# Extract dan tambahkan ke PATH
```

---

## ⚙️ Konfigurasi

### 1. Setup Environment Variables

Copy file `.env.example` menjadi `.env`:

```bash
cp .env.example .env
```

Edit file `.env`:

```env
# PostgreSQL Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=uas_backend
DB_SSLMODE=disable

# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DB=UAS-Backend

# JWT Configuration
JWT_SECRET=your_super_secret_key_here
JWT_EXPIRED=24
JWT_REFRESH_SECRET=your_refresh_secret_key_here
JWT_REFRESH_EXPIRED=168
```

### 2. Setup PostgreSQL Database

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE uas_backend;

# Keluar
\q
```

### 3. Setup MongoDB

MongoDB tidak perlu setup khusus. Database akan dibuat otomatis saat aplikasi pertama kali berjalan.

---

## 🗄 Migrasi Database

### Run Migrations

```bash
# Migrate up (terapkan semua migrasi)
migrate -database "postgres://postgres:your_password@localhost:5432/uas_backend?sslmode=disable" -path db/migrations up

# Migrate down (rollback semua migrasi)
migrate -database "postgres://postgres:your_password@localhost:5432/uas_backend?sslmode=disable" -path db/migrations down

# Migrate to specific version
migrate -database "postgres://..." -path db/migrations goto 2
```

### Migrations Include:

1. **000001** - Create initial tables (users, roles, permissions, students, lecturers, achievement_references)
2. **000002** - Seed default roles (Admin, Mahasiswa, Dosen Wali)
3. **000003** - Seed permissions and role_permissions mapping
4. **000004** - Create default admin and dosen accounts
5. **000005** - Add foreign key constraints
6. **000006** - Add achievement status enum

---

## 🏃 Menjalankan Aplikasi

### Development Mode

```bash
go run main.go
```

### Build & Run

```bash
# Build binary
go build -o uas-backend

# Run binary
./uas-backend
```

### Using Air (Hot Reload)

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

Server akan berjalan di: **http://localhost:3000**

---

## 📚 API Documentation

### Base URL
```
http://localhost:3000/api/v1
```

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/login` | Login user | ❌ |
| GET | `/auth/profile` | Get user profile | ✅ |
| POST | `/auth/logout` | Logout user | ✅ |

### User Management (Admin Only)

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/users` | Get all users | `user:manage` |
| GET | `/users/:id` | Get user by ID | `user:manage` |
| POST | `/users` | Create new user | `user:manage` |
| PUT | `/users/:id` | Update user | `user:manage` |
| DELETE | `/users/:id` | Delete user | `user:manage` |

### Achievement Management

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| POST | `/achievements` | Create achievement | `achievement:create` |
| GET | `/achievements` | Get achievements | `achievement:read` |
| GET | `/achievements/:id` | Get achievement detail | `achievement:read` |
| PUT | `/achievements/:id` | Update achievement | `achievement:update` |
| DELETE | `/achievements/:id` | Delete achievement | `achievement:delete` |
| POST | `/achievements/:id/submit` | Submit for verification | `achievement:update` |
| POST | `/achievements/:id/verify` | Verify achievement | `achievement:verify` |
| POST | `/achievements/:id/reject` | Reject achievement | `achievement:verify` |

### Student & Lecturer Management

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/students` | Get all students | `user:manage` |
| GET | `/students/:id` | Get student by ID | `user:manage` |
| PUT | `/students/:id/advisor` | Assign advisor | `user:manage` |
| GET | `/lecturers` | Get all lecturers | `user:manage` |
| GET | `/lecturers/:id/advisees` | Get advisees | `achievement:verify` |

### Reports & Analytics

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/reports/statistics` | Get achievement statistics | `achievement:read` |
| GET | `/reports/student/:id` | Get student report | `achievement:read` |

### Sample Requests

**Login:**
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**Create Achievement:**
```bash
curl -X POST http://localhost:3000/api/v1/achievements \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "achievement_type": "competition",
    "title": "Juara 1 Hackathon Nasional",
    "description": "Memenangkan kompetisi hackathon",
    "details": {
      "competition_name": "Indonesia Tech Challenge",
      "competition_level": "national",
      "rank": 1
    },
    "tags": ["teknologi", "hackathon"],
    "points": 100
  }'
```

---

## 🧪 Testing

### Setup Postman Collection

1. Import collection: [`docs/postman/UAS-Backend.postman_collection.json`](docs/postman)
2. Import environment: [`docs/postman/UAS-Backend.postman_environment.json`](docs/postman)
3. Run collection atau manual testing

### Testing Flow

1. **Login Admin** → Save token
2. **Create Dosen Wali** → Get lecturer_id
3. **Create Mahasiswa** → Assign advisor
4. **Login Mahasiswa** → Save token
5. **Create Achievement** (draft)
6. **Submit Achievement** (submitted)
7. **Login Dosen** → Save token
8. **Verify/Reject Achievement**
9. **Get Statistics & Reports**

### Default Credentials

| Role | Username | Password |
|------|----------|----------|
| Admin | `admin` | `admin123` |
| Dosen Wali | `dosen` | `dosen123` |

---

## 📁 Project Structure

```
UAS-Backend/
├── app/
│   ├── models/              # Data structures
│   │   ├── achievement.go
│   │   ├── user.go
│   │   ├── auth.go
│   │   └── ...
│   ├── repositories/        # Database layer
│   │   ├── achievement_repository_mongo.go
│   │   ├── user_repository_pg.go
│   │   └── ...
│   └── services/           # Business logic
│       ├── achievement_service.go
│       ├── auth_service.go
│       └── user_service.go
├── db/
│   ├── migrations/         # Database migrations
│   ├── mongo.go           # MongoDB connection
│   └── pg.go              # PostgreSQL connection
├── helper/
│   └── response.go        # Response helpers
├── middleware/
│   └── permissions.go     # RBAC middleware
├── routes/
│   └── index.go          # Route definitions
├── utils/
│   ├── jwt.go            # JWT utilities
│   ├── logger.go         # Logging setup
│   └── password.go       # Password hashing
├── .env.example          # Environment template
├── go.mod
├── go.sum
├── main.go              # Application entry point
└── README.md
```

---

## 👥 Role & Permissions

### Admin
**Full Access:**
- ✅ Manage users (CRUD)
- ✅ Manage students & lecturers
- ✅ View all achievements
- ✅ View statistics & reports
- ✅ Assign advisors

### Mahasiswa
**Student Access:**
- ✅ Create achievement (draft)
- ✅ Update own achievement (draft only)
- ✅ Submit achievement for verification
- ✅ View own achievements
- ❌ Verify/reject achievements

### Dosen Wali
**Advisor Access:**
- ✅ View advisees' achievements
- ✅ Verify achievements
- ✅ Reject achievements with notes
- ✅ View advisees list
- ❌ Create/update achievements

---

## 🔄 Achievement Workflow

```mermaid
graph LR
    A[Draft] -->|Submit| B[Submitted]
    B -->|Verify| C[Verified]
    B -->|Reject| D[Rejected]
    D -->|Revise| A
    A -->|Delete| E[Deleted]
```

**Status Flow:**
1. **Draft** - Mahasiswa membuat prestasi (bisa edit/delete)
2. **Submitted** - Mahasiswa submit untuk verifikasi (tidak bisa edit)
3. **Verified** - Dosen menyetujui prestasi
4. **Rejected** - Dosen menolak dengan catatan (mahasiswa bisa revisi)

---

## 🤝 Kontribusi

Kontribusi sangat diterima! Silakan ikuti langkah berikut:

1. Fork repository ini
2. Buat branch fitur (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

### Development Guidelines

- Ikuti Go coding conventions
- Tambahkan unit tests untuk fitur baru
- Update dokumentasi API jika ada perubahan endpoint
- Gunakan meaningful commit messages

---

## 🐛 Known Issues & TODO

### Known Issues
- [ ] File upload untuk attachments belum diimplementasi
- [ ] Refresh token belum fully implemented
- [ ] Rate limiting belum ada

### TODO
- [ ] Implement file upload (AWS S3/MinIO)
- [ ] Add email notification system
- [ ] Implement caching (Redis)
- [ ] Add API rate limiting
- [ ] Add comprehensive unit tests
- [ ] Swagger/OpenAPI documentation
- [ ] Docker & Docker Compose setup
- [ ] CI/CD pipeline

---

## 👨‍💻 Author

**Akbar Zahron Jiwa Yanu**
- NIM: 434231028
- Kelas: TI C2
- Program Studi: D-IV Teknik Informatika
- Universitas: Universitas Airlangga

**Mata Kuliah:** Pemrograman Backend Lanjut (Praktikum)

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Fiber Framework - https://gofiber.io/
- PostgreSQL - https://www.postgresql.org/
- MongoDB - https://www.mongodb.com/
- Universitas Airlangga - Fakultas Sains dan Teknologi

---

## 📞 Support

Jika ada pertanyaan atau issue, silakan:
- Buat [GitHub Issue](https://github.com/bomboskuy/UAS-Backend/issues)
- Email: akbar.zahron@example.com

---

<div align="center">

Made with ❤️ by Akbar Zahron Jiwa Yanu

**⭐ Star this repository if you find it helpful!**

</div>
```

---

## 💡 Tips Tambahan untuk Repository

### 1. Buat file `.env.example`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=uas_backend
DB_SSLMODE=disable

MONGO_URI=mongodb://localhost:27017
MONGO_DB=UAS-Backend

JWT_SECRET=your_jwt_secret_here
JWT_EXPIRED=24
JWT_REFRESH_SECRET=your_refresh_secret_here
JWT_REFRESH_EXPIRED=168
```

### 2. Buat `.gitignore`:
```
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
uas-backend

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Testing
*.test
coverage.out
```

### 3. Tambahkan folder `docs/`:
```
docs/
├── postman/
│   ├── UAS-Backend.postman_collection.json
│   └── UAS-Backend.postman_environment.json
├── api/
│   └── API_DOCUMENTATION.md
└── database/
    └── ERD.png
```

### 4. Buat LICENSE file (MIT):
```
MIT License

Copyright (c) 2025 Akbar Zahron Jiwa Yanu

Permission is hereby granted, free of charge...
```

README ini sudah lengkap dan professional untuk GitHub repository! 🚀
